<!--
  ResponsePane - Response viewer with tabs, search, copy, HTML preview.
  Closeable. Shows syntax-highlighted body, headers, timing.
-->
<script lang="ts">
  import { response } from '../../stores/app.svelte';
  import { cancelRequest } from '../../lib/http';
  import { formatBytes, formatDuration, statusColor, statusBg, statusReason, isJSON, isHTML } from '../../lib/utils';
  import Icon from '../core/Icon.svelte';
  import { toast } from 'svelte-sonner';
  import ResponseBody from './ResponseBody.svelte';
  import ResponseHeaders from './ResponseHeaders.svelte';
  import ResponseTimeline from './ResponseTimeline.svelte';

  interface Props {
    onClose?: () => void;
    position?: 'bottom' | 'right';
  }

  let { onClose, position = 'bottom' }: Props = $props();

  let activeTab = $state<'body' | 'headers' | 'timeline'>('body');
  let viewMode = $state<'pretty' | 'raw' | 'preview'>('pretty');
  let searchQuery = $state('');
  let showSearch = $state(false);

  function copyBody() {
    if (response.data?.body) {
      navigator.clipboard.writeText(response.data.body);
      toast.success('Response body copied');
    }
  }

  function toggleSearch() {
    showSearch = !showSearch;
    if (!showSearch) searchQuery = '';
  }
</script>

