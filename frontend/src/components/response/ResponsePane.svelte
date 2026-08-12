<!--
  ResponsePane - Full response viewing panel.
  Structure: Status bar → Tabs (Body, Request, Headers, Cookies, Timeline)
  Matches Yaak's HttpResponsePane.
-->
<script lang="ts">
  import { response, ui } from '../../stores/app.svelte';
  import type { ResponseTab } from '../../lib/types';
  import { formatBytes, formatDuration, statusColor, statusBg, statusReason, isJSON, prettyJSON, languageFromContentType } from '../../lib/utils';
  import TabBar from '../core/TabBar.svelte';
  import ResponseBody from './ResponseBody.svelte';
  import ResponseHeaders from './ResponseHeaders.svelte';
  import ResponseTimeline from './ResponseTimeline.svelte';

  const tabs = $derived([
    { id: 'body', label: 'Response' },
    { id: 'headers', label: 'Headers', badge: response.data ? Object.keys(response.data.headers).length : undefined },
    { id: 'timeline', label: 'Timeline' },
    { id: 'info', label: 'Info' },
  ]);
</script>

<div class="flex flex-col h-full bg-surface rounded-md border border-border-subtle overflow-hidden gpu">
  {#if response.loading}
    <!-- Loading -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3">
      <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin"></div>
      <span class="text-xs text-text-subtlest">Sending request...</span>
      <button
        onclick={() => { response.loading = false; }}
        class="px-3 py-1 text-xs border border-border rounded-md
          text-text-subtle hover:text-text hover:border-border-focus
          transition-colors"
      >
        Cancel
      </button>
    </div>

  {:else if response.error && !response.data}
    <!-- Error only -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3 p-6">
      <svg class="w-8 h-8 text-danger/60" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"/>
      </svg>
      <span class="text-sm text-danger text-center max-w-[300px]">{response.error}</span>
    </div>

  {:else if response.data}
    <!-- Status bar -->
    <div class="flex items-center gap-2 px-3 py-1.5 border-b border-border-subtle
      font-mono text-xs shrink-0 hide-scrollbars overflow-x-auto">
      <!-- Status badge -->
      <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded font-bold
        {statusColor(response.data.status)} {statusBg(response.data.status)}">
        {response.data.status}
        <span class="font-normal text-text-subtle">{statusReason(response.data.status)}</span>
      </span>

      <span class="text-text-subtlest">•</span>
      <span class="text-text-subtle">{formatDuration(response.data.duration)}</span>
      <span class="text-text-subtlest">•</span>
      <span class="text-text-subtle">{formatBytes(response.data.size)}</span>

      {#if response.data.redirectCount > 0}
        <span class="text-warning text-[10px] font-sans">
          {response.data.redirectCount} redirect{response.data.redirectCount > 1 ? 's' : ''}
        </span>
      {/if}

      <!-- Pretty/Raw toggle -->
      {#if ui.activeResponseTab === 'body' && isJSON(response.data.contentType)}
        <div class="ml-auto flex items-center gap-0.5">
          <button
            onclick={() => { ui.responseViewMode = 'pretty'; }}
            class="px-2 py-0.5 text-[10px] font-medium rounded
              {ui.responseViewMode === 'pretty' ? 'bg-primary/15 text-primary' : 'text-text-subtlest hover:text-text-subtle'}"
          >Pretty</button>
          <button
            onclick={() => { ui.responseViewMode = 'raw'; }}
            class="px-2 py-0.5 text-[10px] font-medium rounded
              {ui.responseViewMode === 'raw' ? 'bg-primary/15 text-primary' : 'text-text-subtlest hover:text-text-subtle'}"
          >Raw</button>
        </div>
      {/if}
    </div>

    <!-- Error banner (if error + data) -->
    {#if response.error}
      <div class="mx-3 mt-1.5 px-3 py-2 rounded-md bg-danger/10 border border-danger/20 text-xs text-danger">
        {response.error}
      </div>
    {/if}

    <!-- Tabs -->
    <TabBar
      {tabs}
      active={ui.activeResponseTab}
      onchange={(id) => { ui.activeResponseTab = id as ResponseTab; }}
    />

    <!-- Tab content -->
    <div class="flex-1 overflow-hidden">
      {#if ui.activeResponseTab === 'body'}
        <ResponseBody
          body={response.data.body}
          contentType={response.data.contentType}
          pretty={ui.responseViewMode === 'pretty'}
        />
      {:else if ui.activeResponseTab === 'headers'}
        <ResponseHeaders
          requestHeaders={{}}
          responseHeaders={response.data.headers}
        />
      {:else if ui.activeResponseTab === 'timeline'}
        <ResponseTimeline
          dns={response.data.dnsTime}
          connect={response.data.connectTime}
          tls={response.data.tlsTime}
          ttfb={response.data.ttfbTime}
          total={response.data.duration}
        />
      {:else if ui.activeResponseTab === 'info'}
        <div class="p-4 flex flex-col gap-3">
          {#each [
            ['Protocol', response.data.protocol],
            ['Remote Address', response.data.remoteAddr || 'N/A'],
            ['Content-Type', response.data.contentType || 'N/A'],
            ['Response Size', formatBytes(response.data.size)],
            ['Total Duration', formatDuration(response.data.duration)],
          ] as [label, value] (label)}
            <div class="flex items-center justify-between py-1 border-b border-border-subtle/50">
              <span class="text-xs text-text-subtle">{label}</span>
              <span class="text-xs text-text font-mono">{value}</span>
            </div>
          {/each}
        </div>
      {/if}
    </div>

  {:else}
    <!-- Empty state -->
    <div class="flex-1 flex flex-col items-center justify-center gap-4 text-text-subtlest">
      <svg class="w-16 h-16 opacity-20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="0.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5"/>
      </svg>
      <div class="text-center">
        <p class="text-sm text-text-subtle">Enter a URL and send a request</p>
        <p class="text-[11px] text-text-subtlest mt-1">Ctrl+Enter to send</p>
      </div>
    </div>
  {/if}
</div>
