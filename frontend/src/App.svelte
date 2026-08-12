<!--
  App.svelte - Root layout.
  - No protocol tabs at top. User clicks "+" to create new tab (HTTP/gRPC/GraphQL/WebSocket)
  - Left panel (sidebar) is closeable via Ctrl+B
  - Bottom panel (response) is closeable and can be moved to right side
  - Uses @tabler/icons-svelte for icons
  - Uses shadcn-svelte Sonner for toasts
-->
<script lang="ts">
  import { ui, collection, resetRequest } from './stores/app.svelte';
  import { sendRequest, cancelRequest } from './lib/http';
  import { Toaster } from '@/components/ui/sonner';
  import SidebarLayout from './components/core/SidebarLayout.svelte';
  import SplitLayout from './components/core/SplitLayout.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import ResponsePane from './components/response/ResponsePane.svelte';
  import NewTabPopup from './components/NewTabPopup.svelte';
  import Icon from './components/core/Icon.svelte';

  let showNewTabPopup = $state(false);
  let responsePanelPosition = $state<'bottom' | 'right'>('bottom');
  let responsePanelHidden = $state(false);

  // Keyboard shortcuts
  function handleKeydown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      sendRequest();
    }
    if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
      e.preventDefault();
      showNewTabPopup = true;
    }
    if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
      e.preventDefault();
      ui.sidebarHidden = !ui.sidebarHidden;
    }
    if (e.key === 'Escape') {
      showNewTabPopup = false;
    }
  }
</script>

<svelte:document onkeydown={handleKeydown} />

<div class="flex flex-col w-full h-full overflow-hidden bg-background text-foreground">
  <!-- Title bar -->
  <div
    class="flex items-center h-10 px-3 bg-sidebar border-b border-sidebar-border shrink-0 gap-2"
    style="--wails-draggable: drag"
  >
    <!-- Sidebar toggle -->
    <button
      type="button"
      onclick={() => { ui.sidebarHidden = !ui.sidebarHidden; }}
      class="p-1 rounded hover:bg-sidebar-accent text-muted-foreground hover:text-foreground"
      title="Toggle sidebar (Ctrl+B)"
    >
      <Icon name="layout-sidebar-left-collapse" size={16} />
    </button>

    <!-- App title -->
    <span class="text-xs font-semibold text-muted-foreground select-none">Gapiro</span>

    <div class="flex-1"></div>

    <!-- New tab button -->
    <button
      type="button"
      onclick={() => { showNewTabPopup = true; }}
      class="flex items-center gap-1 px-2 py-1 rounded text-xs font-medium
        text-muted-foreground hover:text-foreground hover:bg-sidebar-accent"
      title="New request (Ctrl+N)"
    >
      <Icon name="plus" size={14} />
      New
    </button>

    <!-- Response panel controls -->
    <button
      type="button"
      onclick={() => { responsePanelPosition = responsePanelPosition === 'bottom' ? 'right' : 'bottom'; }}
      class="p-1 rounded hover:bg-sidebar-accent text-muted-foreground hover:text-foreground"
      title="Toggle response panel position"
    >
      {#if responsePanelPosition === 'bottom'}
        <Icon name="layout-sidebar-right" size={16} />
      {:else}
        <Icon name="layout-bottombar" size={16} />
      {/if}
    </button>
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
        <div class="h-full p-1">
          {#if responsePanelHidden}
            <!-- Only request pane -->
            <RequestPane onSend={sendRequest} onCancel={cancelRequest} />
          {:else}
            <SplitLayout
              layout={responsePanelPosition === 'bottom' ? 'vertical' : 'horizontal'}
              storageKey="request-response"
              defaultRatio={0.5}
              minPx={120}
            >
              {#snippet first()}
                <div class="h-full {responsePanelPosition === 'bottom' ? 'pb-0.5' : 'pr-0.5'}">
                  <RequestPane onSend={sendRequest} onCancel={cancelRequest} />
                </div>
              {/snippet}

              {#snippet second()}
                <div class="h-full {responsePanelPosition === 'bottom' ? 'pt-0.5' : 'pl-0.5'}">
                  <ResponsePane
                    onClose={() => { responsePanelHidden = true; }}
                    position={responsePanelPosition}
                  />
                </div>
              {/snippet}
            </SplitLayout>
          {/if}
        </div>
      {/snippet}
    </SidebarLayout>
  </div>
</div>

<!-- New tab popup -->
{#if showNewTabPopup}
  <NewTabPopup onclose={() => { showNewTabPopup = false; }} />
{/if}

<!-- Toast notifications -->
<Toaster position="bottom-right" />
