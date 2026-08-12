<!--
  App.svelte - Root application layout.
  Structure: Sidebar | (RequestPane / ResponsePane) split layout.
  GPU-composited, optimized for Wayland/Hyprland.
-->
<script lang="ts">
  import { ui } from './stores/app.svelte';
  import { sendRequest, cancelRequest } from './lib/http';
  import SidebarLayout from './components/core/SidebarLayout.svelte';
  import SplitLayout from './components/core/SplitLayout.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import ResponsePane from './components/response/ResponsePane.svelte';
  import ProtocolPane from './components/ProtocolPane.svelte';
  let overlay = $state<'none'|'settings'|'environments'|'shortcuts'>('none');
  let appearance = $state<'dark' | 'light'>('dark');
  function toggleAppearance() {
    appearance = appearance === 'dark' ? 'light' : 'dark';
    document.documentElement.dataset.appearance = appearance;
    localStorage.setItem('gapiro:appearance', appearance);
  }
  if (typeof localStorage !== 'undefined') {
    appearance = (localStorage.getItem('gapiro:appearance') as 'dark' | 'light') || 'dark';
  }

  // Keyboard shortcuts
  function handleKeydown(e: KeyboardEvent) {
    // Ctrl+Enter: Send request
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      sendRequest();
    }
    // Ctrl+N: New request
    if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
      e.preventDefault();
      import('./stores/app.svelte').then(m => m.resetRequest());
    }
    // Ctrl+B: Toggle sidebar
    if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
      e.preventDefault();
      ui.sidebarHidden = !ui.sidebarHidden;
    }
  }
</script>

<svelte:document onkeydown={handleKeydown} />

<div class="flex flex-col w-full h-full overflow-hidden gpu" data-appearance={appearance}>
  <!-- Title bar / drag area -->
  <div
    class="flex items-center h-9 px-3 bg-surface-inset border-b border-border-subtle shrink-0"
    style="--wails-draggable: drag"
  >
    <div class="flex items-center gap-2">
      <svg class="w-4 h-4 text-primary" viewBox="0 0 24 24" fill="currentColor">
        <path d="M13 3L4 14h7l-2 7 9-11h-7l2-7z"/>
      </svg>
      <span class="text-[11px] font-semibold text-text-subtle">Gapiro</span>
    </div>
    <div class="flex-1"></div>
    <div class="flex items-center gap-3">
      <span class="text-[10px] text-text-subtlest font-mono">Ctrl+Enter to send</span>
      <button onclick={toggleAppearance} class="h-6 px-2 rounded text-[10px] text-text-subtle hover:text-text hover:bg-surface-highlight" title="Toggle appearance">
        {appearance === 'dark' ? '☼ Light' : '☾ Dark'}
      </button>
      <button onclick={()=>overlay='environments'} class="h-6 px-2 rounded text-[10px] text-text-subtle hover:bg-surface-highlight">Environment ▾</button>
      <button onclick={()=>overlay='settings'} class="h-6 px-2 rounded text-[10px] text-text-subtle hover:bg-surface-highlight">⚙</button>
    </div>
  </div>

  <!-- Main content -->
  <div class="flex-1 overflow-hidden">
    <SidebarLayout
      width={ui.sidebarWidth}
      onWidthChange={(w) => { ui.sidebarWidth = w; }}
      hidden={ui.sidebarHidden}
      onHiddenChange={(h) => { ui.sidebarHidden = h; }}
    >
      {#snippet sidebar()}
        <Sidebar />
      {/snippet}

      {#snippet children()}
        <div class="h-full p-1.5">
          <SplitLayout
            layout="vertical"
            storageKey="request-response"
            defaultRatio={0.5}
            minPx={150}
          >
            {#snippet first()}
              <div class="h-full pb-0.5">
                <ProtocolPane><RequestPane onSend={sendRequest} onCancel={cancelRequest} /></ProtocolPane>
              </div>
            {/snippet}

            {#snippet second()}
              <div class="h-full pt-0.5">
                <ResponsePane />
              </div>
            {/snippet}
          </SplitLayout>
        </div>
      {/snippet}
    </SidebarLayout>
  </div>
</div>
{#if overlay !== 'none'}
  <div class="fixed inset-0 z-[100] bg-black/60 flex items-center justify-center" onclick={()=>overlay='none'}>
    <div class="w-[520px] max-w-[90vw] rounded-xl border border-border bg-surface shadow-lg p-5" onclick={(e)=>e.stopPropagation()}>
      <div class="flex items-center justify-between mb-4"><h2 class="text-base font-semibold">{overlay==='settings'?'Settings':overlay==='environments'?'Environments':'Keyboard shortcuts'}</h2><button onclick={()=>overlay='none'} class="text-text-subtlest hover:text-text">✕</button></div>
      {#if overlay==='settings'}<div class="space-y-3"><label class="flex justify-between text-sm">Theme <span class="text-text-subtle">{appearance}</span></label><label class="flex justify-between text-sm">Request timeout <input value="30" class="w-20 px-2 bg-surface-inset border border-border rounded"/></label><label class="flex justify-between text-sm">Telemetry <input type="checkbox" checked/></label></div>
      {:else if overlay==='environments'}<div class="space-y-2"><div class="p-3 rounded-md bg-surface-highlight border border-primary/30"><div class="text-sm font-medium">Local</div><div class="text-xs text-text-subtlest">http://localhost:3000 · 4 variables</div></div><button class="w-full p-3 rounded-md border border-dashed border-border text-sm text-text-subtle">+ New environment</button></div>
      {:else}<div class="grid grid-cols-2 gap-2 text-xs text-text-subtle">{#each [['Ctrl+Enter','Send request'],['Ctrl+N','New request'],['Ctrl+B','Toggle sidebar'],['Ctrl+K','Command palette']] as s}<div class="flex justify-between p-2 rounded bg-surface-inset"><kbd>{s[0]}</kbd><span>{s[1]}</span></div>{/each}</div>{/if}
    </div>
  </div>
{/if}
