<!--
  RequestPane - Full HTTP request editing panel.
  Structure: CommandBar → Tabs (Body, Params, Headers, Auth, Settings, Info)
-->
<script lang="ts">
  import { request, response, ui } from '../stores/app.svelte';
  import type { HttpMethod, RequestTab } from '../lib/types';
  import CommandBar from './CommandBar.svelte';
  import TabBar from './core/TabBar.svelte';
  import PairEditor from './core/PairEditor.svelte';
  import BodyEditor from './BodyEditor.svelte';
  import AuthEditor from './AuthEditor.svelte';
  import SettingsEditor from './SettingsEditor.svelte';
  import ScriptEditor from './ScriptEditor.svelte';

  interface Props {
    onSend: () => void;
    onCancel: () => void;
  }

  let { onSend, onCancel }: Props = $props();

  const headerCount = $derived(request.headers.filter(h => h.name !== '').length);
  const paramCount = $derived(request.urlParameters.filter(p => p.name !== '').length);
  const authActive = $derived(request.auth.type !== 'none');

  const tabs = $derived([
    { id: 'params', label: 'Params', badge: paramCount || undefined },
    { id: 'auth', label: 'Authorization', badge: authActive ? '•' : undefined },
    { id: 'headers', label: 'Headers', badge: headerCount || undefined },
    { id: 'body', label: 'Body', badge: request.bodyType !== 'none' ? '•' : undefined },
    { id: 'scripts', label: 'Scripts' },
    { id: 'settings', label: 'Settings' },
  ]);
</script>

<div class="flex flex-col h-full bg-surface rounded-lg border border-border-subtle overflow-hidden">
  <!-- Command bar -->
  <CommandBar
    protocol="http"
    method={request.method}
    url={request.url}
    loading={response.loading}
    onMethodChange={(m) => { request.method = m as HttpMethod; }}
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
    {:else if ui.activeRequestTab === 'scripts'}
      <ScriptEditor />
    {:else if ui.activeRequestTab === 'settings'}
      <SettingsEditor />
    {/if}
  </div>
</div>
