// Utility functions used across the app

/** Generate a unique ID with prefix */
export function uid(prefix = ''): string {
  const rand = Math.random().toString(36).substring(2, 10);
  const ts = Date.now().toString(36);
  return prefix ? `${prefix}_${ts}${rand}` : `${ts}${rand}`;
}

/** Create an empty KeyValue row */
export function emptyKV(): import('./types').KeyValue {
  return { id: uid('kv'), key: '', value: '', enabled: true };
}

/** Format bytes to human readable */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
}

/** Format milliseconds to human readable duration */
export function formatDuration(ms: number): string {
  if (ms < 1) return '<1ms';
  if (ms < 1000) return `${Math.round(ms)}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(2)}s`;
  return `${(ms / 60000).toFixed(1)}m`;
}

/** Get color class for HTTP method */
export function methodColor(method: string): string {
  const colors: Record<string, string> = {
    GET: 'text-method-get',
    POST: 'text-method-post',
    PUT: 'text-method-put',
    PATCH: 'text-method-patch',
    DELETE: 'text-method-delete',
    HEAD: 'text-method-head',
    OPTIONS: 'text-method-options',
  };
  return colors[method] || 'text-text-secondary';
}

/** Get status color class */
export function statusColor(status: number): string {
  if (status >= 200 && status < 300) return 'text-success';
  if (status >= 300 && status < 400) return 'text-info';
  if (status >= 400 && status < 500) return 'text-warning';
  if (status >= 500) return 'text-error';
  return 'text-text-secondary';
}

/** Check if content type is JSON */
export function isJSON(contentType: string): boolean {
  return contentType?.includes('json') ?? false;
}

/** Check if content type is HTML */
export function isHTML(contentType: string): boolean {
  return contentType?.includes('html') ?? false;
}

/** Check if content type is XML */
export function isXML(contentType: string): boolean {
  return contentType?.includes('xml') ?? false;
}

/** Try to pretty-print JSON */
export function prettyJSON(str: string): string {
  try {
    return JSON.stringify(JSON.parse(str), null, 2);
  } catch {
    return str;
  }
}

/** Debounce a function */
export function debounce<T extends (...args: any[]) => any>(fn: T, ms: number): T {
  let timer: ReturnType<typeof setTimeout>;
  return ((...args: any[]) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), ms);
  }) as T;
}

/** Clamp a number between min and max */
export function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max);
}
