<!--
  TabStrip - Horizontal strip of open request tabs below the titlebar.
  Each tab carries its protocol accent: a method/protocol chip, the request
  name, a dirty dot, and a close affordance. Trailing "+" opens a new tab.
-->
<script lang="ts">
  import { tabs, switchTab, closeTab } from '../stores/app.svelte';
  import { protocolMeta, methodColor } from '../lib/utils';
  import type { OpenTab } from '../lib/types';
  import Icon from './core/Icon.svelte';

  interface Props {
    onNew: () => void;
  }

  let { onNew }: Props = $props();

  /** HTTP tabs show the method chip; other protocols show their short badge. */
  function chip(tab: OpenTab): { text: string; cls: string } {
    if (tab.protocol === 'http') {
      return { text: tab.request.method, cls: methodColor(tab.request.method) };
    }
    const m = protocolMeta(tab.protocol);
    return { text: m.short, cls: m.text };
  }

  function accentVar(tab: OpenTab): string {
    return `var(--color-${protocolMeta(tab.protocol).color})`;
  }

  function label(tab: OpenTab): string {
    return tab.request.name || tab.request.url || 'Untitled';
  }

  function onMiddleClick(e: MouseEvent, id: string) {
    if (e.button === 1) { e.preventDefault(); closeTab(id); }
  }
</script>

<div class="flex items-stretch h-9 bg-surface-inset border-b border-border-subtle shrink-0 overflow-x-auto hide-scrollbars">
  {#each tabs.open as tab (tab.id)}
    {@const active = tab.id === tabs.activeId}
    {@const c = chip(tab)}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      role="tab"
      tabindex="0"
      aria-selected={active}
      onclick={() => switchTab(tab.id)}
      onauxclick={(e) => onMiddleClick(e, tab.id)}
      onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); switchTab(tab.id); } }}
      class="group relative flex items-center gap-2 pl-3 pr-2 min-w-[132px] max-w-[220px]
        border-r border-border-subtle select-none cursor-pointer
        transition-colors duration-100
        {active ? 'bg-surface' : 'bg-transparent hover:bg-surface-highlight/50'}"
    >
      <!-- Protocol chip -->
      <span class="text-[10px] font-bold font-mono tracking-tight shrink-0 {c.cls}">{c.text}</span>

      <!-- Name -->
      <span class="flex-1 truncate text-xs {active ? 'text-text' : 'text-text-subtle'}">{label(tab)}</span>

      <!-- Dirty dot / close -->
      <span class="relative flex items-center justify-center w-4 h-4 shrink-0">
        {#if tab.dirty}
          <span
            class="absolute w-1.5 h-1.5 rounded-full group-hover:opacity-0 transition-opacity"
            style="background: {accentVar(tab)}"
          ></span>
        {/if}
        <button
          type="button"
          onclick={(e) => { e.stopPropagation(); closeTab(tab.id); }}
          class="flex items-center justify-center w-4 h-4 rounded
            text-text-subtlest hover:text-text hover:bg-surface-active
            {tab.dirty ? 'opacity-0 group-hover:opacity-100' : 'opacity-60 group-hover:opacity-100'}
            transition-opacity"
          aria-label="Close tab"
          title="Close tab"
        >
          <Icon name="x" size={12} />
        </button>
      </span>

      <!-- Active accent underline -->
      {#if active}
        <div class="absolute bottom-0 left-0 right-0 h-[2px]" style="background: {accentVar(tab)}"></div>
      {/if}
    </div>
  {/each}

  <!-- New tab -->
  <button
    type="button"
    onclick={onNew}
    class="flex items-center justify-center w-9 shrink-0 text-text-subtlest
      hover:text-text hover:bg-surface-highlight/50 transition-colors"
    title="New request (Ctrl+N)"
    aria-label="New request"
  >
    <Icon name="plus" size={16} />
  </button>

  <div class="flex-1 min-w-0" style="--wails-draggable: drag"></div>
</div>
