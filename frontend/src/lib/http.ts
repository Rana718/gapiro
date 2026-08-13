import { request, response, requestRuntime } from '../stores/app.svelte';
import type { ResponseData } from './types';

let HttpService: any = null;
let CollectionService: any = null;
let GraphQLService: any = null;
let GrpcService: any = null;
let WebSocketService: any = null;
let Events: any = null;

// Load Wails bindings + runtime (resolves once; awaited by executors below).
export const bindingsReady: Promise<void> = (async () => {
  try {
    const mod = await import('../../bindings/changeme/services');
    HttpService = mod?.HttpService ?? null;
    CollectionService = mod?.CollectionService ?? null;
    GraphQLService = mod?.GraphQLService ?? null;
    GrpcService = mod?.GrpcService ?? null;
    WebSocketService = (mod as any)?.WebSocketService ?? null;
  } catch { /* browser dev — bindings absent, fall back to fetch */ }
  try {
    const rt = await import('@wailsio/runtime');
    Events = rt?.Events ?? null;
  } catch { /* runtime absent in browser dev */ }
})();

// ─── Shared execution state ───────────────────────────────────────────────────

let activePromise: any = null; // Wails CancellablePromise
let abortController: AbortController | null = null;
let requestGeneration = 0;
let activeRequestId: string | null = null;

/** Build a full ResponseData from a partial, filling the timing/meta fields. */
function makeResponseData(p: Partial<ResponseData>): ResponseData {
  return {
    status: 0, statusText: '', headers: {}, body: '', size: 0,
    duration: 0, dnsTime: 0, connectTime: 0, tlsTime: 0, ttfbTime: 0,
    protocol: '', remoteAddr: '', contentType: '', redirectCount: 0,
    ...p,
    size: p.size ?? (p.body ? new Blob([p.body]).size : 0),
  };
}

// ─── Protocol dispatcher ──────────────────────────────────────────────────────

/** Run the active request against the correct backend, keyed off request.protocol. */
export async function execute(): Promise<void> {
  switch (request.protocol) {
    case 'graphql': return executeGraphQL();
    case 'grpc': return executeGrpc();
    // WebSocket is driven by its own pane (connect/send), not this one-shot path.
    case 'websocket': return;
    default: return sendRequest();
  }
}

// ─── HTTP ─────────────────────────────────────────────────────────────────────

export async function sendRequest(): Promise<void> {
  if (!request.url || response.loading) return;

  // Cancel any in-flight request first
  cancelRequest();
  const generation = ++requestGeneration;
  const requestId = `http_${Date.now()}_${generation}`;
  activeRequestId = requestId;

  response.loading = true;
  response.error = null;
  response.data = null;
  // Safety valve for a bridge that cannot deliver cancellation (e.g. blocked native DNS).
  const watchdog = window.setTimeout(() => {
    if (generation === requestGeneration && response.loading) cancelRequest();
  }, Math.max(60000, (request.settings.timeout || 30) * 1000 + 5000));

  const payload = {
    requestId,
    method: request.method,
    url: request.url,
    headers: request.headers.filter(h => h.name !== '').map(h => ({ key: h.name, value: h.value, enabled: h.enabled })),
    queryParams: request.urlParameters.filter(p => p.name !== '').map(p => ({ key: p.name, value: p.value, enabled: p.enabled })),
    bodyType: request.bodyType,
    body: request.body,
    formData: request.formData.filter(f => f.name !== '').map(f => ({ key: f.name, value: f.value, enabled: f.enabled })),
    timeout: request.settings.timeout,
    followRedirects: request.settings.followRedirects,
    verifySSL: request.settings.verifySSL,
  };

  try {
    await bindingsReady;
    let result: ResponseData;

    if (HttpService) {
      // Wails binding — returns CancellablePromise
      // Store the raw Call promise for cancellation (before .then transforms it)
      const promise = HttpService.SendRequest(payload);
      activePromise = promise;
      result = await promise;
    } else {
      // Browser fetch fallback
      abortController = new AbortController();
      result = await fetchFallback(payload, abortController.signal);
    }

    if (generation === requestGeneration) {
      // Go zero values arrive as null for maps. Normalize at the boundary so
      // cancelled/error responses can never crash reactive UI rendering.
      result.headers = result.headers ?? {};
      result.body = result.body ?? '';
      result.contentType = result.contentType ?? '';
      result.protocol = result.protocol ?? '';
      result.remoteAddr = result.remoteAddr ?? '';
      response.data = result;
      if (result.error) response.error = result.error;
    }
  } catch (err: any) {
    if (generation === requestGeneration) {
      const msg = err?.message ?? String(err) ?? 'Request failed';
      // Don't show "cancelled" as an error
      if (!msg.includes('cancel') && !msg.includes('abort')) {
        response.error = msg;
      }
    }
  } finally {
    window.clearTimeout(watchdog);
    if (generation === requestGeneration) {
      response.loading = false;
    }
    if (generation === requestGeneration) {
      activePromise = null;
      abortController = null;
      activeRequestId = null;
    }
  }
}

