<!--
  App.svelte - Root application layout.
  Structure: TitleBar → Sidebar | (RequestPane / ResponsePane)
  GPU-composited panels with resizable splits.
-->
<script lang="ts">
  import { requestState, responseState } from './stores/app.svelte';
  import type { ResponsePayload } from './lib/types';
  import TitleBar from './components/TitleBar.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import ResponsePane from './components/ResponsePane.svelte';
  import ResizablePanel from './components/ResizablePanel.svelte';

  // Import the Wails-generated binding for HttpService
  // This will be available after `wails3 generate bindings`
  let HttpService: any;
  try {
    // Dynamic import - bindings generated at build time
    import('../bindings/changeme').then(mod => {
      HttpService = mod.HttpService;
    });
  } catch {
    // Fallback: will be available when bindings are generated
  }

  async function sendRequest() {
    if (!requestState.url) return;
    if (responseState.loading) return;

    responseState.loading = true;
    responseState.error = null;

    try {
      const payload = {
        method: requestState.method,
        url: requestState.url,
        headers: requestState.headers.filter(h => h.key !== ''),
        queryParams: requestState.queryParams.filter(p => p.key !== ''),
        bodyType: requestState.bodyType,
        body: requestState.body,
        formData: requestState.formData.filter(f => f.key !== ''),
        timeout: requestState.timeout,
        followRedirects: requestState.followRedirects,
        verifySSL: requestState.verifySSL,
      };

      let response: ResponsePayload;

      if (HttpService) {
        response = await HttpService.SendRequest(payload);
      } else {
        // Fallback for dev without Wails backend - use fetch
        response = await devFallbackRequest(payload);
      }

      if (response.error) {
        responseState.error = response.error;
      }
      responseState.response = response;
    } catch (err: any) {
      responseState.error = err?.message || 'Unknown error';
      responseState.response = null;
    } finally {
      responseState.loading = false;
    }
  }

  function cancelRequest() {
    // In a full implementation, this would abort the request via Wails
    responseState.loading = false;
  }

  /** Dev fallback using browser fetch (limited but functional for testing UI) */
  async function devFallbackRequest(payload: any): Promise<ResponsePayload> {
    const start = performance.now();
    try {
      let url = payload.url;
      if (!url.includes('://')) url = 'http://' + url;

      // Add query params
      const enabledParams = payload.queryParams.filter((p: any) => p.enabled && p.key);
      if (enabledParams.length > 0) {
        const u = new URL(url);
        enabledParams.forEach((p: any) => u.searchParams.append(p.key, p.value));
        url = u.toString();
      }

      const headers: Record<string, string> = {};
      payload.headers.filter((h: any) => h.enabled && h.key).forEach((h: any) => {
        headers[h.key] = h.value;
      });

      const init: RequestInit = {
        method: payload.method,
        headers,
      };

      if (payload.bodyType !== 'none' && !['GET', 'HEAD'].includes(payload.method)) {
        if (payload.bodyType === 'json' || payload.bodyType === 'text') {
          init.body = payload.body;
        }
      }

      const resp = await fetch(url, init);
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
        ttfbTime: Math.round(duration * 0.6),
        protocol: 'HTTP/2',
        remoteAddr: '',
        contentType: resp.headers.get('content-type') || '',
        redirectCount: 0,
      };
    } catch (err: any) {
      return {
        status: 0,
        statusText: '',
        headers: {},
        body: '',
        size: 0,
        duration: Math.round(performance.now() - start),
        dnsTime: 0,
        connectTime: 0,
        tlsTime: 0,
        ttfbTime: 0,
        error: err?.message || 'Request failed',
        protocol: '',
        remoteAddr: '',
        contentType: '',
        redirectCount: 0,
      };
    }
  }

  // Keyboard shortcut: Ctrl+Enter to send
  function handleGlobalKeydown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      sendRequest();
    }
    if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
      e.preventDefault();
      // New request shortcut
      requestState.url = '';
      requestState.method = 'GET';
    }
  }
</script>

<svelte:document onkeydown={handleGlobalKeydown} />

<div class="flex flex-col w-full h-full overflow-hidden gpu-layer">
  <!-- Title Bar -->
  <TitleBar />

  <!-- Main Content: Sidebar | Workspace -->
  <div class="flex-1 overflow-hidden">
    <ResizablePanel
      direction="horizontal"
      initialRatio={0.18}
      minFirst={180}
      minSecond={600}
      storageKey="sidebar"
    >
      {#snippet first()}
        <Sidebar />
      {/snippet}

      {#snippet second()}
        <!-- Workspace: Request / Response split -->
        <div class="flex flex-col h-full p-1.5 gap-1.5 bg-surface-2">
          <ResizablePanel
            direction="vertical"
            initialRatio={0.5}
            minFirst={200}
            minSecond={150}
            storageKey="req-res"
          >
            {#snippet first()}
              <RequestPane onSend={sendRequest} onCancel={cancelRequest} />
            {/snippet}

            {#snippet second()}
              <ResponsePane />
            {/snippet}
          </ResizablePanel>
        </div>
      {/snippet}
    </ResizablePanel>
  </div>
</div>
