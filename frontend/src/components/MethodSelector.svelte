<!--
  MethodSelector - HTTP method dropdown.
  type="button" prevents form submission when inside a form.
-->
<script lang="ts">
  import type { HttpMethod } from '../lib/types';
  import { HTTP_METHODS, methodColor } from '../lib/utils';

  interface Props {
    method: HttpMethod;
    onchange: (m: HttpMethod) => void;
  }

  let { method, onchange }: Props = $props();
  let open = $state(false);

  function select(m: HttpMethod) {
    onchange(m);
    open = false;
  }
</script>

<svelte:document onclick={() => { if (open) open = false; }} />

<div class="relative">
  <button
    type="button"
    onclick={(e) => { e.stopPropagation(); open = !open; }}
    class="flex items-center gap-1 px-2.5 h-full min-h-sm rounded-l-md
      border-r border-border bg-surface-highlight hover:bg-surface-active
      text-xs font-bold select-none
      {methodColor(method)}"
  >
    {method}
    <svg class="w-2.5 h-2.5 text-text-subtlest" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
    </svg>
  </button>

  {#if open}
    <div
      class="absolute top-full left-0 mt-1 z-50 min-w-[110px] py-1
        bg-surface-active border border-border rounded-lg shadow-lg"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => { if (e.key === 'Escape') open = false; }}
      role="listbox"
      tabindex="-1"
    >
      {#each HTTP_METHODS as m (m)}
        <button
          type="button"
          onclick={() => select(m)}
          class="w-full px-3 py-1.5 text-left text-xs font-bold
            hover:bg-surface-highlight
            {methodColor(m)} {m === method ? 'bg-surface-highlight' : ''}"
          role="option"
          aria-selected={m === method}
        >
          {m}
        </button>
      {/each}
    </div>
  {/if}
</div>
