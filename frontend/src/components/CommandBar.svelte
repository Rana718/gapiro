<!--
  CommandBar - The protocol-aware command input: [chip] [address] [action].
  The heart of the app. The leading chip reflects the protocol (HTTP method
  dropdown, or a GQL/gRPC/WS badge), and the action verb + accent adapt.
-->
<script lang="ts">
  import type { HttpMethod, Protocol } from '../lib/types';
  import { protocolMeta } from '../lib/utils';
  import MethodSelector from './MethodSelector.svelte';
  import Icon from './core/Icon.svelte';

  interface Props {
    protocol: Protocol;
    method?: HttpMethod;
    url: string;
    loading: boolean;
    placeholder?: string;
    /** Disable the action (e.g. gRPC with no method selected). */
    disabled?: boolean;
    onMethodChange?: (m: HttpMethod) => void;
    onUrlChange: (url: string) => void;
    onSend: () => void;
    onCancel: () => void;
  }

  let {
    protocol, method = 'GET', url, loading, placeholder = 'https://api.example.com',
    disabled = false, onMethodChange, onUrlChange, onSend, onCancel,
  }: Props = $props();

  const meta = $derived(protocolMeta(protocol));
  const accent = $derived(`var(--color-${meta.color})`);

  function submit(e: Event) {
    e.preventDefault();
    if (loading) onCancel();
    else if (!disabled) onSend();
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (loading) onCancel();
      else if (!disabled) onSend();
    }
  }
</script>

<form onsubmit={submit} class="flex items-stretch gap-2 px-3 py-2.5">
  <div
    class="relative flex items-stretch flex-1 min-w-0 h-9 rounded-lg border border-border
      bg-surface-inset focus-within:border-border-focus transition-colors duration-100"
  >
    <!-- Leading chip -->
    {#if protocol === 'http'}
      <MethodSelector {method} onchange={(m) => onMethodChange?.(m)} />
    {:else}
      <div
        class="flex items-center gap-1.5 px-3 border-r border-border select-none"
        style="color: {accent}"
        title={meta.label}
      >
        <Icon name={meta.icon} size={14} />
        <span class="text-[11px] font-bold font-mono tracking-tight">{meta.short}</span>
      </div>
    {/if}

    <!-- Address input -->
    <input
      type="text"
      value={url}
      oninput={(e) => onUrlChange((e.target as HTMLInputElement).value)}
      onkeydown={onKeydown}
      {placeholder}
      class="flex-1 min-w-0 rounded-r-lg bg-transparent px-3 text-xs text-text font-mono
        placeholder:text-placeholder border-0 focus:outline-none"
      spellcheck="false"
      autocomplete="off"
      autocorrect="off"
    />
  </div>

  <!-- Action -->
  {#if loading}
    <button
      type="button"
      data-cancel-request
      onpointerdown={(e) => { e.preventDefault(); e.stopPropagation(); onCancel(); }}
      class="flex items-center gap-1.5 h-9 px-4 rounded-lg text-xs font-semibold shrink-0
        bg-danger/10 text-danger hover:bg-danger/20 transition-colors"
      title="Cancel (Esc)"
    >
      <span class="size-3.5 border-2 border-danger/40 border-t-danger rounded-full spinner"></span>
      Cancel
    </button>
  {:else}
    <button
      type="submit"
      {disabled}
      class="flex items-center gap-1.5 h-9 px-4 rounded-lg text-xs font-semibold shrink-0
        text-white transition-opacity disabled:opacity-40 disabled:cursor-not-allowed hover:opacity-90"
      style="background: {accent}"
      title="{meta.action} (Ctrl+Enter)"
    >
      <Icon name={protocol === 'websocket' ? 'plug-connected' : 'send'} size={14} />
      {meta.action}
    </button>
  {/if}
</form>
