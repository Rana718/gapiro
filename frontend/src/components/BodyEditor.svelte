<!--
  BodyEditor - Request body editor with type selector.
  Supports: none, JSON, XML, text, form-urlencoded, form-data, GraphQL, binary.
-->
<script lang="ts">
  import { request } from '../stores/app.svelte';
  import type { BodyType } from '../lib/types';
  import { BODY_TYPES } from '../lib/utils';
  import PairEditor from './core/PairEditor.svelte';

  function setBodyType(bt: string) {
    request.bodyType = bt as BodyType;
    // Auto-set method to POST for body types that need it
    if (bt !== 'none' && request.method === 'GET') {
      request.method = 'POST';
    }
  }
</script>

<div class="flex flex-col h-full">
  <!-- Body type selector -->
  <div class="flex items-center gap-1 px-3 py-1.5 border-b border-border-subtle overflow-x-auto hide-scrollbars">
    {#each BODY_TYPES as bt (bt.id)}
      <button
        onclick={() => setBodyType(bt.id)}
        class="px-2 py-1 text-[11px] font-medium rounded whitespace-nowrap
          transition-colors duration-50
          {request.bodyType === bt.id
            ? 'bg-primary/15 text-primary border border-primary/30'
            : 'text-text-subtlest hover:text-text-subtle hover:bg-surface-highlight border border-transparent'}"
      >
        {bt.short}
      </button>
    {/each}
  </div>

  <!-- Body content -->
  <div class="flex-1 overflow-hidden">
    {#if request.bodyType === 'none'}
      <div class="flex items-center justify-center h-full text-text-subtlest text-sm">
        This request does not have a body
      </div>
    {:else if request.bodyType === 'form-urlencoded' || request.bodyType === 'form-data'}
      <PairEditor
        pairs={request.formData}
        namePlaceholder="Field"
        valuePlaceholder="Value"
        onchange={(p) => { request.formData = p; }}
      />
    {:else}
      <!-- Text-based body (JSON, XML, text, GraphQL) -->
      <div class="relative w-full h-full">
        {#await import('./core/CodeEditor.svelte') then { default: CodeEditor }}
          <CodeEditor
            value={request.body}
            onchange={(v) => { request.body = v; }}
            language={request.bodyType === 'graphql' ? 'graphql' : request.bodyType}
            placeholder={request.bodyType === 'json' ? '{\n  "key": "value"\n}'
              : request.bodyType === 'graphql' ? 'query {\n  \n}'
              : request.bodyType === 'xml' ? '<root>\n  \n</root>'
              : 'Enter request body...'}
          />
        {/await}
      </div>
    {/if}
  </div>
</div>
