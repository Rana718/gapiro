<!--
  UrlBar - The main URL input bar with method selector and send button.
-->
<script lang="ts">
  import type { HttpMethod } from '../lib/types';
  import MethodSelector from './MethodSelector.svelte';

  interface Props {
    method: HttpMethod;
    url: string;
    loading: boolean;
    onMethodChange: (method: HttpMethod) => void;
    onUrlChange: (url: string) => void;
    onSend: () => void;
    onCancel: () => void;
  }

  let { method, url, loading, onMethodChange, onUrlChange, onSend, onCancel }: Props = $props();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      e.preventDefault();
      onSend();
    }
  }
</script>

<div class="flex items-center gap-0 px-3 py-2">
  <div class="flex items-center flex-1 border border-border-default rounded
    focus-within:border-accent transition-colors duration-100 bg-surface-0">
    <!-- Method selector -->
    <MethodSelector {method} onchange={onMethodChange} />

    <!-- URL input -->
    <input
      type="text"
      value={url}
      oninput={(e) => onUrlChange((e.target as HTMLInputElement).value)}
      onkeydown={handleKeydown}
      placeholder="Enter URL or paste cURL"
      class="flex-1 bg-transparent px-3 py-1.5 text-sm text-text-primary
        placeholder:text-text-muted/60 border-0 focus:outline-none"
      spellcheck="false"
      autocomplete="off"
    />

    <!-- Send/Cancel button -->
    {#if loading}
      <button
        onclick={onCancel}
        class="flex items-center gap-1.5 px-4 py-1.5 mr-0.5 rounded-r
          bg-error hover:bg-error/80 text-white text-xs font-semibold
          cursor-pointer transition-colors duration-75"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        </svg>
        Cancel
      </button>
    {:else}
      <button
        onclick={onSend}
        class="flex items-center gap-1.5 px-4 py-1.5 mr-0.5 rounded-r
          bg-accent hover:bg-accent-hover text-white text-xs font-semibold
          cursor-pointer transition-colors duration-75"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 12h14"/>
        </svg>
        Send
      </button>
    {/if}
  </div>
</div>
