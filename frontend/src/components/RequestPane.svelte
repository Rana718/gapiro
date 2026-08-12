<!--
  RequestPane - Full request editing panel.
  Structure: UrlBar → Tabs (Body, Params, Headers, Auth, Settings, Info)
  Matches Yaak's HttpRequestPane.
-->
<script lang="ts">
  import { request, response, ui } from '../stores/app.svelte';
  import type { HttpMethod, BodyType, RequestTab, Pair } from '../lib/types';
  import { BODY_TYPES, emptyPair } from '../lib/utils';
  import UrlBar from './UrlBar.svelte';
  import TabBar from './core/TabBar.svelte';
  import PairEditor from './core/PairEditor.svelte';
  import BodyEditor from './BodyEditor.svelte';
  import AuthEditor from './AuthEditor.svelte';
  import SettingsEditor from './SettingsEditor.svelte';

  interface Props {
    onSend: () => void;
    onCancel: () => void;
  }

  let { onSend, onCancel }: Props = $props();

  // Count non-empty items for badges
  const headerCount = $derived(request.headers.filter(h => h.name !== '').length);
  const paramCount = $derived(request.urlParameters.filter(p => p.name !== '').length);
  const bodyCount = $derived(() => {
    if (request.bodyType === 'form-urlencoded' || request.bodyType === 'form-data') {
      return request.formData.filter(f => f.name !== '').length;
    }
    return request.body ? 1 : 0;
  });

  const tabs = $derived([
    { id: 'body', label: request.bodyType === 'none' ? 'Body' : BODY_TYPES.find(b => b.id === request.bodyType)?.short ?? 'Body', badge: bodyCount() || undefined },
    { id: 'params', label: 'Params', badge: paramCount || undefined },
    { id: 'headers', label: 'Headers', badge: headerCount || undefined },
    { id: 'auth', label: 'Auth' },
    { id: 'settings', label: 'Settings' },
    { id: 'description', label: 'Info' },
  ]);
</script>

<div class="flex flex-col h-full bg-surface rounded-md border border-border-subtle overflow-hidden gpu">
  <!-- URL Bar -->
  <UrlBar
    method={request.method}
    url={request.url}
    loading={response.loading}
    onMethodChange={(m) => { request.method = m; }}
    onUrlChange={(u) => { request.url = u; }}
    {onSend}
    {onCancel}
  />

  <!-- Tabs -->
  <TabBar
    {tabs}
    active={ui.activeRequestTab}
    onchange={(id) => { ui.activeRequestTab = id as RequestTab; }}
  />

  <!-- Tab content -->
  <div class="flex-1 overflow-hidden">
    {#if ui.activeRequestTab === 'body'}
      <BodyEditor />
    {:else if ui.activeRequestTab === 'params'}
      <PairEditor
        pairs={request.urlParameters}
        namePlaceholder="Parameter"
        valuePlaceholder="Value"
        onchange={(p) => { request.urlParameters = p; }}
      />
    {:else if ui.activeRequestTab === 'headers'}
      <PairEditor
        pairs={request.headers}
        namePlaceholder="Header"
        valuePlaceholder="Value"
        onchange={(p) => { request.headers = p; }}
      />
    {:else if ui.activeRequestTab === 'auth'}
      <AuthEditor />
    {:else if ui.activeRequestTab === 'settings'}
      <SettingsEditor />
    {:else if ui.activeRequestTab === 'description'}
      <div class="p-3 flex flex-col gap-3 h-full">
        <input
          type="text"
          value={request.name}
          oninput={(e) => { request.name = (e.target as HTMLInputElement).value; }}
          placeholder="Request Name"
          class="w-full bg-transparent text-lg font-semibold text-text
            placeholder:text-placeholder border-0 focus:outline-none px-0"
        />
        <textarea
          value={request.description}
          oninput={(e) => { request.description = (e.target as HTMLTextAreaElement).value; }}
          placeholder="Add a description for this request..."
          class="flex-1 w-full bg-surface-inset text-xs text-text font-mono
            p-3 rounded-md border border-border-subtle resize-none
            placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        ></textarea>
      </div>
    {/if}
  </div>
</div>