export function cancelRequest(): void {
  requestGeneration++;
  // UI state must be committed before touching the native bridge. A stalled
  // Wails call can otherwise prevent Svelte/WebKit from painting this change.
  response.loading = false;
  response.error = null;
  const promise = activePromise;
  const controller = abortController;
  const requestId = activeRequestId;
  activePromise = null;
  abortController = null;
  activeRequestId = null;

  // Yield a frame so the Cancel button visibly responds even when the native
  // bridge is wedged. Abort work only after the UI has painted.
  requestAnimationFrame(() => {
    window.setTimeout(() => {
      if (controller) controller.abort('Cancelled by user');
      if (requestId && HttpService?.CancelRequest) {
        try { void HttpService.CancelRequest(requestId); } catch { /* best effort */ }
      }
      if (promise && typeof promise.cancel === 'function') {
        try { void promise.cancel('Cancelled by user'); } catch { /* expected */ }
      }
    }, 0);
  });
}

// Let the store cancel in-flight work when switching/closing tabs (avoids a circular import).
requestRuntime.cancel = cancelRequest;

// ─── GraphQL ────────────────────────────────────────────────────────────────

/** Combine a GraphQL data/errors pair into one displayable JSON document. */
function mergeGraphQLBody(data: string, errors: string): string {
  const parse = (s: string, fallback: any) => { try { return s ? JSON.parse(s) : fallback; } catch { return s || fallback; } };
  if (errors) return JSON.stringify({ data: parse(data, null), errors: parse(errors, null) }, null, 2);
  return data || 'null';
}

/** Execute the active GraphQL request and write the result into the response store. */
export async function executeGraphQL(): Promise<void> {
  if (!request.url || response.loading) return;
  cancelRequest();
  const generation = ++requestGeneration;

  response.loading = true;
  response.error = null;
  response.data = null;
  const start = performance.now();

  const gql = request.graphql ?? { query: '', variables: '{}' };
  const headers = request.headers.filter(h => h.name !== '').map(h => ({ name: h.name, value: h.value, enabled: h.enabled }));

  try {
    await bindingsReady;
    let result: any;
    if (GraphQLService) {
      const promise = GraphQLService.ExecuteQuery({ url: request.url, query: gql.query, variables: gql.variables, headers });
      activePromise = promise;
      result = await promise;
    } else {
      result = await graphqlFetchFallback(request.url, gql.query, gql.variables, headers);
    }
    if (generation !== requestGeneration) return;

    const duration = result?.duration ?? (performance.now() - start);
    if (result?.error) {
      response.error = result.error;
      response.data = makeResponseData({ status: 0, statusText: 'Error', duration, contentType: 'application/json', protocol: 'GraphQL' });
      return;
    }
    const hasErrors = !!result?.errors;
    response.data = makeResponseData({
      status: 200,
      statusText: hasErrors ? '200 (with errors)' : 'OK',
      body: mergeGraphQLBody(result?.data ?? '', result?.errors ?? ''),
      duration,
      contentType: 'application/json',
      protocol: 'GraphQL',
    });
  } catch (err: any) {
    if (generation === requestGeneration) {
      const msg = err?.message ?? String(err) ?? 'GraphQL request failed';
      if (!msg.includes('cancel') && !msg.includes('abort')) response.error = msg;
    }
  } finally {
    if (generation === requestGeneration) {
      response.loading = false;
      activePromise = null;
    }
  }
}

export async function sendGraphQLQuery(url: string, query: string, variables: string, headers: any[]): Promise<any> {
  await bindingsReady;
  if (GraphQLService) return await GraphQLService.ExecuteQuery({ url, query, variables, headers });
  return await graphqlFetchFallback(url, query, variables, headers);
}

