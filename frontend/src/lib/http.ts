import { request, response } from '../stores/app.svelte';
import type { ResponseData } from './types';

let HttpService: any = null;

// Try to load Wails bindings (available at runtime)
try {
  import('../../bindings/changeme').then(mod => {
    HttpService = mod.HttpService;
  }).catch(() => {});
} catch {}

export async function sendRequest(): Promise<void> {
  if (!request.url || response.loading) return;

  response.loading = true;
  response.error = null;
  response.data = null;

  const payload = {
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
    let result: ResponseData;
    if (HttpService) {
      result = await HttpService.SendRequest(payload);
    } else {
      result = await fetchFallback(payload);
    }
    response.data = result;
    if (result.error) response.error = result.error;
  } catch (err: any) {
    response.error = err?.message ?? 'Request failed';
  } finally {
    response.loading = false;
  }
}

export function cancelRequest(): void {
  response.loading = false;
}

/** Browser fetch fallback for dev mode without Wails backend */
async function fetchFallback(payload: any): Promise<ResponseData> {
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

    const init: RequestInit = { method: payload.method, headers };

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
      duration: Math.round(duration),
      dnsTime: 0,
      connectTime: 0,
      tlsTime: 0,
      ttfbTime: Math.round(duration * 0.3),
      protocol: 'HTTP/2',
      remoteAddr: '',
      contentType: resp.headers.get('content-type') ?? '',
      redirectCount: 0,
    };
  } catch (err: any) {
    return {
      status: 0, statusText: '', headers: {}, body: '', size: 0,
      duration: Math.round(performance.now() - start),
      dnsTime: 0, connectTime: 0, tlsTime: 0, ttfbTime: 0,
      error: err?.message ?? 'Request failed',
      protocol: '', remoteAddr: '', contentType: '', redirectCount: 0,
    };
  }
}
