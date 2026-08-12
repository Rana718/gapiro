<!--
  RequestPane - The left/top panel containing URL bar and request configuration tabs.
-->
<script lang="ts">
  import { requestState, responseState, uiState } from '../stores/app.svelte';
  import type { HttpMethod, BodyType, RequestTab } from '../lib/types';
  import { emptyKV } from '../lib/utils';
  import UrlBar from './UrlBar.svelte';
  import TabBar from './TabBar.svelte';
  import KeyValueEditor from './KeyValueEditor.svelte';
  import CodeEditor from './CodeEditor.svelte';
  import BodyTypeSelector from './BodyTypeSelector.svelte';

  interface Props {
    onSend: () => void;
    onCancel: () => void;
  }

  let { onSend, onCancel }: Props = $props();

  const tabs = $derived([
    { id: 'params', label: 'Params', badge: requestState.queryParams.filter(p => p.key !== '').length || undefined },
    { id: 'headers', label: 'Headers', badge: requestState.headers.filter(h => h.key !== '').length || undefined },
    { id: 'body', label: 'Body' },
    { id: 'settings', label: 'Settings' },
  ]);
</script>

<div class="flex flex-col h-full bg-surface-1 rounded-lg border border-border-subtle gpu-layer">
  <!-- URL Bar -->
  <UrlBar
    method={requestState.method}
    url={requestState.url}
    loading={responseState.loading}
    onMethodChange={(m) => { requestState.method = m as HttpMethod; }}
    onUrlChange={(u) => { requestState.url = u; }}
    {onSend}
    {onCancel}
  />

  <!-- Tab Navigation -->
  <TabBar
    {tabs}
    active={uiState.activeRequestTab}
    onchange={(id) => { uiState.activeRequestTab = id as RequestTab; }}
  />

  <!-- Tab Content -->
  <div class="flex-1 overflow-auto">
    {#if uiState.activeRequestTab === 'params'}
      <KeyValueEditor
        items={requestState.queryParams}
        keyPlaceholder="Parameter"
        valuePlaceholder="Value"
        onchange={(items) => { requestState.queryParams = items; }}
      />
    {:else if uiState.activeRequestTab === 'headers'}
      <KeyValueEditor
        items={requestState.headers}
        keyPlaceholder="Header"
        valuePlaceholder="Value"
        onchange={(items) => { requestState.headers = items; }}
      />
    {:else if uiState.activeRequestTab === 'body'}
      <div class="flex flex-col h-full">
        <!-- Body type selector -->
        <BodyTypeSelector
          bodyType={requestState.bodyType}
          onchange={(bt) => { requestState.bodyType = bt as BodyType; }}
        />

        <!-- Body content -->
        {#if requestState.bodyType === 'none'}
          <div class="flex-1 flex items-center justify-center text-text-muted text-sm">
            This request does not have a body
          </div>
        {:else if requestState.bodyType === 'json'}
          <CodeEditor
            value={requestState.body}
            onchange={(v) => { requestState.body = v; }}
            language="json"
            placeholder={'{\n  "key": "value"\n}'}
          />
        {:else if requestState.bodyType === 'text'}
          <CodeEditor
            value={requestState.body}
            onchange={(v) => { requestState.body = v; }}
            language="text"
            placeholder="Enter request body..."
          />
        {:else if requestState.bodyType === 'form-urlencoded' || requestState.bodyType === 'form-data'}
          <KeyValueEditor
            items={requestState.formData}
            keyPlaceholder="Field"
            valuePlaceholder="Value"
            onchange={(items) => { requestState.formData = items; }}
          />
        {/if}
      </div>
    {:else if uiState.activeRequestTab === 'settings'}
      <div class="p-4 flex flex-col gap-3">
        <!-- Timeout -->
        <label class="flex items-center justify-between">
          <span class="text-xs text-text-secondary">Timeout (seconds)</span>
          <input
            type="number"
            min="0"
            max="300"
            value={requestState.timeout}
            oninput={(e) => { requestState.timeout = parseInt((e.target as HTMLInputElement).value) || 0; }}
            class="w-20 px-2 py-1 text-xs bg-surface-2 border border-border-default rounded
              text-text-primary focus:outline-none focus:border-accent"
          />
        </label>

        <!-- Follow redirects -->
        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-xs text-text-secondary">Follow redirects</span>
          <input
            type="checkbox"
            checked={requestState.followRedirects}
            onchange={() => { requestState.followRedirects = !requestState.followRedirects; }}
            class="w-4 h-4 rounded border-border-default bg-surface-2
              checked:bg-accent checked:border-accent cursor-pointer"
          />
        </label>

        <!-- Verify SSL -->
        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-xs text-text-secondary">Verify SSL certificate</span>
          <input
            type="checkbox"
            checked={requestState.verifySSL}
            onchange={() => { requestState.verifySSL = !requestState.verifySSL; }}
            class="w-4 h-4 rounded border-border-default bg-surface-2
              checked:bg-accent checked:border-accent cursor-pointer"
          />
        </label>
      </div>
    {/if}
  </div>
</div>
