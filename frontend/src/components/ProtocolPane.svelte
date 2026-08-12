<script lang="ts">
  import type { Snippet } from 'svelte';
  let { children }: { children: Snippet } = $props();
  let protocol = $state<'http'|'websocket'|'graphql'|'grpc'|'sql'>('http');
  let endpoint = $state('');
  let message = $state('');
  let connected = $state(false);
  const protocols = [
    {id:'http', label:'HTTP'}, {id:'websocket', label:'WebSocket'},
    {id:'graphql', label:'GraphQL'}, {id:'grpc', label:'gRPC'}, {id:'sql', label:'SQL'}
  ];
</script>
<div class="flex flex-col h-full bg-surface rounded-md border border-border-subtle overflow-hidden">
  <div class="flex items-center gap-1 px-3 py-2 border-b border-border-subtle shrink-0">
    {#each protocols as p}<button class="protocol-button {protocol===p.id?'active':''}" onclick={()=>protocol=p.id as typeof protocol}>
      <svg viewBox="0 0 24 24" aria-hidden="true">
        {#if p.id==='http'}<path d="M5 12h14M13 6l6 6-6 6"/>
        {:else if p.id==='websocket'}<path d="M8 8a6 6 0 0 1 8 0M6 5a9 9 0 0 1 12 0M8 16a6 6 0 0 0 8 0M6 19a9 9 0 0 0 12 0"/>
        {:else if p.id==='graphql'}<path d="m12 3 7.8 4.5v9L12 21l-7.8-4.5v-9L12 3Z M4.5 7.7h15M4.5 16.3h15M12 3v18"/>
        {:else if p.id==='grpc'}<path d="M4 7h7v10H4zM13 4h7v7h-7zM13 13h7v7h-7zM11 9h2M11 15h2"/>
        {:else}<ellipse cx="12" cy="5" rx="7" ry="3"/><path d="M5 5v7c0 1.7 3.1 3 7 3s7-1.3 7-3V5M5 12v7c0 1.7 3.1 3 7 3s7-1.3 7-3v-7"/>{/if}
      </svg><span>{p.label}</span>
    </button>{/each}
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

<style>
  .protocol-button{height:30px;padding:0 10px;display:inline-flex;align-items:center;gap:7px;border-radius:4px;color:var(--color-text-subtlest);font-size:11px;font-weight:600}
  .protocol-button:hover{background:var(--color-surface-highlight);color:var(--color-text)}
  .protocol-button.active{background:color-mix(in srgb,var(--color-primary) 14%,transparent);color:var(--color-primary)}
  .protocol-button svg{width:15px;height:15px;fill:none;stroke:currentColor;stroke-width:1.7;stroke-linecap:round;stroke-linejoin:round;flex:none}
</style>
