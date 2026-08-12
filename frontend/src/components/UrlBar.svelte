<!--
  UrlBar - Main URL input with method selector and send button.
  Matches Yaak's UrlBar structure: [Method] [URL Input] [Send]
-->
<script lang="ts">
  import type { HttpMethod } from '../lib/types';
  import MethodSelector from './MethodSelector.svelte';

  interface Props {
    method: HttpMethod;
    url: string;
    loading: boolean;
    onMethodChange: (m: HttpMethod) => void;
    onUrlChange: (url: string) => void;
    onSend: () => void;
    onCancel: () => void;
  }

  let { method, url, loading, onMethodChange, onUrlChange, onSend, onCancel }: Props = $props();

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (loading) onCancel();
    else onSend();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (loading) onCancel();
      else onSend();
    }
  }
</script>

<form onsubmit={handleSubmit} class="flex items-stretch px-3 py-2">
  <div class="flex items-stretch flex-1 rounded-md border border-border
    bg-surface-inset focus-within:border-border-focus transition-colors duration-100">
    <!-- Method dropdown -->
    <MethodSelector {method} onchange={onMethodChange} />

    <!-- URL input -->
    <input
      type="text"
      value={url}
      oninput={(e) => onUrlChange((e.target as HTMLInputElement).value)}
      onkeydown={handleKeydown}
      placeholder="https://api.example.com/endpoint"
      class="flex-1 bg-transparent px-3 py-0 text-sm text-text font-mono
        placeholder:text-placeholder border-0 focus:outline-none
        min-w-0"
      spellcheck="false"
      autocomplete="off"
      autocorrect="off"
    />

    <!-- Send / Cancel -->
    {#if loading}
      <button
        type="button"
        onclick={onCancel}
        class="flex items-center justify-center w-8 h-full mr-0.5 my-0.5 rounded-md
          text-text-subtle hover:text-danger hover:bg-danger/10
          transition-colors duration-75"
        title="Cancel Request"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    {:else}
      <button
        type="submit"
        class="flex items-center justify-center w-8 h-full mr-0.5 my-0.5 rounded-md
          text-text-subtle hover:text-primary hover:bg-primary/10
          transition-colors duration-75"
        title="Send Request (Ctrl+Enter)"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 12h14"/>
        </svg>
      </button>
    {/if}
  </div>
</form>