<div class="flex flex-col h-full bg-surface rounded-lg border border-border-subtle overflow-hidden">
  {#if response.loading}
    <!-- Loading -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3">
      <div class="size-5 border-2 border-primary/30 border-t-primary rounded-full spinner"></div>
      <span class="text-xs text-text-subtle">Sending request…</span>
      <button
        type="button"
        onclick={cancelRequest}
        class="px-3 py-1 text-xs border border-border rounded-md
          text-text-subtle hover:text-text hover:border-border-focus transition-colors"
      >
        Cancel
      </button>
    </div>

  {:else if response.error && !response.data}
    <!-- Error -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3 p-6">
      <Icon name="alert-circle" size={28} class="text-danger" />
      <div class="text-danger text-sm text-center max-w-[320px] break-words">{response.error}</div>
    </div>

  {:else if response.data}
    <!-- Status bar -->
    <div class="flex items-center gap-2 px-3 h-9 border-b border-border-subtle shrink-0">
      <!-- Status -->
      <span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-bold
        {statusColor(response.data.status)} {statusBg(response.data.status)}">
        {response.data.status || '—'}
      </span>
      {#if statusReason(response.data.status)}
        <span class="text-[11px] text-text-subtle">{statusReason(response.data.status)}</span>
      {:else if response.data.statusText}
        <span class="text-[11px] text-text-subtle">{response.data.statusText}</span>
      {/if}
      <span class="text-[11px] text-text-subtlest">•</span>
      <span class="text-[11px] text-text-subtle font-mono">{formatDuration(response.data.duration)}</span>
      <span class="text-[11px] text-text-subtlest">•</span>
      <span class="text-[11px] text-text-subtle font-mono">{formatBytes(response.data.size)}</span>
      {#if response.data.protocol}
        <span class="text-[11px] text-text-subtlest">•</span>
        <span class="text-[11px] text-text-subtlest font-mono">{response.data.protocol}</span>
      {/if}

      <div class="flex-1"></div>

      <!-- Actions -->
      <div class="flex items-center gap-0.5">
        <!-- View mode for body tab -->
        {#if activeTab === 'body'}
          {#if isJSON(response.data.contentType)}
            <button type="button" onclick={() => { viewMode = 'pretty'; }} class="px-1.5 py-0.5 text-[10px] rounded {viewMode === 'pretty' ? 'bg-primary/15 text-primary' : 'text-text-subtle hover:text-text'}" title="Pretty">Pretty</button>
            <button type="button" onclick={() => { viewMode = 'raw'; }} class="px-1.5 py-0.5 text-[10px] rounded {viewMode === 'raw' ? 'bg-primary/15 text-primary' : 'text-text-subtle hover:text-text'}" title="Raw">Raw</button>
          {/if}
          {#if isHTML(response.data.contentType)}
            <button type="button" onclick={() => { viewMode = viewMode === 'preview' ? 'pretty' : 'preview'; }} class="p-1 rounded {viewMode === 'preview' ? 'bg-primary/15 text-primary' : 'text-text-subtle hover:text-text'}" title="Preview HTML">
              <Icon name="eye" size={14} />
            </button>
          {/if}
        {/if}

        <!-- Search -->
        <button type="button" onclick={toggleSearch} class="p-1 rounded {showSearch ? 'bg-primary/15 text-primary' : 'text-text-subtle hover:text-text'}" title="Search (Ctrl+F)">
          <Icon name="search" size={14} />
        </button>

        <!-- Copy -->
        <button type="button" onclick={copyBody} class="p-1 rounded text-text-subtle hover:text-text transition-colors" title="Copy response body">
          <Icon name="copy" size={14} />
        </button>

        <!-- Close panel -->
        {#if onClose}
          <button type="button" onclick={onClose} class="p-1 rounded text-text-subtle hover:text-text transition-colors" title="Hide response panel">
            <Icon name="x" size={14} />
          </button>
        {/if}
      </div>
    </div>

    <!-- Search bar -->
    {#if showSearch}
      <div class="flex items-center px-3 py-1.5 border-b border-border-subtle gap-2">
        <Icon name="search" size={12} class="text-text-subtlest" />
        <!-- svelte-ignore a11y_autofocus -->
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Search response…"
          class="flex-1 bg-transparent text-xs text-text placeholder:text-placeholder border-0 focus:outline-none"
          autofocus
        />
        {#if searchQuery}
          <button type="button" onclick={() => { searchQuery = ''; }} class="text-text-subtle hover:text-text">
            <Icon name="x" size={12} />
          </button>
        {/if}
      </div>
    {/if}

    <!-- Tabs -->
    <div class="flex items-center gap-0 border-b border-border-subtle px-3 shrink-0">
      {#each [['body', 'Body'], ['headers', 'Headers'], ['timeline', 'Timeline']] as [id, label] (id)}
        <button
          type="button"
          onclick={() => { activeTab = id as any; }}
          class="relative px-3 py-1.5 text-xs font-medium transition-colors
            {activeTab === id ? 'text-text' : 'text-text-subtle hover:text-text'}"
        >
          {label}
          {#if activeTab === id}
            <div class="absolute bottom-0 left-2 right-2 h-[2px] bg-primary rounded-t-full"></div>
          {/if}
        </button>
      {/each}
    </div>

    <!-- Tab content -->
    <div class="flex-1 overflow-hidden">
      {#if activeTab === 'body'}
        {#if viewMode === 'preview' && isHTML(response.data.contentType)}
          <iframe
            srcdoc={response.data.body}
            class="w-full h-full border-0 bg-white"
            sandbox="allow-same-origin"
            title="HTML Preview"
          ></iframe>
        {:else}
          <ResponseBody
            body={response.data.body}
            contentType={response.data.contentType}
            pretty={viewMode === 'pretty'}
            {searchQuery}
          />
        {/if}
      {:else if activeTab === 'headers'}
        <ResponseHeaders
          requestHeaders={{}}
          responseHeaders={response.data.headers}
        />
      {:else if activeTab === 'timeline'}
        <ResponseTimeline
          dns={response.data.dnsTime}
          connect={response.data.connectTime}
          tls={response.data.tlsTime}
          ttfb={response.data.ttfbTime}
          total={response.data.duration}
        />
      {/if}
    </div>

  {:else}
    <!-- Empty state -->
    <div class="flex-1 flex flex-col items-center justify-center gap-2 text-center select-none">
      <Icon name="send" size={30} class="text-text-subtlest/40" />
      <p class="text-sm text-text-subtle">Send a request to see the response</p>
      <p class="text-[11px] text-text-subtlest">Ctrl+Enter to send</p>
    </div>
  {/if}
</div>
