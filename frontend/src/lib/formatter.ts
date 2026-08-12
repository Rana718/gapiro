// Off-main-thread JSON formatter using Web Worker
// Prevents UI blocking when formatting large responses

let worker: Worker | null = null;
let requestId = 0;
const pending = new Map<number, { resolve: (v: string) => void; reject: (e: any) => void }>();

function getWorker(): Worker {
  if (!worker) {
    worker = new Worker(new URL('./format.worker.ts', import.meta.url), { type: 'module' });
    worker.onmessage = (e) => {
      const { id, result } = e.data;
      const p = pending.get(id);
      if (p) {
        pending.delete(id);
        p.resolve(result);
      }
    };
    worker.onerror = (e) => {
      // If worker fails, reject all pending
      for (const [id, p] of pending) {
        p.reject(e);
        pending.delete(id);
      }
    };
  }
  return worker;
}

/**
 * Format JSON in a Web Worker (non-blocking).
 * Falls back to sync if worker unavailable.
 */
export function formatJSONAsync(json: string): Promise<string> {
  // Small payloads: format synchronously (no overhead)
  if (json.length < 10_000) {
    try {
      return Promise.resolve(JSON.stringify(JSON.parse(json), null, 2));
    } catch {
      return Promise.resolve(json);
    }
  }

  // Large payloads: use worker
  return new Promise((resolve, reject) => {
    const id = ++requestId;
    pending.set(id, { resolve, reject });
    try {
      getWorker().postMessage({ type: 'format-json', data: json, id });
    } catch {
      pending.delete(id);
      // Fallback sync
      try {
        resolve(JSON.stringify(JSON.parse(json), null, 2));
      } catch {
        resolve(json);
      }
    }

    // Timeout: don't hang forever
    setTimeout(() => {
      if (pending.has(id)) {
        pending.delete(id);
        resolve(json); // Return raw if worker too slow
      }
    }, 5000);
  });
}

/** Terminate the worker (cleanup) */
export function terminateWorker(): void {
  if (worker) {
    worker.terminate();
    worker = null;
    pending.clear();
  }
}
