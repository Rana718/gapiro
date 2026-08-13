import { uid, emptyPair } from '../lib/utils';
import type {
  Pair, HttpMethod, BodyType, RequestTab, ResponseTab, ResponseViewMode,
  AuthConfig, RequestConfig, RequestSettings, ResponseData, Folder,
  Protocol, OpenTab
} from '../lib/types';

// ─── Factories ──────────────────────────────────────────────────────────────

function defaultSettings(): RequestSettings {
  return { timeout: 30, followRedirects: true, verifySSL: true, maxRedirects: 10 };
}

function defaultAuth(): AuthConfig {
  return { type: 'none' };
}

const DEFAULT_GRAPHQL_QUERY = 'query {\n  \n}';

function defaultName(p: Protocol): string {
  return p === 'http' ? 'New Request'
    : p === 'graphql' ? 'New GraphQL Query'
    : p === 'grpc' ? 'New gRPC Call'
    : 'New WebSocket';
}

/** Build a fresh request config for a protocol, with protocol-appropriate defaults. */
export function newRequestConfig(protocol: Protocol, name?: string): RequestConfig {
  const now = Date.now();
  const cfg: RequestConfig = {
    id: uid('rq'),
    name: name ?? defaultName(protocol),
    protocol,
    method: protocol === 'graphql' ? 'POST' : 'GET',
    url: '',
    headers: [emptyPair()],
    urlParameters: [emptyPair()],
    bodyType: 'none',
    body: '',
    formData: [emptyPair()],
    auth: defaultAuth(),
    settings: defaultSettings(),
    description: '',
    preRequestScript: '',
    postResponseScript: '',
    createdAt: now,
    updatedAt: now,
  };
  if (protocol === 'graphql') {
    cfg.graphql = { query: DEFAULT_GRAPHQL_QUERY, variables: '{}' };
    cfg.headers = [
      { id: uid('p'), name: 'Content-Type', value: 'application/json', enabled: true },
      emptyPair(),
    ];
  } else if (protocol === 'grpc') {
    cfg.grpc = { protoFile: '', importPaths: [], fullMethod: '', message: '{}', metadata: [emptyPair()] };
  } else if (protocol === 'websocket') {
    cfg.websocket = { protocols: '' };
  }
  return cfg;
}

// ─── Live Request State ─────────────────────────────────────────────────────
// Mirrors the active tab. All editors bind directly to this object.

export const request = $state<RequestConfig>(newRequestConfig('http'));

// ─── Response State ───────────────────────────────────────────────────────────
// $state.raw for data to avoid deep proxy on large response bodies (can be MBs).

export const response = {
  get loading() { return _responseLoading; },
  set loading(v: boolean) { _responseLoading = v; },
  get data() { return _responseData; },
  set data(v: ResponseData | null) { _responseData = v; },
  get error() { return _responseError; },
  set error(v: string | null) { _responseError = v; },
};
let _responseLoading = $state(false);
let _responseData: ResponseData | null = $state.raw(null);
let _responseError: string | null = $state(null);

/** Cancel hook — http.ts registers its cancelRequest here to avoid a circular import. */
export const requestRuntime = { cancel: null as (null | (() => void)) };

// ─── UI State ─────────────────────────────────────────────────────────────────

export const ui = $state({
  activeRequestTab: 'body' as RequestTab,
  activeResponseTab: 'body' as ResponseTab,
  responseViewMode: 'pretty' as ResponseViewMode,
  sidebarWidth: 250,
  sidebarHidden: false,
  splitRatio: 0.5,
  splitLayout: 'vertical' as 'horizontal' | 'vertical',
});

// ─── Collection State ─────────────────────────────────────────────────────────

export const collection = $state({
  requests: [] as RequestConfig[],
  folders: [] as Folder[],
  activeRequestId: null as string | null,
});

// ─── Tabs State ─────────────────────────────────────────────────────────────

