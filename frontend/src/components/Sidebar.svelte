<!--
  Sidebar - Saved request collection with folders. "+" opens the New Request
  chooser; requests can be filed into folders via the right-click menu and are
  rendered nested under them. Each row shows its protocol chip.
-->
<script lang="ts">
  import {
    collection, request, openSavedTab, saveCurrentRequest,
    deleteRequest, duplicateRequest, createFolder, moveRequestToFolder,
  } from '../stores/app.svelte';
  import { methodColor, protocolMeta } from '../lib/utils';
  import type { RequestConfig } from '../lib/types';
  import Icon from './core/Icon.svelte';

  interface Props {
    onNew: () => void;
  }
  let { onNew }: Props = $props();

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

  const rootRequests = $derived(filteredRequests.filter(r => !r.folderId));
  const byFolder = (folderId: string) => filteredRequests.filter(r => r.folderId === folderId);

  /** HTTP rows show the method; other protocols show their short badge. */
  function chip(req: RequestConfig): { text: string; cls: string } {
    if (!req.protocol || req.protocol === 'http') return { text: req.method, cls: methodColor(req.method) };
    const m = protocolMeta(req.protocol);
    return { text: m.short, cls: m.text };
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
    contextMenu = { x: Math.min(e.clientX, window.innerWidth - 170), y: e.clientY, requestId: reqId };
  }
  function closeContextMenu() { contextMenu = null; }
</script>

<svelte:document onclick={closeContextMenu} />

