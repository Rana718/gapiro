<!--
  Sidebar - Request collection navigation with folders and requests.
  Matches Yaak's sidebar structure.
-->
<script lang="ts">
  import { collection, request, loadRequest, resetRequest, saveCurrentRequest, deleteRequest, duplicateRequest, createFolder } from '../stores/app.svelte';
  import { methodColor } from '../lib/utils';
  import type { RequestConfig } from '../lib/types';

  let filterText = $state('');
  let showNewFolder = $state(false);
  let newFolderName = $state('');
  let contextMenu = $state<{ x: number; y: number; requestId: string } | null>(null);

  const filteredRequests = $derived(
    filterText
      ? collection.requests.filter(r =>
          r.name.toLowerCase().includes(filterText.toLowerCase()) ||
          r.url.toLowerCase().includes(filterText.toLowerCase()))
      : collection.requests
  );

  function handleNew() {
    resetRequest();
  }

  function handleSelect(req: RequestConfig) {
    loadRequest(req);
  }

  function handleSave() {
    saveCurrentRequest();
  }

  function handleCreateFolder() {
    if (newFolderName.trim()) {
      createFolder(newFolderName.trim());
      newFolderName = '';
      showNewFolder = false;
    }
  }

  function handleContextMenu(e: MouseEvent, reqId: string) {
    e.preventDefault();
    contextMenu = { x: e.clientX, y: e.clientY, requestId: reqId };
  }

  function closeContextMenu() {
    contextMenu = null;
  }
</script>

<svelte:document onclick={closeContextMenu} />

<div class="flex flex-col h-full bg-surface-inset">
  <!-- Header -->
  <div class="flex items-center justify-between px-3 py-2 border-b border-border-subtle shrink-0">
    <span class="text-[10px] font-semibold text-text-subtlest uppercase tracking-wider">Collection</span>
    <div class="flex items-center gap-0.5">
      <button
        onclick={() => { showNewFolder = !showNewFolder; }}
        class="p-1 rounded text-text-subtlest hover:text-text-subtle hover:bg-surface-highlight
          transition-colors"
        title="New Folder"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 10v4m-2-2h4m-7 6h10a2 2 0 002-2V8l-4-4H7a2 2 0 00-2 2v10a2 2 0 002 2z"/>
        </svg>
      </button>
      <button
        onclick={handleNew}
        class="p-1 rounded text-text-subtlest hover:text-text-subtle hover:bg-surface-highlight
          transition-colors"
        title="New Request"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
        </svg>
      </button>
    </div>
  </div>

  <!-- Search/Filter -->
  <div class="px-2 py-1.5 border-b border-border-subtle shrink-0">
    <input
      type="text"
      bind:value={filterText}
      placeholder="Filter..."
      class="w-full px-2.5 py-1 text-xs bg-surface border border-border-subtle rounded-md
        text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
    />
  </div>

  <!-- New folder input -->
  {#if showNewFolder}
    <div class="px-2 py-1.5 border-b border-border-subtle">
      <form onsubmit={(e) => { e.preventDefault(); handleCreateFolder(); }}>
        <input
          type="text"
          bind:value={newFolderName}
          placeholder="Folder name..."
          class="w-full px-2.5 py-1 text-xs bg-surface border border-border-subtle rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </form>
    </div>
  {/if}

  <!-- Request list -->
  <div class="flex-1 overflow-y-auto py-1">
    {#if collection.folders.length === 0 && filteredRequests.length === 0}
      <div class="flex flex-col items-center justify-center h-full gap-2 p-4 text-center">
        <svg class="w-10 h-10 text-text-subtlest/30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="0.8">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/>
        </svg>
        <p class="text-[11px] text-text-subtlest">No saved requests</p>
        <p class="text-[10px] text-text-subtlest/60">Send a request and save it here</p>
      </div>
    {:else}
      <!-- Folders -->
      {#each collection.folders as folder (folder.id)}
        <div class="px-1">
          <button
            onclick={() => { folder.expanded = !folder.expanded; }}
            class="w-full flex items-center gap-1.5 px-2 py-1 rounded text-left
              text-xs text-text-subtle hover:bg-surface-highlight transition-colors"
          >
            <svg class="w-3 h-3 transition-transform duration-100 {folder.expanded ? 'rotate-90' : ''}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
            </svg>
            <svg class="w-3.5 h-3.5 text-text-subtlest" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V7z"/>
            </svg>
            <span class="truncate">{folder.name}</span>
          </button>
        </div>
      {/each}

      <!-- Requests -->
      {#each filteredRequests as req (req.id)}
        <div class="px-1">
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            onclick={() => handleSelect(req)}
            oncontextmenu={(e) => handleContextMenu(e, req.id)}
            onkeydown={(e) => { if (e.key === 'Enter') handleSelect(req); }}
            role="button"
            tabindex="0"
            class="w-full flex items-center gap-2 px-2 py-1.5 rounded group
              transition-colors duration-50
              {collection.activeRequestId === req.id
                ? 'bg-primary/8 border border-primary/20'
                : 'hover:bg-surface-highlight border border-transparent'}"
          >
            <span class="text-[9px] font-bold w-[34px] shrink-0 text-right font-mono {methodColor(req.method)}">
              {req.method}
            </span>
            <span class="text-xs text-text truncate flex-1">{req.name || req.url || 'Untitled'}</span>
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <!-- Save button -->
  {#if request.url}
    <div class="px-2 py-2 border-t border-border-subtle shrink-0">
      <button
        onclick={handleSave}
        class="w-full flex items-center justify-center gap-1.5 px-3 py-1.5
          text-xs font-medium text-primary bg-primary/8 rounded-md
          hover:bg-primary/15 transition-colors"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V7l-4-4z"/>
          <path stroke-linecap="round" stroke-linejoin="round" d="M7 3v5h8V3M7 14h10"/>
        </svg>
        Save
      </button>
    </div>
  {/if}

  <!-- Context menu -->
  {#if contextMenu}
    <div
      class="fixed z-50 py-1 min-w-[140px] bg-surface-active border border-border rounded-lg shadow-lg"
      style="left: {contextMenu.x}px; top: {contextMenu.y}px;"
    >
      <button
        onclick={() => { duplicateRequest(contextMenu!.requestId); closeContextMenu(); }}
        class="w-full px-3 py-1.5 text-xs text-left text-text hover:bg-surface-highlight transition-colors"
      >Duplicate</button>
      <button
        onclick={() => { deleteRequest(contextMenu!.requestId); closeContextMenu(); }}
        class="w-full px-3 py-1.5 text-xs text-left text-danger hover:bg-danger/10 transition-colors"
      >Delete</button>
    </div>
  {/if}
</div>