export const tabs = $state({
  open: [] as OpenTab[],
  activeId: null as string | null,
});

// ─── Persistence ──────────────────────────────────────────────────────────────

if (typeof localStorage !== 'undefined') {
  try {
    const saved = localStorage.getItem('gapiro:collection');
    if (saved) {
      const parsed = JSON.parse(saved);
      collection.requests = (parsed.requests ?? []).map(migrateConfig);
      collection.folders = parsed.folders ?? [];
    }
  } catch { /* ignore corrupt local state */ }
}

function persistCollection() {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('gapiro:collection', JSON.stringify({
      requests: collection.requests,
      folders: collection.folders,
    }));
  }
}

function persistTabs() {
  if (typeof localStorage === 'undefined') return;
  // Responses are transient — persist only the request configs + which tab is active.
  const slim = tabs.open.map(t => ({
    id: t.id,
    protocol: t.protocol,
    request: t.id === tabs.activeId ? snapshotLiveRequest() : t.request,
    savedId: t.savedId,
    dirty: t.dirty,
  }));
  localStorage.setItem('gapiro:tabs', JSON.stringify({ open: slim, activeId: tabs.activeId }));
}

/** Back-fill fields on configs saved before the protocol/tab model existed. */
function migrateConfig(c: any): RequestConfig {
  return {
    ...c,
    protocol: (c.protocol as Protocol) ?? 'http',
    headers: c.headers ?? [],
    urlParameters: c.urlParameters ?? [],
    formData: c.formData ?? [],
    auth: c.auth ?? defaultAuth(),
    settings: c.settings ?? defaultSettings(),
    preRequestScript: c.preRequestScript ?? '',
    postResponseScript: c.postResponseScript ?? '',
  };
}

// ─── Live ↔ Snapshot helpers ────────────────────────────────────────────────

function withTrailingEmpty(pairs: Pair[] | undefined): Pair[] {
  const arr = pairs ? pairs.map(p => ({ ...p })) : [];
  const last = arr[arr.length - 1];
  if (!last || last.name !== '' || last.value !== '') arr.push(emptyPair());
  return arr;
}

/** Detached deep copy of the live request (for storing in a tab / collection). */
function snapshotLiveRequest(): RequestConfig {
  return $state.snapshot(request) as RequestConfig;
}

// Baseline of the active tab's request at load/save time — used for dirty detection.
let activeBaseline = '';

function serializeLive(): string {
  const s = snapshotLiveRequest() as any;
  // Ignore volatile/derived fields so cosmetic changes don't flag the tab dirty.
  return JSON.stringify({ ...s, id: undefined, createdAt: undefined, updatedAt: undefined });
}

/** Recompute the active tab's dirty flag against its baseline. Cheap; call from an effect. */
export function refreshDirty() {
  const t = tabs.open.find(t => t.id === tabs.activeId);
  if (!t) return;
  const dirty = serializeLive() !== activeBaseline;
  if (t.dirty !== dirty) t.dirty = dirty;
}

/** Assign a config into the live request object, ensuring editors have a trailing blank row. */
function hydrateLiveRequest(cfg: RequestConfig) {
  request.id = cfg.id;
  request.name = cfg.name;
  request.protocol = cfg.protocol ?? 'http';
  request.method = cfg.method;
  request.url = cfg.url;
  request.headers = withTrailingEmpty(cfg.headers);
  request.urlParameters = withTrailingEmpty(cfg.urlParameters);
  request.bodyType = cfg.bodyType;
  request.body = cfg.body;
  request.formData = withTrailingEmpty(cfg.formData);
  request.auth = structuredClone($state.snapshot(cfg.auth) ?? defaultAuth());
  request.settings = { ...cfg.settings };
  request.description = cfg.description;
  request.preRequestScript = cfg.preRequestScript ?? '';
  request.postResponseScript = cfg.postResponseScript ?? '';
  request.folderId = cfg.folderId;
  request.graphql = cfg.graphql ? { ...cfg.graphql } : undefined;
  request.grpc = cfg.grpc
    ? { ...cfg.grpc, importPaths: [...(cfg.grpc.importPaths ?? [])], metadata: withTrailingEmpty(cfg.grpc.metadata) }
    : undefined;
  request.websocket = cfg.websocket ? { ...cfg.websocket } : undefined;
}

