// Reactive stores for the application state (Svelte 5 runes)
import { uid, emptyKV } from '../lib/utils';
import type { KeyValue, HttpMethod, BodyType, RequestTab, ResponseTab, ResponsePayload, SavedRequest, Folder } from '../lib/types';

// ─── Request State ────────────────────────────────────────────────────

export const requestState = $state({
  method: 'GET' as HttpMethod,
  url: '',
  headers: [emptyKV()] as KeyValue[],
  queryParams: [emptyKV()] as KeyValue[],
  bodyType: 'none' as BodyType,
  body: '',
  formData: [emptyKV()] as KeyValue[],
  timeout: 30,
  followRedirects: true,
  verifySSL: true,
});

// ─── Response State ───────────────────────────────────────────────────

export const responseState = $state({
  loading: false,
  response: null as ResponsePayload | null,
  error: null as string | null,
});

// ─── UI State ─────────────────────────────────────────────────────────

export const uiState = $state({
  activeRequestTab: 'params' as RequestTab,
  activeResponseTab: 'body' as ResponseTab,
  sidebarWidth: 240,
  sidebarCollapsed: false,
  splitRatio: 0.5, // request/response split
  responseBodyPretty: true,
});

// ─── Collection State ─────────────────────────────────────────────────

export const collectionState = $state({
  requests: [] as SavedRequest[],
  folders: [] as Folder[],
  activeRequestId: null as string | null,
});

// ─── Actions ──────────────────────────────────────────────────────────

export function resetRequest() {
  requestState.method = 'GET';
  requestState.url = '';
  requestState.headers = [emptyKV()];
  requestState.queryParams = [emptyKV()];
  requestState.bodyType = 'none';
  requestState.body = '';
  requestState.formData = [emptyKV()];
  requestState.timeout = 30;
  requestState.followRedirects = true;
  requestState.verifySSL = true;
}

export function loadRequest(saved: SavedRequest) {
  requestState.method = saved.method;
  requestState.url = saved.url;
  requestState.headers = [...saved.headers, emptyKV()];
  requestState.queryParams = [...saved.queryParams, emptyKV()];
  requestState.bodyType = saved.bodyType;
  requestState.body = saved.body;
  requestState.formData = [...saved.formData, emptyKV()];
  collectionState.activeRequestId = saved.id;
}

export function saveCurrentRequest(name?: string) {
  const existing = collectionState.requests.find(r => r.id === collectionState.activeRequestId);

  if (existing) {
    existing.method = requestState.method;
    existing.url = requestState.url;
    existing.headers = requestState.headers.filter(h => h.key !== '');
    existing.queryParams = requestState.queryParams.filter(p => p.key !== '');
    existing.bodyType = requestState.bodyType;
    existing.body = requestState.body;
    existing.formData = requestState.formData.filter(f => f.key !== '');
    existing.updatedAt = Date.now();
  } else {
    const newReq: SavedRequest = {
      id: uid('rq'),
      name: name || requestState.url || 'Untitled Request',
      method: requestState.method,
      url: requestState.url,
      headers: requestState.headers.filter(h => h.key !== ''),
      queryParams: requestState.queryParams.filter(p => p.key !== ''),
      bodyType: requestState.bodyType,
      body: requestState.body,
      formData: requestState.formData.filter(f => f.key !== ''),
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };
    collectionState.requests.push(newReq);
    collectionState.activeRequestId = newReq.id;
  }
}

export function deleteRequest(id: string) {
  const idx = collectionState.requests.findIndex(r => r.id === id);
  if (idx !== -1) {
    collectionState.requests.splice(idx, 1);
    if (collectionState.activeRequestId === id) {
      collectionState.activeRequestId = null;
      resetRequest();
    }
  }
}

export function createFolder(name: string, parentId?: string) {
  collectionState.folders.push({
    id: uid('fl'),
    name,
    parentId,
    expanded: true,
  });
}