export async function introspectGraphQL(url: string, headers: any[]): Promise<any> {
  await bindingsReady;
  if (GraphQLService) return await GraphQLService.IntrospectSchema(url, headers);
  return null;
}

/** Browser-dev GraphQL POST fallback (mirrors the Go ExecuteQuery response shape). */
async function graphqlFetchFallback(url: string, query: string, variables: string, headers: any[]): Promise<any> {
  const start = performance.now();
  try {
    let vars: any = {};
    try { vars = variables ? JSON.parse(variables) : {}; } catch { /* leave empty */ }
    const h: Record<string, string> = { 'Content-Type': 'application/json' };
    headers.filter((x: any) => x.enabled && x.name).forEach((x: any) => { h[x.name] = x.value; });
    const resp = await fetch(url, { method: 'POST', headers: h, body: JSON.stringify({ query, variables: vars }) });
    const json = await resp.json().catch(() => ({}));
    return {
      data: json?.data !== undefined ? JSON.stringify(json.data, null, 2) : '',
      errors: json?.errors !== undefined ? JSON.stringify(json.errors, null, 2) : '',
      duration: performance.now() - start,
    };
  } catch (err: any) {
    return { error: err?.message ?? 'GraphQL request failed', duration: performance.now() - start };
  }
}

// ─── gRPC ───────────────────────────────────────────────────────────────────

export async function parseProtoFile(filePath: string, importPaths: string[]): Promise<any> {
  await bindingsReady;
  if (GrpcService) return await GrpcService.ParseProtoFile(filePath, importPaths);
  return [];
}

export async function sendGrpcUnary(address: string, protoFile: string, importPaths: string[], fullMethod: string, message: string, metadata: any[]): Promise<any> {
  await bindingsReady;
  if (GrpcService) return await GrpcService.SendUnary(address, protoFile, importPaths, fullMethod, message, metadata);
  return { error: 'gRPC requires the desktop app', duration: 0 };
}

/** Execute the active gRPC unary call and write the result into the response store. */
export async function executeGrpc(): Promise<void> {
  const g = request.grpc;
  if (!request.url || !g?.fullMethod || response.loading) return;
  cancelRequest();
  const generation = ++requestGeneration;

  response.loading = true;
  response.error = null;
  response.data = null;
  const start = performance.now();

  const metadata = (g.metadata ?? []).filter(m => m.name !== '').map(m => ({ id: m.id, name: m.name, value: m.value, enabled: m.enabled }));

  try {
    await bindingsReady;
    if (!GrpcService) {
      response.error = 'gRPC requires the desktop app (run with `task dev`)';
      response.loading = false;
      return;
    }
    const promise = GrpcService.SendUnary(request.url, g.protoFile, g.importPaths ?? [], g.fullMethod, g.message, metadata);
    activePromise = promise;
    const result = await promise;
    if (generation !== requestGeneration) return;

    const duration = result?.duration ?? (performance.now() - start);
    if (result?.error) {
      response.error = result.error;
      response.data = makeResponseData({ status: 0, statusText: 'RPC Error', duration, headers: result?.metadata ?? {}, contentType: 'application/grpc', protocol: 'gRPC' });
      return;
    }
    response.data = makeResponseData({
      status: 200,
      statusText: 'OK',
      body: result?.data ?? '',
      headers: result?.metadata ?? {},
      duration,
      contentType: 'application/json',
      protocol: 'gRPC',
    });
  } catch (err: any) {
    if (generation === requestGeneration) {
      const msg = err?.message ?? String(err) ?? 'gRPC call failed';
      if (!msg.includes('cancel') && !msg.includes('abort')) response.error = msg;
    }
  } finally {
    if (generation === requestGeneration) {
      response.loading = false;
      activePromise = null;
    }
  }
}

// ─── WebSocket ────────────────────────────────────────────────────────────────

export async function wsConnect(connID: string, url: string, headers: any[], subprotocols = ''): Promise<void> {
  await bindingsReady;
  if (!WebSocketService) throw new Error('WebSocket requires the desktop app (run with `task dev`)');
  return await WebSocketService.Connect(connID, url, headers ?? [], subprotocols);
}

export async function wsSend(connID: string, message: string): Promise<void> {
  await bindingsReady;
  if (!WebSocketService) throw new Error('WebSocket requires the desktop app');
  return await WebSocketService.Send(connID, message);
}

