<script lang="ts">
  import type { Snippet } from 'svelte';
  let { children }: { children: Snippet } = $props();
  let protocol = $state<'http'|'websocket'|'graphql'|'grpc'|'sql'>('http');
  let endpoint = $state('');
  let message = $state('');
  let connected = $state(false);
  const protocols = [
    {id:'http', label:'HTTP', icon:'↗'}, {id:'websocket', label:'WebSocket', icon:'◉'},
    {id:'graphql', label:'GraphQL', icon:'◇'}, {id:'grpc', label:'gRPC', icon:'▣'}, {id:'sql', label:'SQL', icon:'⌘'}
  ];
</script>
<div class="flex flex-col h-full bg-surface rounded-md border border-border-subtle overflow-hidden">
  <div class="flex items-center gap-1 px-3 py-2 border-b border-border-subtle shrink-0">
    {#each protocols as p}<button class="px-3 py-1.5 text-xs rounded-md {protocol===p.id?'bg-primary/15 text-primary':'text-text-subtlest hover:bg-surface-highlight'}" onclick={()=>protocol=p.id as typeof protocol}>{p.icon} {p.label}</button>{/each}
  </div>
  {#if protocol === 'http'}
    <div class="flex-1 min-h-0">{@render children()}</div>
  {:else}
    <div class="flex flex-col h-full">
      <div class="flex items-center gap-2 p-3 border-b border-border-subtle"><input bind:value={endpoint} placeholder={protocol==='sql'?'postgres://user:pass@host/db':'Service endpoint…'} class="flex-1 h-9 px-3 rounded-md bg-surface-inset border border-border text-sm font-mono focus:outline-none focus:border-border-focus"/><button onclick={()=>connected=!connected} class="h-9 px-4 rounded-md {connected?'bg-success/15 text-success':'bg-primary text-white'}">{connected?'Disconnect':'Connect'}</button></div>
      {#if protocol==='websocket'}<div class="flex-1 p-3 flex flex-col gap-3"><div class="rounded-md border border-border-subtle bg-surface-inset p-3 text-xs text-text-subtle">{connected?'Connected and ready for messages':'Enter a WebSocket URL to begin'}</div><textarea bind:value={message} placeholder="Message payload…" class="flex-1 p-3 rounded-md bg-surface-inset border border-border font-mono text-xs resize-none"/><button class="self-end px-4 py-2 rounded-md bg-primary text-white" onclick={()=>message='' }>Send message</button></div>
      {:else if protocol==='graphql'}<div class="grid grid-cols-2 gap-2 flex-1 p-3"><textarea placeholder="GraphQL query…" class="p-3 rounded-md bg-surface-inset border border-border font-mono text-xs resize-none"></textarea><textarea placeholder="Variables (JSON)" class="p-3 rounded-md bg-surface-inset border border-border font-mono text-xs resize-none"></textarea></div>
      {:else if protocol==='grpc'}<div class="flex-1 p-4 text-sm text-text-subtle">Select a .proto service and method to configure unary or streaming calls.</div>
      {:else}<div class="flex-1 p-3"><textarea placeholder="SELECT * FROM users LIMIT 50;" class="w-full h-full p-3 rounded-md bg-surface-inset border border-border font-mono text-xs resize-none"/></div>{/if}
    </div>
  {/if}
</div>
