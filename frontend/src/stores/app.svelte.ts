import { uid, emptyPair } from '../lib/utils';
import type {
  Pair, HttpMethod, BodyType, RequestTab, ResponseTab, ResponseViewMode,
  AuthType, AuthConfig, RequestConfig, RequestSettings, ResponseData, ResponseState, Folder
} from '../lib/types';

// ─── Request State ────────────────────────────────────────────────────────────

function defaultSettings(): RequestSettings {
  return { timeout: 30, followRedirects: true, verifySSL: true, maxRedirects: 10 };
}

function defaultAuth(): AuthConfig {
  return { type: 'none' };
}

export const request = $state({
  method: 'GET' as HttpMethod,
  url: '',
  headers: [emptyPair()] as Pair[],
  urlParameters: [emptyPair()] as Pair[],
  bodyType: 'none' as BodyType,
  body: '',
  formData: [emptyPair()] as Pair[],
  auth: defaultAuth(),
  settings: defaultSettings(),
  description: '',
  name: '',
});

// ─── Response State ───────────────────────────────────────────────────────────

export const response: ResponseState = $state({
  loading: false,
  data: null,
  error: null,
});

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

// Keep the local workspace useful between launches (desktop apps should feel stateful).
if (typeof localStorage !== 'undefined') {
  try {
    const saved = localStorage.getItem('gapiro:collection');
    if (saved) Object.assign(collection, JSON.parse(saved));
  } catch { /* ignore corrupt local state */ }
}

function persistCollection() {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('gapiro:collection', JSON.stringify(collection));
  }
}

// ─── Actions ──────────────────────────────────────────────────────────────────

export function resetRequest() {
  request.method = 'GET';
  request.url = '';
  request.headers = [emptyPair()];
  request.urlParameters = [emptyPair()];
  request.bodyType = 'none';
  request.body = '';
  request.formData = [emptyPair()];
  request.auth = defaultAuth();
  request.settings = defaultSettings();
  request.description = '';
  request.name = '';
  collection.activeRequestId = null;
}

export function loadRequest(saved: RequestConfig) {
  request.method = saved.method;
  request.url = saved.url;
  request.headers = [...saved.headers, emptyPair()];
  request.urlParameters = [...saved.urlParameters, emptyPair()];
  request.bodyType = saved.bodyType;
  request.body = saved.body;
  request.formData = [...saved.formData, emptyPair()];
  request.auth = { ...saved.auth };
  request.settings = { ...saved.settings };
  request.description = saved.description;
  request.name = saved.name;
  collection.activeRequestId = saved.id;
}

export function saveCurrentRequest(name?: string) {
  const existing = collection.requests.find(r => r.id === collection.activeRequestId);
  const cleanPairs = (pairs: Pair[]) => pairs.filter(p => p.name !== '');

  if (existing) {
    Object.assign(existing, {
      name: name ?? existing.name,
      method: request.method,
      url: request.url,
      headers: cleanPairs(request.headers),
      urlParameters: cleanPairs(request.urlParameters),
      bodyType: request.bodyType,
      body: request.body,
      formData: cleanPairs(request.formData),
      auth: { ...request.auth },
      settings: { ...request.settings },
      description: request.description,
      updatedAt: Date.now(),
    });
  } else {
    const newReq: RequestConfig = {
      id: uid('rq'),
      name: name ?? (request.name || request.url || 'Untitled'),
      method: request.method,
      url: request.url,
      headers: cleanPairs(request.headers),
      urlParameters: cleanPairs(request.urlParameters),
      bodyType: request.bodyType,
      body: request.body,
      formData: cleanPairs(request.formData),
      auth: { ...request.auth },
      settings: { ...request.settings },
      description: request.description,
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };
    collection.requests.push(newReq);
    collection.activeRequestId = newReq.id;
  }
  persistCollection();
}

export function deleteRequest(id: string) {
  const idx = collection.requests.findIndex(r => r.id === id);
  if (idx !== -1) {
    collection.requests.splice(idx, 1);
    if (collection.activeRequestId === id) resetRequest();
  }
  persistCollection();
}

export function duplicateRequest(id: string) {
  const orig = collection.requests.find(r => r.id === id);
  if (!orig) return;
  const dup: RequestConfig = {
    ...structuredClone(orig),
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
