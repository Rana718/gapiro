import type { Pair, HttpMethod } from "./types";

/** Generate unique ID with optional prefix */
export function uid(prefix = ""): string {
   const r = Math.random().toString(36).slice(2, 10);
   const t = Date.now().toString(36);
   return prefix ? `${prefix}_${t}${r}` : `${t}${r}`;
}

/** Create empty key-value pair */
export function emptyPair(): Pair {
   return { id: uid("p"), name: "", value: "", enabled: true };
}

/** Format bytes */
export function formatBytes(bytes: number): string {
   if (bytes === 0) return "0 B";
   const k = 1024;
   const sizes = ["B", "KB", "MB", "GB"];
   const i = Math.floor(Math.log(bytes) / Math.log(k));
   return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
}

/** Format duration in ms — shows µs/ns for sub-millisecond */
export function formatDuration(ms: number): string {
   if (ms <= 0) return "0 ms";
   if (ms < 0.001) return `${Math.round(ms * 1_000_000)} ns`;
   if (ms < 1) return `${Math.round(ms * 1000)} µs`;
   if (ms < 1000) return `${Math.round(ms)} ms`;
   if (ms < 60000) return `${(ms / 1000).toFixed(2)} s`;
   return `${(ms / 60000).toFixed(1)} min`;
}

/** Protocol display metadata — icon, label, accent color, and the send-action verb. */
export const PROTOCOL_META: Record<
   string,
   { label: string; short: string; icon: string; color: string; text: string; action: string }
> = {
   http: { label: "HTTP Request", short: "HTTP", icon: "world", color: "protocol-http", text: "text-protocol-http", action: "Send" },
   graphql: { label: "GraphQL", short: "GQL", icon: "brand-graphql", color: "protocol-graphql", text: "text-protocol-graphql", action: "Query" },
   grpc: { label: "gRPC", short: "gRPC", icon: "server", color: "protocol-grpc", text: "text-protocol-grpc", action: "Invoke" },
   websocket: { label: "WebSocket", short: "WS", icon: "plug-connected", color: "protocol-websocket", text: "text-protocol-websocket", action: "Connect" },
};

/** Metadata for a protocol, falling back to HTTP for unknown values. */
export function protocolMeta(protocol: string) {
   return PROTOCOL_META[protocol] ?? PROTOCOL_META.http;
}

/** Get CSS color class for HTTP method */
export function methodColor(method: string): string {
   const map: Record<string, string> = {
      GET: "text-method-get",
      POST: "text-method-post",
      PUT: "text-method-put",
      PATCH: "text-method-patch",
      DELETE: "text-method-delete",
      HEAD: "text-method-head",
      OPTIONS: "text-method-options",
   };
   return map[method] ?? "text-text-subtle";
}

/** Get CSS color class for HTTP status */
export function statusColor(status: number): string {
   if (status >= 200 && status < 300) return "text-success";
   if (status >= 300 && status < 400) return "text-info";
   if (status >= 400 && status < 500) return "text-warning";
   if (status >= 500) return "text-danger";
   return "text-text-subtle";
}

/** Get background color for status badge */
export function statusBg(status: number): string {
   if (status >= 200 && status < 300) return "bg-success/10";
   if (status >= 300 && status < 400) return "bg-info/10";
   if (status >= 400 && status < 500) return "bg-warning/10";
   if (status >= 500) return "bg-danger/10";
   return "bg-surface-highlight";
}

/** Check if content type is JSON */
export function isJSON(ct: string | null | undefined): boolean {
   return ct?.includes("json") ?? false;
}

/** Check if content type is HTML */
export function isHTML(ct: string | null | undefined): boolean {
   return ct?.includes("html") ?? false;
}

/** Check if content type is XML */
export function isXML(ct: string | null | undefined): boolean {
   return ct?.includes("xml") ?? false;
}

/** Detect language from content type */
export function languageFromContentType(ct: string | null | undefined): string {
   if (!ct) return "text";
   if (ct.includes("json")) return "json";
   if (ct.includes("html")) return "html";
   if (ct.includes("xml")) return "xml";
   if (ct.includes("javascript")) return "javascript";
   if (ct.includes("css")) return "css";
   return "text";
}

/** Pretty print JSON (returns original on failure) */
export function prettyJSON(str: string): string {
   try {
      return JSON.stringify(JSON.parse(str), null, 2);
   } catch {
      return str;
   }
}

/** Clamp number */
export function clamp(value: number, min: number, max: number): number {
   return Math.min(Math.max(value, min), max);
}

/** Status code reason text */
export function statusReason(status: number): string {
   const reasons: Record<number, string> = {
      200: "OK",
      201: "Created",
      204: "No Content",
      301: "Moved Permanently",
      302: "Found",
      304: "Not Modified",
      400: "Bad Request",
      401: "Unauthorized",
      403: "Forbidden",
      404: "Not Found",
      405: "Method Not Allowed",
      408: "Request Timeout",
      409: "Conflict",
      422: "Unprocessable Entity",
      429: "Too Many Requests",
      500: "Internal Server Error",
      502: "Bad Gateway",
      503: "Service Unavailable",
      504: "Gateway Timeout",
   };
   return reasons[status] ?? "";
}

/** Default HTTP methods */
export const HTTP_METHODS: HttpMethod[] = [
   "GET",
   "POST",
   "PUT",
   "PATCH",
   "DELETE",
   "HEAD",
   "OPTIONS",
];

/** Default body types for dropdown */
export const BODY_TYPES = [
   { id: "none", label: "No Body", short: "Body" },
   { id: "json", label: "JSON", short: "JSON" },
   { id: "xml", label: "XML", short: "XML" },
   { id: "text", label: "Plain Text", short: "Text" },
   { id: "form-urlencoded", label: "Form URL Encoded", short: "Form" },
   { id: "form-data", label: "Multipart Form", short: "Multipart" },
   { id: "graphql", label: "GraphQL", short: "GraphQL" },
   { id: "binary", label: "Binary File", short: "Binary" },
] as const;

import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import type { Snippet } from "svelte";

export type WithElementRef<T, E extends HTMLElement = HTMLElement> = T & { ref?: E | null };
export type WithoutChild<T> = Omit<T, 'child'>;
export type WithoutChildren<T> = Omit<T, 'children'>;
export type WithoutChildrenOrChild<T> = Omit<T, 'children' | 'child'> & { children?: Snippet };

export function cn(...inputs: ClassValue[]) {
   return twMerge(clsx(inputs));
}
