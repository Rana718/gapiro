// ─── Core Types ─────────────────────────────────────────────────────────────

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD' | 'OPTIONS';

export type BodyType = 'none' | 'json' | 'xml' | 'text' | 'form-urlencoded' | 'form-data' | 'graphql' | 'binary';

export type RequestTab = 'body' | 'params' | 'headers' | 'auth' | 'settings' | 'description';

export type ResponseTab = 'body' | 'request' | 'headers' | 'cookies' | 'timeline' | 'info';

export type ResponseViewMode = 'pretty' | 'raw';

export type AuthType = 'none' | 'basic' | 'bearer' | 'api-key';

// ─── Data Structures ────────────────────────────────────────────────────────

export interface Pair {
  id: string;
  name: string;
  value: string;
  enabled: boolean;
  readOnly?: boolean;
}

export interface AuthConfig {
  type: AuthType;
  basic?: { username: string; password: string };
  bearer?: { token: string; prefix: string };
  apiKey?: { key: string; value: string; addTo: 'header' | 'query' };
}

// ─── Request/Response ───────────────────────────────────────────────────────

export interface RequestConfig {
  id: string;
  name: string;
  method: HttpMethod;
  url: string;
  headers: Pair[];
  urlParameters: Pair[];
  bodyType: BodyType;
  body: string;
  formData: Pair[];
  auth: AuthConfig;
  settings: RequestSettings;
  description: string;
  folderId?: string;
  createdAt: number;
  updatedAt: number;
}

export interface RequestSettings {
  timeout: number;
  followRedirects: boolean;
  verifySSL: boolean;
  maxRedirects: number;
}

export interface ResponseData {
  status: number;
  statusText: string;
  headers: Record<string, string>;
  body: string;
  size: number;
  duration: number;
  dnsTime: number;
  connectTime: number;
  tlsTime: number;
  ttfbTime: number;
  error?: string;
  protocol: string;
  remoteAddr: string;
  contentType: string;
  redirectCount: number;
}

export interface ResponseState {
  loading: boolean;
  data: ResponseData | null;
  error: string | null;
}

// ─── Collection ─────────────────────────────────────────────────────────────

export interface Folder {
  id: string;
  name: string;
  parentId?: string;
  expanded: boolean;
}

// ─── Timeline Events ────────────────────────────────────────────────────────

export interface TimelineEvent {
  type: 'dns' | 'connect' | 'tls' | 'send_headers' | 'send_body' | 'receive_headers' | 'receive_body' | 'redirect' | 'info';
  timestamp: number;
  duration: number;
  detail?: string;
}
