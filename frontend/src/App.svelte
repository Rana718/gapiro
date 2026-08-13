<!--
  App.svelte - Root layout.
  Titlebar → TabStrip (open requests) → Sidebar | main area.
  The main area routes to the protocol editor for the active tab
  (HTTP / GraphQL / gRPC / WebSocket) and pairs it with the response panel
  (except WebSocket, which owns its full pane). "+" / Ctrl+N opens a new tab.
-->
<script lang="ts">
  import { ui, request, tabs, refreshDirty, closeTab } from './stores/app.svelte';
  import { execute, cancelRequest } from './lib/http';
  import { Toaster } from '@/components/ui/sonner';
  import SidebarLayout from './components/core/SidebarLayout.svelte';
  import SplitLayout from './components/core/SplitLayout.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import TabStrip from './components/TabStrip.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import GraphQLPane from './components/protocols/GraphQLPane.svelte';
  import GrpcPane from './components/protocols/GrpcPane.svelte';
  import WebSocketPane from './components/protocols/WebSocketPane.svelte';
  import ResponsePane from './components/response/ResponsePane.svelte';
  import NewTabPopup from './components/NewTabPopup.svelte';
  import Icon from './components/core/Icon.svelte';

  let showNewTabPopup = $state(false);
  let responsePanelPosition = $state<'bottom' | 'right'>('bottom');
  let responsePanelHidden = $state(false);

  const isWebSocket = $derived(request.protocol === 'websocket');

  // Keep each tab's dirty indicator in sync with live edits.
  $effect(() => { refreshDirty(); });

  /** Run the active request against its backend, revealing the response panel. */
  function runRequest() {
    if (!isWebSocket) responsePanelHidden = false;
    execute();
  }
  function handleKeydown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') { e.preventDefault(); runRequest(); }
    else if ((e.ctrlKey || e.metaKey) && e.key === 'n') { e.preventDefault(); showNewTabPopup = true; }
    else if ((e.ctrlKey || e.metaKey) && e.key === 'w') { e.preventDefault(); if (tabs.activeId) closeTab(tabs.activeId); }
    else if ((e.ctrlKey || e.metaKey) && e.key === 'b') { e.preventDefault(); ui.sidebarHidden = !ui.sidebarHidden; }
    else if (e.key === 'Escape' && showNewTabPopup) { showNewTabPopup = false; }
  }
</script>

<svelte:document onkeydown={handleKeydown} />

{#snippet requestEditor()}
  {#if request.protocol === 'graphql'}
    <GraphQLPane onSend={runRequest} onCancel={cancelRequest} />
  {:else if request.protocol === 'grpc'}
    <GrpcPane onSend={runRequest} onCancel={cancelRequest} />
  {:else}
    <RequestPane onSend={runRequest} onCancel={cancelRequest} />
  {/if}
{/snippet}

<div class="flex flex-col w-full h-full overflow-hidden bg-surface-inset text-text">
  <!-- Title bar -->
  <div
    class="flex items-center h-10 px-2.5 bg-surface-inset border-b border-border-subtle shrink-0 gap-1"
    style="--wails-draggable: drag"
  >
    <button
      type="button"
      onclick={() => { ui.sidebarHidden = !ui.sidebarHidden; }}
      class="flex items-center justify-center size-7 rounded-md text-text-subtle hover:text-text hover:bg-surface-highlight transition-colors"
      title="Toggle sidebar (Ctrl+B)"
    >
      <Icon name={ui.sidebarHidden ? 'layout-sidebar-left-expand' : 'layout-sidebar-left-collapse'} size={16} />
    </button>

    <span class="text-xs font-semibold tracking-tight text-text-subtle select-none px-1">Gapiro</span>

    <div class="flex-1"></div>

    {#if !isWebSocket}
      <button
        type="button"
        onclick={() => { responsePanelPosition = responsePanelPosition === 'bottom' ? 'right' : 'bottom'; }}
        disabled={responsePanelHidden}
        class="flex items-center justify-center size-7 rounded-md text-text-subtle hover:text-text hover:bg-surface-highlight disabled:opacity-40 transition-colors"
        title="Response panel: {responsePanelPosition === 'bottom' ? 'move right' : 'move bottom'}"
      >
        <Icon name={responsePanelPosition === 'bottom' ? 'layout-sidebar-right' : 'layout-bottombar'} size={16} />
      </button>
      <button
        type="button"
        onclick={() => { responsePanelHidden = !responsePanelHidden; }}
        class="flex items-center justify-center size-7 rounded-md transition-colors
          {responsePanelHidden ? 'text-text-subtle hover:text-text hover:bg-surface-highlight' : 'text-primary bg-primary/10 hover:bg-primary/20'}"
        title={responsePanelHidden ? 'Show response panel' : 'Hide response panel'}
      >
        <Icon name="layout-bottombar" size={16} />
      </button>
    {/if}
  </div>

  <!-- Tab strip -->
  <TabStrip onNew={() => { showNewTabPopup = true; }} />

  <!-- Main content -->
  <div class="flex-1 overflow-hidden">
    <SidebarLayout
      width={ui.sidebarWidth}
      onWidthChange={(w) => { ui.sidebarWidth = w; }}
      hidden={ui.sidebarHidden}
      onHiddenChange={(h) => { ui.sidebarHidden = h; }}
    >
      {#snippet sidebar()}
        <Sidebar onNew={() => { showNewTabPopup = true; }} />
      {/snippet}

      {#snippet children()}
        <div class="h-full p-1">
          {#if isWebSocket}
            {#key request.id}
              <WebSocketPane />
            {/key}
          {:else if responsePanelHidden}
            {@render requestEditor()}
          {:else}
            <SplitLayout
              layout={responsePanelPosition === 'bottom' ? 'vertical' : 'horizontal'}
              storageKey="request-response"
              defaultRatio={0.5}
              minPx={120}
            >
              {#snippet first()}
                <div class="h-full {responsePanelPosition === 'bottom' ? 'pb-0.5' : 'pr-0.5'}">
                  {@render requestEditor()}
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

{#if showNewTabPopup}
  <NewTabPopup onclose={() => { showNewTabPopup = false; }} />
{/if}

<Toaster position="bottom-right" />
