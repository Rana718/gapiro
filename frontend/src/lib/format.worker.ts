// Web Worker for heavy JSON operations off the main thread
// This prevents UI freezing when formatting large JSON responses

self.onmessage = (e: MessageEvent) => {
  const { type, data, id } = e.data;

  switch (type) {
    case 'format-json': {
      try {
        const parsed = JSON.parse(data);
        const formatted = JSON.stringify(parsed, null, 2);
        self.postMessage({ id, result: formatted, error: null });
      } catch {
        self.postMessage({ id, result: data, error: null });
      }
      break;
    }
    case 'search-body': {
      const { body, query } = data;
      const lines = body.split('\n');
      const matches: number[] = [];
      const lowerQuery = query.toLowerCase();
      for (let i = 0; i < lines.length; i++) {
        if (lines[i].toLowerCase().includes(lowerQuery)) {
          matches.push(i);
        }
      }
      self.postMessage({ id, result: matches, error: null });
      break;
    }
    default:
      self.postMessage({ id, result: null, error: 'Unknown message type' });
  }
};