export async function wsClose(connID: string): Promise<void> {
  await bindingsReady;
  if (WebSocketService) { try { await WebSocketService.Close(connID); } catch { /* already closed */ } }
}

/** Subscribe to streamed WebSocket frames. Returns an unsubscribe fn. */
export async function wsOnEvent(handler: (payload: any) => void): Promise<() => void> {
  await bindingsReady;
  if (!Events?.On) return () => {};
  return Events.On('ws:event', (ev: any) => handler(ev?.data));
}

// ─── Collection persistence ─────────────────────────────────────────────────

export async function loadCollection(): Promise<{ requests: any[]; folders: any[] }> {
  await bindingsReady;
  if (CollectionService) {
    const [requests, folders] = await Promise.all([
      CollectionService.ListRequests('ws_default'),
      CollectionService.ListFolders('ws_default'),
    ]);
    return { requests: requests ?? [], folders: folders ?? [] };
  }
  return { requests: [], folders: [] };
}

export async function saveRequestToDb(req: any): Promise<any> {
  if (CollectionService) return await CollectionService.SaveRequest(req);
  return req;
}

export async function deleteRequestFromDb(id: string): Promise<void> {
  if (CollectionService) await CollectionService.DeleteRequest(id);
}

export async function duplicateRequestInDb(id: string): Promise<any> {
  if (CollectionService) return await CollectionService.DuplicateRequest(id);
  return null;
}

export async function loadEnvironments(): Promise<any[]> {
  if (CollectionService) return await CollectionService.ListEnvironments('ws_default') ?? [];
  return [];
}

export async function saveEnvironment(env: any): Promise<any> {
  if (CollectionService) return await CollectionService.SaveEnvironment(env);
  return env;
}

export async function deleteEnvironment(id: string): Promise<void> {
  if (CollectionService) await CollectionService.DeleteEnvironment(id);
}

export async function getResponseHistory(requestId: string): Promise<any[]> {
  if (CollectionService) return await CollectionService.ListResponses(requestId, 20) ?? [];
  return [];
}

// ─── Fetch Fallback (browser dev mode) ──────────────────────────────────────

async function fetchFallback(payload: any, signal: AbortSignal): Promise<ResponseData> {
  const start = performance.now();
  try {
    let url = payload.url;
    if (!url.includes('://')) url = 'https://' + url;

    const u = new URL(url);
    payload.queryParams.filter((p: any) => p.enabled && p.key).forEach((p: any) => {
      u.searchParams.append(p.key, p.value);
    });

    const headers: Record<string, string> = {};
    payload.headers.filter((h: any) => h.enabled && h.key).forEach((h: any) => {
      headers[h.key] = h.value;
    });

    const init: RequestInit = { method: payload.method, headers, signal };

    if (payload.bodyType !== 'none' && !['GET', 'HEAD'].includes(payload.method)) {
      if (['json', 'text', 'xml', 'graphql'].includes(payload.bodyType)) {
        init.body = payload.body;
      } else if (payload.bodyType === 'form-urlencoded') {
        const params = new URLSearchParams();
        payload.formData.filter((f: any) => f.enabled && f.key).forEach((f: any) => params.append(f.key, f.value));
        init.body = params;
      }
    }

    const resp = await fetch(u.toString(), init);
    const body = await resp.text();
    const duration = performance.now() - start;

    const respHeaders: Record<string, string> = {};
    resp.headers.forEach((v, k) => { respHeaders[k] = v; });

    return {
      status: resp.status,
      statusText: `${resp.status} ${resp.statusText}`,
      headers: respHeaders,
      body,
      size: new Blob([body]).size,
      duration,
      dnsTime: 0,
      connectTime: 0,
      tlsTime: 0,
      ttfbTime: duration * 0.3,
      protocol: 'HTTP/2',
      remoteAddr: '',
      contentType: resp.headers.get('content-type') ?? '',
      redirectCount: 0,
    };
  } catch (err: any) {
    if (signal.aborted) throw new Error('Request cancelled');
    return {
      status: 0, statusText: '', headers: {}, body: '', size: 0,
      duration: performance.now() - start,
      dnsTime: 0, connectTime: 0, tlsTime: 0, ttfbTime: 0,
      error: err?.message ?? 'Request failed',
      protocol: '', remoteAddr: '', contentType: '', redirectCount: 0,
    };
  }
}
