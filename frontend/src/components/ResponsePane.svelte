<!--
  ResponsePane - The right/bottom panel showing response data with tabs.
-->
<script lang="ts">
  import { responseState, uiState } from '../stores/app.svelte';
  import type { ResponseTab } from '../lib/types';
  import { formatBytes, formatDuration, statusColor, isJSON, prettyJSON } from '../lib/utils';
  import TabBar from './TabBar.svelte';
  import CodeEditor from './CodeEditor.svelte';
  import StatusBadge from './StatusBadge.svelte';
  import TimingChart from './TimingChart.svelte';

  const tabs = $derived([
    { id: 'body', label: 'Body' },
    { id: 'headers', label: 'Headers', badge: responseState.response ? Object.keys(responseState.response.headers).length : undefined },
    { id: 'timing', label: 'Timing' },
    { id: 'info', label: 'Info' },
  ]);

  const formattedBody = $derived(() => {
    if (!responseState.response?.body) return '';
    if (uiState.responseBodyPretty && isJSON(responseState.response.contentType)) {
      return prettyJSON(responseState.response.body);
    }
    return responseState.response.body;
  });

  const bodyLanguage = $derived(() => {
    if (!responseState.response?.contentType) return 'text';
    const ct = responseState.response.contentType;
    if (ct.includes('json')) return 'json';
    if (ct.includes('html')) return 'html';
    if (ct.includes('xml')) return 'xml';
    return 'text';
  });
</script>

<div class="flex flex-col h-full bg-surface-1 rounded-lg border border-border-subtle gpu-layer">
  {#if responseState.loading}
    <!-- Loading state -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3">
      <div class="w-6 h-6 border-2 border-accent/30 border-t-accent rounded-full animate-spin"></div>
      <span class="text-xs text-text-muted">Sending request...</span>
    </div>
  {:else if responseState.error && !responseState.response}
    <!-- Error state -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3 p-6">
      <svg class="w-8 h-8 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"/>
      </svg>
      <span class="text-sm text-error text-center">{responseState.error}</span>
    </div>
  {:else if responseState.response}
    <!-- Status bar -->
    <div class="flex items-center gap-3 px-3 py-1.5 border-b border-border-subtle">
      <StatusBadge status={responseState.response.status} statusText={responseState.response.statusText} />
      <span class="text-[11px] text-text-muted">{formatDuration(responseState.response.duration)}</span>
      <span class="text-[11px] text-text-muted">{formatBytes(responseState.response.size)}</span>
      {#if responseState.response.redirectCount > 0}
        <span class="text-[11px] text-warning">{responseState.response.redirectCount} redirect(s)</span>
      {/if}

      <!-- Pretty toggle for body tab -->
      {#if uiState.activeResponseTab === 'body' && isJSON(responseState.response.contentType)}
        <button
          onclick={() => { uiState.responseBodyPretty = !uiState.responseBodyPretty; }}
          class="ml-auto px-2 py-0.5 text-[10px] font-medium rounded cursor-pointer
            {uiState.responseBodyPretty
              ? 'bg-accent/15 text-accent'
              : 'text-text-muted hover:text-text-secondary'}"
        >
          Pretty
        </button>
      {/if}
    </div>

    <!-- Tabs -->
    <TabBar
      {tabs}
      active={uiState.activeResponseTab}
      onchange={(id) => { uiState.activeResponseTab = id as ResponseTab; }}
    />

    <!-- Tab content -->
    <div class="flex-1 overflow-auto">
      {#if uiState.activeResponseTab === 'body'}
        <CodeEditor
          value={formattedBody()}
          readonly
          language={bodyLanguage()}
        />
      {:else if uiState.activeResponseTab === 'headers'}
        <div class="p-2">
          <table class="w-full text-xs">
            <thead>
              <tr class="text-left text-text-muted border-b border-border-subtle">
                <th class="px-2 py-1.5 font-semibold">Header</th>
                <th class="px-2 py-1.5 font-semibold">Value</th>
              </tr>
            </thead>
            <tbody>
              {#each Object.entries(responseState.response.headers) as [key, value] (key)}
                <tr class="border-b border-border-subtle/50 hover:bg-surface-2/50">
                  <td class="px-2 py-1.5 text-accent font-medium">{key}</td>
                  <td class="px-2 py-1.5 text-text-primary font-mono break-all">{value}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {:else if uiState.activeResponseTab === 'timing'}
        <TimingChart
          dns={responseState.response.dnsTime}
          connect={responseState.response.connectTime}
          tls={responseState.response.tlsTime}
          ttfb={responseState.response.ttfbTime}
          total={responseState.response.duration}
        />
      {:else if uiState.activeResponseTab === 'info'}
        <div class="p-4 flex flex-col gap-2 text-xs">
          <div class="flex justify-between">
            <span class="text-text-muted">Protocol</span>
            <span class="text-text-primary font-mono">{responseState.response.protocol}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-text-muted">Remote Address</span>
            <span class="text-text-primary font-mono">{responseState.response.remoteAddr || 'N/A'}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-text-muted">Content-Type</span>
            <span class="text-text-primary font-mono">{responseState.response.contentType || 'N/A'}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-text-muted">Response Size</span>
            <span class="text-text-primary font-mono">{formatBytes(responseState.response.size)}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-text-muted">Total Time</span>
            <span class="text-text-primary font-mono">{formatDuration(responseState.response.duration)}</span>
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <!-- Empty state -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3 text-text-muted">
      <svg class="w-12 h-12 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
        <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5"/>
      </svg>
      <p class="text-sm">Enter a URL and click Send</p>
      <p class="text-[11px] text-text-muted/60">Response will appear here</p>
    </div>
  {/if}
</div>
