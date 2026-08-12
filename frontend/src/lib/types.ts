// Types shared across the application

export interface KeyValue {
  id: string;
  key: string;
  value: string;
  enabled: boolean;
}

export interface RequestPayload {
  method: string;
  url: string;
  headers: KeyValue[];
  queryParams: KeyValue[];
  bodyType: BodyType;
  body: string;
  formData: KeyValue[];
  timeout: number;
  followRedirects: boolean;
  verifySSL: boolean;
}

export interface ResponsePayload {
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

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD' | 'OPTIONS';

export type BodyType = 'none' | 'json' | 'text' | 'form-urlencoded' | 'form-data' | 'binary';

export type RequestTab = 'params' | 'headers' | 'body' | 'auth' | 'settings';

export type ResponseTab = 'body' | 'headers' | 'timing' | 'info';

export interface SavedRequest {
  id: string;
  name: string;
  method: HttpMethod;
  url: string;
  headers: KeyValue[];
  queryParams: KeyValue[];
  bodyType: BodyType;
  body: string;
  formData: KeyValue[];
  folderId?: string;
  createdAt: number;
  updatedAt: number;
}

export interface Folder {
  id: string;
  name: string;
  parentId?: string;
  expanded: boolean;
}
