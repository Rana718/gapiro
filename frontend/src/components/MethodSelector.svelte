<!--
  MethodSelector - HTTP method dropdown with color coding.
-->
<script lang="ts">
  import type { HttpMethod } from '../lib/types';
  import { methodColor } from '../lib/utils';

  interface Props {
    method: HttpMethod;
    onchange: (method: HttpMethod) => void;
  }

  let { method, onchange }: Props = $props();
  let open = $state(false);

  const methods: HttpMethod[] = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'];

  function select(m: HttpMethod) {
    onchange(m);
    open = false;
  }

  function handleClickOutside(e: MouseEvent) {
    if (open) open = false;
  }
</script>

<svelte:document onclick={handleClickOutside} />

<div class="relative">
  <button
    onclick={(e) => { e.stopPropagation(); open = !open; }}
    class="flex items-center gap-1 px-2.5 py-1.5 rounded-l border border-r-0
      border-border-default bg-surface-2 hover:bg-surface-3
      text-xs font-bold cursor-pointer transition-colors duration-75
      {methodColor(method)}"
  >
    {method}
    <svg class="w-3 h-3 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
    </svg>
  </button>

  {#if open}
    <div
      class="absolute top-full left-0 mt-1 z-50 min-w-[120px] py-1
        bg-surface-3 border border-border-default rounded-lg shadow-xl shadow-black/30
        gpu-layer"
      onclick={(e) => e.stopPropagation()}
      role="listbox"
    >
      {#each methods as m (m)}
        <button
          onclick={() => select(m)}
          class="w-full px-3 py-1.5 text-left text-xs font-bold cursor-pointer
            hover:bg-surface-4 transition-colors duration-75
            {methodColor(m)} {m === method ? 'bg-surface-4' : ''}"
        >
          {m}
        </button>
      {/each}
    </div>
  {/if}
</div>
