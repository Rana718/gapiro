<!--
  NewTabPopup - "Create New" chooser. Picking a protocol opens a fresh draft tab
  via newDraftTab(), which seeds the correct protocol config so the right editor
  and backend are used. Keyboard: ↑/↓ to move, Enter to pick, Esc to close.
-->
<script lang="ts">
  import { newDraftTab } from '../stores/app.svelte';
  import { PROTOCOL_META } from '../lib/utils';
  import type { Protocol } from '../lib/types';
  import Icon from './core/Icon.svelte';

  interface Props {
    onclose: () => void;
  }

  let { onclose }: Props = $props();

  const options: { id: Protocol; desc: string }[] = [
    { id: 'http', desc: 'REST call with any method' },
    { id: 'graphql', desc: 'Query, mutation & schema docs' },
    { id: 'grpc', desc: 'Unary RPC from a .proto' },
    { id: 'websocket', desc: 'Live bidirectional stream' },
  ];

  let focused = $state(0);

  function select(protocol: Protocol) {
    newDraftTab(protocol);
    onclose();
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') { onclose(); return; }
    if (e.key === 'ArrowDown') { e.preventDefault(); focused = (focused + 1) % options.length; }
    else if (e.key === 'ArrowUp') { e.preventDefault(); focused = (focused - 1 + options.length) % options.length; }
    else if (e.key === 'Enter') { e.preventDefault(); select(options[focused].id); }
  }
</script>

<svelte:window onkeydown={onKeydown} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="fixed inset-0 z-50 flex items-start justify-center pt-[18vh] bg-black/60 backdrop-blur-[2px]"
  onclick={onclose}
>
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="w-[420px] bg-popover border border-border rounded-xl shadow-2xl overflow-hidden"
    onclick={(e) => e.stopPropagation()}
    role="dialog"
    aria-label="Create new request"
  >
    <div class="flex items-center justify-between px-4 h-11 border-b border-border-subtle">
      <span class="text-[10px] font-semibold uppercase tracking-wider text-text-subtlest">Create new request</span>
      <button
        type="button"
        onclick={onclose}
        class="flex items-center justify-center size-6 rounded-md text-text-subtlest hover:text-text hover:bg-surface-highlight transition-colors"
        aria-label="Close"
      >
        <Icon name="x" size={14} />
      </button>
    </div>

    <div class="p-2">
      {#each options as opt, i (opt.id)}
        {@const m = PROTOCOL_META[opt.id]}
        {@const accent = `var(--color-${m.color})`}
        <button
          type="button"
          onclick={() => select(opt.id)}
          onmouseenter={() => { focused = i; }}
          class="w-full flex items-center gap-3 px-2.5 py-2.5 rounded-lg text-left transition-colors
            {focused === i ? 'bg-surface-highlight' : ''}"
        >
          <div
            class="flex items-center justify-center size-9 rounded-lg shrink-0"
            style="background: color-mix(in srgb, {accent} 14%, transparent); color: {accent}"
          >
            <Icon name={m.icon} size={18} />
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-text">{m.label}</div>
            <div class="text-xs text-text-subtle truncate">{opt.desc}</div>
          </div>
          <span class="text-[10px] font-bold font-mono tracking-tight shrink-0" style="color: {accent}">{m.short}</span>
        </button>
      {/each}
    </div>
  </div>
</div>
