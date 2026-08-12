<!--
  BodyEditor - Request body editor with type selector.
  Uses CodeMirror only for JSON/XML (lazy loaded), plain textarea for text.
-->
<script lang="ts">
  import { request } from '../stores/app.svelte';
  import type { BodyType } from '../lib/types';
  import PairEditor from './core/PairEditor.svelte';

  const bodyTypes: { id: BodyType; label: string }[] = [
    { id: 'none', label: 'None' },
    { id: 'json', label: 'JSON' },
    { id: 'text', label: 'Text' },
    { id: 'xml', label: 'XML' },
    { id: 'form-urlencoded', label: 'Form' },
    { id: 'form-data', label: 'Multipart' },
  ];

  let useCodeMirror = $state(false);

  // Only load CodeMirror for JSON/XML after a short delay (avoids blocking render)
  $effect(() => {
    const needsCM = request.bodyType === 'json' || request.bodyType === 'xml';
    if (needsCM && !useCodeMirror) {
      // Delay loading to avoid jank
      const t = setTimeout(() => { useCodeMirror = true; }, 100);
      return () => clearTimeout(t);
    }
    if (!needsCM) useCodeMirror = false;
  });
</script>

<div class="flex flex-col h-full">
  <!-- Body type selector -->
  <div class="flex items-center gap-1 px-3 py-1.5 border-b border-border-subtle shrink-0">
    {#each bodyTypes as bt (bt.id)}
      <button
        type="button"
        onclick={() => { request.bodyType = bt.id; }}
        class="px-2 py-1 text-[11px] font-medium rounded whitespace-nowrap
          {request.bodyType === bt.id
            ? 'bg-primary/15 text-primary'
            : 'text-text-subtlest hover:text-text-subtle hover:bg-surface-highlight'}"
      >
        {bt.label}
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
    {:else if useCodeMirror}
      <!-- CodeMirror for JSON/XML (lazy) -->
      {#await import('./core/CodeEditor.svelte') then { default: CodeEditor }}
        <CodeEditor
          value={request.body}
          onchange={(v) => { request.body = v; }}
          language={request.bodyType === 'xml' ? 'xml' : 'json'}
          placeholder={request.bodyType === 'json' ? '{\n  "key": "value"\n}' : '<root>\n  \n</root>'}
        />
      {/await}
    {:else}
      <!-- Plain textarea for text and initial load -->
      <textarea
        value={request.body}
        oninput={(e) => { request.body = (e.target as HTMLTextAreaElement).value; }}
        placeholder={request.bodyType === 'json' ? '{\n  "key": "value"\n}'
          : request.bodyType === 'xml' ? '<root>\n  \n</root>'
          : 'Enter request body...'}
        class="w-full h-full bg-surface-inset text-text text-xs font-mono
          p-3 resize-none border-0 focus:outline-none
          placeholder:text-placeholder leading-relaxed"
        spellcheck="false"
      ></textarea>
    {/if}
  </div>
</div>
