<!--
  NewTabPopup - Popup shown when user clicks "+" to create a new request.
  Creates a new saved request in the collection and switches to it.
-->
<script lang="ts">
  import { resetRequest, request, collection, saveCurrentRequest, loadRequest } from '../stores/app.svelte';
  import { uid, emptyPair } from '../lib/utils';
  import type { RequestConfig } from '../lib/types';
  import Icon from './core/Icon.svelte';

  interface Props {
    onclose: () => void;
  }

  let { onclose }: Props = $props();

  const options = [
    { id: 'http', label: 'HTTP Request', desc: 'REST API call with any method', icon: 'world', color: 'text-success' },
    { id: 'graphql', label: 'GraphQL', desc: 'Query or mutation', icon: 'brand-graphql', color: 'text-[#e535ab]' },
    { id: 'grpc', label: 'gRPC', desc: 'Unary or streaming RPC', icon: 'server', color: 'text-info' },
    { id: 'websocket', label: 'WebSocket', desc: 'Real-time bidirectional', icon: 'plug-connected', color: 'text-warning' },
  ];

  function select(type: string) {
    // Reset current state
    resetRequest();

    // Set appropriate defaults based on type
    const name = type === 'http' ? 'New Request'
      : type === 'graphql' ? 'New GraphQL Query'
      : type === 'grpc' ? 'New gRPC Call'
      : 'New WebSocket';

    request.name = name;

    if (type === 'graphql') {
      request.method = 'POST';
      request.bodyType = 'json';
      request.body = JSON.stringify({
        query: '{\n  \n}',
        variables: {}
      }, null, 2);
      request.headers = [
        { id: uid('p'), name: 'Content-Type', value: 'application/json', enabled: true },
        emptyPair(),
      ];
    } else if (type === 'grpc') {
      request.bodyType = 'json';
      request.body = '{\n  \n}';
    }

    // Save to collection immediately so it appears in sidebar
    const now = Date.now();
    const newReq: RequestConfig = {
      id: uid('rq'),
      name,
      method: request.method,
      url: request.url,
      headers: request.headers.filter(h => h.name !== ''),
      urlParameters: [],
      bodyType: request.bodyType,
      body: request.body,
      formData: [],
      auth: { ...request.auth },
      settings: { ...request.settings },
      description: '',
      createdAt: now,
      updatedAt: now,
    };
    collection.requests.push(newReq);
    collection.activeRequestId = newReq.id;

    onclose();
  }
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
  onclick={onclose}
  onkeydown={(e) => { if (e.key === 'Escape') onclose(); }}
>
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="w-[380px] bg-popover border border-border rounded-xl shadow-2xl p-4"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <h3 class="text-sm font-semibold text-foreground mb-3">Create New</h3>
    <div class="flex flex-col gap-1">
      {#each options as opt (opt.id)}
        <button
          type="button"
          onclick={() => select(opt.id)}
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-left
            hover:bg-accent group"
        >
          <div class="p-1.5 rounded-md bg-accent group-hover:bg-background {opt.color}">
            <Icon name={opt.icon} size={18} />
          </div>
          <div>
            <div class="text-sm font-medium text-foreground">{opt.label}</div>
            <div class="text-xs text-muted-foreground">{opt.desc}</div>
          </div>
        </button>
      {/each}
    </div>
  </div>
</div>