/** Write the live request + response back into the active tab's snapshot. */
function syncActiveTab() {
  const t = tabs.open.find(t => t.id === tabs.activeId);
  if (!t) return;
  t.request = snapshotLiveRequest();
  t.protocol = request.protocol;
  t.response = { data: response.data, error: response.error };
}

// ─── Tab Actions ──────────────────────────────────────────────────────────────

function makeTab(cfg: RequestConfig, savedId: string | null): OpenTab {
  return {
    id: uid('tab'),
    protocol: cfg.protocol,
    request: cfg,
    response: { data: null, error: null },
    savedId,
    dirty: false,
  };
}

/** Open a brand-new draft request for a protocol and focus it. */
export function newDraftTab(protocol: Protocol, name?: string): OpenTab {
  const cfg = newRequestConfig(protocol, name);
  const tab = makeTab(cfg, null);
  syncActiveTab();
  tabs.open.push(tab);
  activate(tab);
  persistTabs();
  return tab;
}

/** Open a saved request in a tab (reusing an existing tab if already open). */
export function openSavedTab(saved: RequestConfig) {
  const existing = tabs.open.find(t => t.savedId === saved.id);
  if (existing) { switchTab(existing.id); return; }
  const tab = makeTab(structuredClone($state.snapshot(saved)) as RequestConfig, saved.id);
  syncActiveTab();
  tabs.open.push(tab);
  activate(tab);
  persistTabs();
}

/** Switch to an already-open tab, preserving the current tab's edits. */
export function switchTab(id: string) {
  if (id === tabs.activeId) return;
  if (response.loading) requestRuntime.cancel?.();
  syncActiveTab();
  const t = tabs.open.find(t => t.id === id);
  if (!t) return;
  activate(t);
  persistTabs();
}

function activate(t: OpenTab) {
  tabs.activeId = t.id;
  hydrateLiveRequest(t.request);
  response.loading = false;
  response.data = t.response.data;
  response.error = t.response.error;
  collection.activeRequestId = t.savedId;
  activeBaseline = serializeLive();
}

/** Close a tab; focus a neighbour, or open a fresh draft if it was the last one. */
export function closeTab(id: string) {
  const idx = tabs.open.findIndex(t => t.id === id);
  if (idx === -1) return;
  const wasActive = tabs.activeId === id;
  if (wasActive && response.loading) requestRuntime.cancel?.();
  tabs.open.splice(idx, 1);

  if (!wasActive) { persistTabs(); return; }

  if (tabs.open.length === 0) {
    tabs.activeId = null;
    newDraftTab('http');
    return;
  }
  const next = tabs.open[Math.min(idx, tabs.open.length - 1)];
  activate(next);
  persistTabs();
}

// ─── Request Actions ──────────────────────────────────────────────────────────

export function resetRequest() {
  hydrateLiveRequest(newRequestConfig('http'));
  response.data = null;
  response.error = null;
  collection.activeRequestId = null;
}

/** Load a saved request into the active tab (legacy single-doc entry point). */
export function loadRequest(saved: RequestConfig) {
  openSavedTab(saved);
}

