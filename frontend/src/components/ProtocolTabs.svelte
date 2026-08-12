<!--
  ProtocolTabs - Simple protocol selector at the top of the request area.
  Lightweight - just buttons, no extra DOM.
-->
<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();

  let protocol = $state<'http' | 'websocket' | 'graphql' | 'grpc'>('http');

  const protocols = [
    { id: 'http', label: 'HTTP' },
    { id: 'websocket', label: 'WebSocket' },
    { id: 'graphql', label: 'GraphQL' },
    { id: 'grpc', label: 'gRPC' },
  ] as const;
</script>

<div class="flex flex-col h-full">
  <!-- Protocol tabs -->
  <div class="flex items-center gap-0.5 px-3 py-1 border-b border-border-subtle shrink-0 bg-surface">
    {#each protocols as p (p.id)}
      <button
        type="button"
        onclick={() => { protocol = p.id; }}
        class="px-2.5 py-1 text-[11px] font-medium rounded
          {protocol === p.id
            ? 'bg-surface-highlight text-text'
            : 'text-text-subtlest hover:text-text-subtle'}"
      >
        {p.label}
      </button>
    {/each}
  </div>

  <!-- Content -->
  <div class="flex-1 overflow-hidden">
    {#if protocol === 'http'}
      {@render children()}
    {:else if protocol === 'websocket'}
      <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
        WebSocket support coming soon
      </div>
    {:else if protocol === 'graphql'}
      <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
        GraphQL support coming soon
      </div>
    {:else if protocol === 'grpc'}
      <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
        gRPC support coming soon
      </div>
    {/if}
  </div>
</div>
