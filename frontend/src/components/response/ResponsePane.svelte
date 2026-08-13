<!--
  ResponsePane - Response viewer matching Postman style.
  Tabs: Body | Cookies | Headers | Test Results
  Status bar with 200 OK, time, size.
  Body tab has JSON dropdown and Pretty/Preview/Visualize toggles.
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

  let activeTab = $state<'body' | 'cookies' | 'headers' | 'test-results'>('body');
  let viewMode = $state<'pretty' | 'raw' | 'preview' | 'visualize'>('pretty');
  let searchQuery = $state('');
  let showSearch = $state(false);

  const headerCount = $derived(() => {
    if (!response.data?.headers) return 0;
    return Object.keys(response.data.headers).length;
  });

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

  function saveResponse() {
    if (!response.data?.body) return;
    const blob = new Blob([response.data.body], { type: response.data.contentType || 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'response.' + (isJSON(response.data.contentType) ? 'json' : 'txt');
    a.click();
    URL.revokeObjectURL(url);
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
    <!-- Tabs + Status bar (Postman style: tabs left, status right) -->
    <div class="flex items-center justify-between border-b border-border-subtle shrink-0 px-3 h-10">
      <!-- Tabs -->
      <div class="flex items-center gap-0">
        {#each [
          { id: 'body', label: 'Body' },
          { id: 'cookies', label: 'Cookies' },
          { id: 'headers', label: `Headers`, badge: headerCount() },
          { id: 'test-results', label: 'Test Results' },
        ] as tab (tab.id)}
          <button
            type="button"
            onclick={() => { activeTab = tab.id as any; }}
            class="relative px-3 py-2 text-xs font-medium transition-colors
              {activeTab === tab.id ? 'text-text' : 'text-text-subtle hover:text-text'}"
          >
            {tab.label}
            {#if tab.badge}
              <span class="ml-1 text-[10px] text-text-subtlest">{tab.badge}</span>
            {/if}
            {#if activeTab === tab.id}
              <div class="absolute bottom-0 left-2 right-2 h-[2px] bg-primary rounded-t-full"></div>
            {/if}
          </button>
        {/each}
      </div>

      <!-- Status info -->
      <div class="flex items-center gap-2">
        <span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-bold
          {statusColor(response.data.status)} {statusBg(response.data.status)}">
          {response.data.status} {statusReason(response.data.status) || response.data.statusText || ''}
        </span>
        <span class="text-[11px] text-text-subtle font-mono">{formatDuration(response.data.duration)}</span>
        <span class="text-[11px] text-text-subtlest">·</span>
        <span class="text-[11px] text-text-subtle font-mono">{formatBytes(response.data.size)}</span>

        <!-- Save Response -->
        <button
          type="button"
          onclick={saveResponse}
          class="ml-2 flex items-center gap-1 px-2 py-1 text-[11px] text-text-subtle hover:text-text rounded hover:bg-surface-highlight transition-colors"
          title="Save Response"
        >
          <Icon name="download" size={12} />
          <span>Save Response</span>
        </button>
      </div>
    </div>

    <!-- Body sub-toolbar: format selector + view mode -->
    {#if activeTab === 'body'}
      <div class="flex items-center justify-between px-3 h-8 border-b border-border-subtle shrink-0 bg-surface-inset/50">
        <!-- Format dropdown -->
        <div class="flex items-center gap-2">
          {#if isJSON(response.data.contentType)}
            <div class="flex items-center gap-1 px-2 py-0.5 rounded border border-border-subtle bg-surface text-xs">
              <span class="text-text-subtle">{'{}'}</span>
              <span class="text-text font-medium">JSON</span>
              <Icon name="chevron-down" size={10} class="text-text-subtlest" />
            </div>
          {:else if isHTML(response.data.contentType)}
            <div class="flex items-center gap-1 px-2 py-0.5 rounded border border-border-subtle bg-surface text-xs">
              <span class="text-text font-medium">HTML</span>
            </div>
          {:else}
            <div class="flex items-center gap-1 px-2 py-0.5 rounded border border-border-subtle bg-surface text-xs">
              <span class="text-text font-medium">Text</span>
            </div>
          {/if}

          <!-- View mode toggles -->
          <div class="flex items-center gap-0 ml-2">
            <button type="button" onclick={() => { viewMode = 'pretty'; }}
              class="px-2 py-0.5 text-[11px] rounded-l border border-border-subtle
                {viewMode === 'pretty' ? 'bg-primary/15 text-primary font-medium border-primary/30' : 'text-text-subtle hover:text-text bg-surface'}">
              Pretty
            </button>
            <button type="button" onclick={() => { viewMode = 'raw'; }}
              class="px-2 py-0.5 text-[11px] border-y border-border-subtle
                {viewMode === 'raw' ? 'bg-primary/15 text-primary font-medium border-primary/30' : 'text-text-subtle hover:text-text bg-surface'}">
              Raw
            </button>
            <button type="button" onclick={() => { viewMode = 'preview'; }}
              class="px-2 py-0.5 text-[11px] border-y border-border-subtle
                {viewMode === 'preview' ? 'bg-primary/15 text-primary font-medium border-primary/30' : 'text-text-subtle hover:text-text bg-surface'}">
              Preview
            </button>
            <button type="button" onclick={() => { viewMode = 'visualize'; }}
              class="px-2 py-0.5 text-[11px] rounded-r border border-border-subtle
                {viewMode === 'visualize' ? 'bg-primary/15 text-primary font-medium border-primary/30' : 'text-text-subtle hover:text-text bg-surface'}">
              Visualize
            </button>
          </div>
        </div>

        <!-- Right side actions -->
        <div class="flex items-center gap-1">
          <button type="button" onclick={toggleSearch}
            class="p-1 rounded {showSearch ? 'bg-primary/15 text-primary' : 'text-text-subtle hover:text-text'}" title="Search (Ctrl+F)">
            <Icon name="search" size={13} />
          </button>
          <button type="button" onclick={copyBody}
            class="p-1 rounded text-text-subtle hover:text-text transition-colors" title="Copy">
            <Icon name="copy" size={13} />
          </button>
          {#if onClose}
            <button type="button" onclick={onClose}
              class="p-1 rounded text-text-subtle hover:text-text transition-colors" title="Close">
              <Icon name="x" size={13} />
            </button>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Search bar -->
    {#if showSearch && activeTab === 'body'}
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
      {:else if activeTab === 'cookies'}
        <div class="flex items-center justify-center h-full text-xs text-text-subtlest">
          No cookies
        </div>
      {:else if activeTab === 'test-results'}
        <div class="flex items-center justify-center h-full text-xs text-text-subtlest">
          No test results
        </div>
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