export function saveCurrentRequest(name?: string) {
  const cleanPairs = (pairs: Pair[]) => pairs.filter(p => p.name !== '');
  const activeTab = tabs.open.find(t => t.id === tabs.activeId);
  const savedId = activeTab?.savedId ?? collection.activeRequestId;
  const existing = collection.requests.find(r => r.id === savedId);

  const fields = {
    name: name ?? request.name,
    protocol: request.protocol,
    method: request.method,
    url: request.url,
    headers: cleanPairs(request.headers),
    urlParameters: cleanPairs(request.urlParameters),
    bodyType: request.bodyType,
    body: request.body,
    formData: cleanPairs(request.formData),
    auth: structuredClone($state.snapshot(request.auth)),
    settings: { ...request.settings },
    description: request.description,
    preRequestScript: request.preRequestScript,
    postResponseScript: request.postResponseScript,
    folderId: request.folderId,
    graphql: request.graphql ? { ...request.graphql } : undefined,
    grpc: request.grpc
      ? { ...request.grpc, importPaths: [...request.grpc.importPaths], metadata: cleanPairs(request.grpc.metadata) }
      : undefined,
    websocket: request.websocket ? { ...request.websocket } : undefined,
    updatedAt: Date.now(),
  };

  let targetId: string;
  if (existing) {
    Object.assign(existing, fields);
    targetId = existing.id;
  } else {
    const newReq: RequestConfig = {
      id: uid('rq'),
      createdAt: Date.now(),
      ...fields,
      name: name ?? (request.name || request.url || 'Untitled'),
    } as RequestConfig;
    collection.requests.push(newReq);
    targetId = newReq.id;
  }

  request.name = fields.name;
  collection.activeRequestId = targetId;
  if (activeTab) { activeTab.savedId = targetId; activeTab.dirty = false; }
  activeBaseline = serializeLive();
  persistCollection();
  persistTabs();
}

export function deleteRequest(id: string) {
  const idx = collection.requests.findIndex(r => r.id === id);
  if (idx !== -1) collection.requests.splice(idx, 1);
  // Detach any open tabs that pointed at this saved request.
  for (const t of tabs.open) if (t.savedId === id) { t.savedId = null; t.dirty = true; }
  if (collection.activeRequestId === id) collection.activeRequestId = null;
  persistCollection();
  persistTabs();
}

export function duplicateRequest(id: string) {
  const orig = collection.requests.find(r => r.id === id);
  if (!orig) return;
  const dup: RequestConfig = {
    ...structuredClone($state.snapshot(orig)),
    id: uid('rq'),
    name: `${orig.name} (copy)`,
    createdAt: Date.now(),
    updatedAt: Date.now(),
  };
  collection.requests.push(dup);
  persistCollection();
}

export function createFolder(name: string, parentId?: string) {
  collection.folders.push({ id: uid('fl'), name, parentId, expanded: true });
  persistCollection();
}

/** Assign a saved request to a folder (or to the root when folderId is undefined). */
export function moveRequestToFolder(requestId: string, folderId: string | undefined) {
  const req = collection.requests.find(r => r.id === requestId);
  if (!req) return;
  req.folderId = folderId;
  req.updatedAt = Date.now();
  if (collection.activeRequestId === requestId) request.folderId = folderId;
  persistCollection();
}

// ─── Bootstrap: restore open tabs, or start with one HTTP draft ─────────────────

(function initTabs() {
  if (typeof localStorage !== 'undefined') {
    try {
      const raw = localStorage.getItem('gapiro:tabs');
      if (raw) {
        const parsed = JSON.parse(raw);
        const restored: OpenTab[] = (parsed.open ?? []).map((t: any) => ({
          id: t.id ?? uid('tab'),
          protocol: (t.protocol as Protocol) ?? 'http',
          request: migrateConfig(t.request),
          response: { data: null, error: null },
          savedId: t.savedId ?? null,
          dirty: t.dirty ?? false,
        }));
        if (restored.length) {
          tabs.open = restored;
          const active = restored.find(t => t.id === parsed.activeId) ?? restored[0];
          activate(active);
          return;
        }
      }
    } catch { /* ignore corrupt tab state */ }
  }
  // First launch (or corrupt state): open a single HTTP draft.
  const tab = makeTab(newRequestConfig('http'), null);
  tabs.open = [tab];
  activate(tab);
})();
