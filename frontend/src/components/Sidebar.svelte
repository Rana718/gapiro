<!--
  Sidebar - Request collection navigation with folders and saved requests.
-->
<script lang="ts">
  import { collectionState, requestState, loadRequest, resetRequest, saveCurrentRequest, createFolder, deleteRequest } from '../stores/app.svelte';
  import { methodColor } from '../lib/utils';
  import type { SavedRequest } from '../lib/types';

  let newFolderName = $state('');
  let showNewFolder = $state(false);

  function handleNewRequest() {
    resetRequest();
    collectionState.activeRequestId = null;
  }

  function handleSaveRequest() {
    const name = requestState.url || 'Untitled Request';
    saveCurrentRequest(name);
  }

  function handleSelectRequest(req: SavedRequest) {
    loadRequest(req);
  }

  function handleCreateFolder() {
    if (newFolderName.trim()) {
      createFolder(newFolderName.trim());
      newFolderName = '';
      showNewFolder = false;
    }
  }
</script>

<div class="flex flex-col h-full bg-surface-0 gpu-layer">
  <!-- Sidebar Header -->
  <div class="flex items-center justify-between px-3 py-2 border-b border-border-subtle">
    <span class="text-xs font-semibold text-text-secondary uppercase tracking-wider">Collection</span>
    <div class="flex items-center gap-1">
      <!-- New folder button -->
      <button
        onclick={() => { showNewFolder = !showNewFolder; }}
        class="p-1 rounded text-text-muted hover:text-text-secondary hover:bg-surface-2
          transition-colors cursor-pointer"
        aria-label="New folder"
        title="New folder"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M2 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V7z"/>
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 10v4m-2-2h4"/>
        </svg>
      </button>
      <!-- New request button -->
      <button
        onclick={handleNewRequest}
        class="p-1 rounded text-text-muted hover:text-text-secondary hover:bg-surface-2
          transition-colors cursor-pointer"
        aria-label="New request"
        title="New request"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
        </svg>
      </button>
    </div>
  </div>

  <!-- New folder input -->
  {#if showNewFolder}
    <div class="px-2 py-1.5 border-b border-border-subtle">
      <form onsubmit={(e) => { e.preventDefault(); handleCreateFolder(); }}>
        <input
          type="text"
          bind:value={newFolderName}
          placeholder="Folder name..."
          class="w-full px-2 py-1 text-xs bg-surface-2 border border-border-default rounded
            text-text-primary placeholder:text-text-muted/50 focus:outline-none focus:border-accent"
          autofocus
        />
      </form>
    </div>
  {/if}

  <!-- Request list -->
  <div class="flex-1 overflow-y-auto py-1">
    {#if collectionState.folders.length === 0 && collectionState.requests.length === 0}
      <div class="flex flex-col items-center justify-center h-full gap-2 p-4 text-center">
        <svg class="w-8 h-8 text-text-muted/30" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/>
        </svg>
        <p class="text-[11px] text-text-muted">No requests yet</p>
        <p class="text-[10px] text-text-muted/60">Send a request and save it</p>
      </div>
    {:else}
      <!-- Folders -->
      {#each collectionState.folders as folder (folder.id)}
        <div class="px-1">
          <button
            onclick={() => { folder.expanded = !folder.expanded; }}
            class="w-full flex items-center gap-1.5 px-2 py-1 rounded text-left
              text-xs text-text-secondary hover:bg-surface-2 cursor-pointer transition-colors"
          >
            <svg class="w-3 h-3 transition-transform {folder.expanded ? 'rotate-90' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
            </svg>
            <svg class="w-3.5 h-3.5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V7z"/>
            </svg>
            <span class="truncate">{folder.name}</span>
          </button>
        </div>
      {/each}

      <!-- Requests -->
      {#each collectionState.requests as req (req.id)}
        <div class="px-1">
          <div
            onclick={() => handleSelectRequest(req)}
            class="w-full flex items-center gap-2 px-2 py-1.5 rounded text-left group
              cursor-pointer transition-colors
              {collectionState.activeRequestId === req.id
                ? 'bg-accent/10 border border-accent/20'
                : 'hover:bg-surface-2 border border-transparent'}"
            role="button"
            tabindex="0"
            onkeydown={(e) => { if (e.key === 'Enter') handleSelectRequest(req); }}
          >
            <span class="text-[10px] font-bold w-9 shrink-0 {methodColor(req.method)}">
              {req.method}
            </span>
            <span class="text-xs text-text-primary truncate flex-1">{req.name}</span>
            <button
              onclick={(e) => { e.stopPropagation(); deleteRequest(req.id); }}
              class="p-0.5 rounded opacity-0 group-hover:opacity-100
                text-text-muted hover:text-error hover:bg-error/10
                transition-all cursor-pointer"
              aria-label="Delete request"
            >
              <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <!-- Sidebar Footer - Save current request -->
  {#if requestState.url}
    <div class="px-2 py-2 border-t border-border-subtle">
      <button
        onclick={handleSaveRequest}
        class="w-full flex items-center justify-center gap-1.5 px-3 py-1.5
          text-xs font-medium text-accent bg-accent/10 rounded
          hover:bg-accent/20 cursor-pointer transition-colors"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V7l-4-4z"/>
          <path stroke-linecap="round" stroke-linejoin="round" d="M7 3v5h8V3M7 14h10"/>
        </svg>
        Save
      </button>
    </div>
  {/if}
</div>