<div class="flex flex-col h-full bg-surface-inset">
  <!-- Header -->
  <div class="flex items-center justify-between px-3 h-9 border-b border-border-subtle shrink-0">
    <span class="text-[10px] font-semibold text-text-subtlest uppercase tracking-wider">Collection</span>
    <div class="flex items-center gap-0.5">
      <button
        onclick={() => { showNewFolder = !showNewFolder; }}
        class="flex items-center justify-center size-6 rounded-md text-text-subtlest hover:text-text hover:bg-surface-highlight transition-colors"
        title="New folder"
      >
        <Icon name="folder-plus" size={14} />
      </button>
      <button
        onclick={onNew}
        class="flex items-center justify-center size-6 rounded-md text-text-subtlest hover:text-text hover:bg-surface-highlight transition-colors"
        title="New request (Ctrl+N)"
      >
        <Icon name="plus" size={15} />
      </button>
    </div>
  </div>

  <!-- Search/Filter -->
  <div class="px-2 py-1.5 border-b border-border-subtle shrink-0">
    <div class="flex items-center gap-1.5 px-2.5 h-7 bg-surface border border-border-subtle rounded-md focus-within:border-border-focus transition-colors">
      <Icon name="search" size={12} class="text-text-subtlest shrink-0" />
      <input
        type="text"
        bind:value={filterText}
        placeholder="Filter requests…"
        class="flex-1 min-w-0 bg-transparent text-xs text-text placeholder:text-placeholder border-0 focus:outline-none"
      />
    </div>
  </div>

  <!-- New folder input -->
  {#if showNewFolder}
    <div class="px-2 py-1.5 border-b border-border-subtle">
      <form onsubmit={(e) => { e.preventDefault(); handleCreateFolder(); }}>
        <!-- svelte-ignore a11y_autofocus -->
        <input
          type="text"
          bind:value={newFolderName}
          placeholder="Folder name…"
          autofocus
          onblur={() => { if (!newFolderName.trim()) showNewFolder = false; }}
          class="w-full px-2.5 py-1 text-xs bg-surface border border-border-subtle rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </form>
    </div>
  {/if}

  <!-- Request list -->
  <div class="flex-1 overflow-y-auto py-1">
    {#if collection.folders.length === 0 && filteredRequests.length === 0}
      <div class="flex flex-col items-center justify-center h-full gap-2 p-4 text-center select-none">
        <Icon name="folder" size={32} class="text-text-subtlest/30" />
        <p class="text-[11px] text-text-subtle">No saved requests</p>
        <p class="text-[10px] text-text-subtlest">Send a request and save it here</p>
      </div>
    {:else}
      <!-- Folders with nested requests -->
      {#each collection.folders as folder (folder.id)}
        {@const reqs = byFolder(folder.id)}
        <div class="px-1">
          <button
            onclick={() => { folder.expanded = !folder.expanded; }}
            class="w-full flex items-center gap-1.5 px-2 py-1 rounded-md text-left
              text-xs text-text-subtle hover:bg-surface-highlight transition-colors"
          >
            <Icon name="chevron-right" size={12} class="text-text-subtlest transition-transform duration-100 {folder.expanded ? 'rotate-90' : ''}" />
            <Icon name="folder" size={14} class="text-text-subtlest" />
            <span class="truncate flex-1">{folder.name}</span>
            {#if reqs.length}<span class="text-[10px] text-text-subtlest tabular-nums">{reqs.length}</span>{/if}
          </button>

          {#if folder.expanded}
            <div class="ml-3 border-l border-border-subtle pl-1">
              {#each reqs as req (req.id)}
                {@const c = chip(req)}
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div
                  onclick={() => openSavedTab(req)}
                  oncontextmenu={(e) => handleContextMenu(e, req.id)}
                  onkeydown={(e) => { if (e.key === 'Enter') openSavedTab(req); }}
                  role="button"
                  tabindex="0"
                  class="w-full flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer border transition-colors duration-75
                    {collection.activeRequestId === req.id ? 'bg-primary/10 border-primary/25' : 'border-transparent hover:bg-surface-highlight'}"
                >
                  <span class="text-[9px] font-bold w-[34px] shrink-0 text-right font-mono {c.cls}">{c.text}</span>
                  <span class="text-xs text-text truncate flex-1">{req.name || req.url || 'Untitled'}</span>
                </div>
              {/each}
              {#if !reqs.length}
                <div class="px-2 py-1 text-[10px] text-text-subtlest italic select-none">Empty</div>
              {/if}
            </div>
          {/if}
        </div>
      {/each}

      <!-- Root-level requests -->
      {#each rootRequests as req (req.id)}
        {@const c = chip(req)}
        <div class="px-1">
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            onclick={() => openSavedTab(req)}
            oncontextmenu={(e) => handleContextMenu(e, req.id)}
            onkeydown={(e) => { if (e.key === 'Enter') openSavedTab(req); }}
            role="button"
            tabindex="0"
            class="w-full flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer border transition-colors duration-75
              {collection.activeRequestId === req.id ? 'bg-primary/10 border-primary/25' : 'border-transparent hover:bg-surface-highlight'}"
          >
            <span class="text-[9px] font-bold w-[34px] shrink-0 text-right font-mono {c.cls}">{c.text}</span>
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
        onclick={() => saveCurrentRequest()}
        class="w-full flex items-center justify-center gap-1.5 px-3 py-1.5
          text-xs font-medium text-primary bg-primary/10 rounded-md
          hover:bg-primary/20 transition-colors"
      >
        <Icon name="save" size={14} />
        Save request
      </button>
    </div>
  {/if}

  <!-- Context menu -->
  {#if contextMenu}
    {@const cm = contextMenu}
    <div
      class="fixed z-50 py-1 min-w-[160px] bg-surface-active border border-border rounded-lg shadow-xl"
      style="left: {cm.x}px; top: {cm.y}px;"
    >
      <button
        onclick={() => { duplicateRequest(cm.requestId); closeContextMenu(); }}
        class="w-full flex items-center gap-2 px-3 py-1.5 text-xs text-left text-text hover:bg-surface-highlight transition-colors"
      ><Icon name="copy" size={13} /> Duplicate</button>

      {#if collection.folders.length}
        <div class="my-1 border-t border-border-subtle"></div>
        <div class="px-3 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-text-subtlest">Move to</div>
        {#each collection.folders as folder (folder.id)}
          <button
            onclick={() => { moveRequestToFolder(cm.requestId, folder.id); closeContextMenu(); }}
            class="w-full flex items-center gap-2 px-3 py-1.5 text-xs text-left text-text hover:bg-surface-highlight transition-colors"
          ><Icon name="folder" size={13} class="text-text-subtlest" /> {folder.name}</button>
        {/each}
        <button
          onclick={() => { moveRequestToFolder(cm.requestId, undefined); closeContextMenu(); }}
          class="w-full flex items-center gap-2 px-3 py-1.5 text-xs text-left text-text-subtle hover:bg-surface-highlight transition-colors"
        ><Icon name="arrow-up" size={13} /> Root</button>
      {/if}

      <div class="my-1 border-t border-border-subtle"></div>
      <button
        onclick={() => { deleteRequest(cm.requestId); closeContextMenu(); }}
        class="w-full flex items-center gap-2 px-3 py-1.5 text-xs text-left text-danger hover:bg-danger/10 transition-colors"
      ><Icon name="trash" size={13} /> Delete</button>
    </div>
  {/if}
</div>
