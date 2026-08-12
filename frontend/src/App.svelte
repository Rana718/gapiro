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

<div class="flex flex-col w-full h-full overflow-hidden gpu">
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
    <span class="text-[10px] text-text-subtlest font-mono">Ctrl+Enter to send</span>
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
                <RequestPane onSend={sendRequest} onCancel={cancelRequest} />
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
