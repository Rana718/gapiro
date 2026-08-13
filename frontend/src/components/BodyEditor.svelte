<!--
  BodyEditor - Request body editor with Postman-style type selector.
  Radio buttons: none | form-data | x-www-form-urlencoded | raw | binary | GraphQL
  When "raw" is selected, a dropdown shows format: JSON, Text, XML, HTML, JavaScript.
-->
<script lang="ts">
  import { request } from '../stores/app.svelte';
  import type { BodyType } from '../lib/types';
  import PairEditor from './core/PairEditor.svelte';
  import FormDataEditor from './FormDataEditor.svelte';
  import * as Select from './ui/select/index';

  // Map between our internal types and Postman-style radio groups
  type RadioOption = 'none' | 'form-data' | 'x-www-form-urlencoded' | 'raw' | 'binary' | 'graphql';
  type RawFormat = 'json' | 'text' | 'xml' | 'html' | 'javascript';

  const rawFormats: { id: RawFormat; label: string }[] = [
    { id: 'json', label: 'JSON' },
    { id: 'text', label: 'Text' },
    { id: 'xml', label: 'XML' },
    { id: 'html', label: 'HTML' },
    { id: 'javascript', label: 'JavaScript' },
  ];

  // Derive which radio is active from the internal bodyType
  let activeRadio = $derived.by<RadioOption>(() => {
    switch (request.bodyType) {
      case 'none': return 'none';
      case 'form-data': return 'form-data';
      case 'form-urlencoded': return 'x-www-form-urlencoded';
      case 'binary': return 'binary';
      case 'graphql': return 'graphql';
      case 'json': case 'xml': case 'text': return 'raw';
      default: return 'none';
    }
  });

  // Derive the raw format from the internal bodyType
  let rawFormat = $derived.by<RawFormat>(() => {
    switch (request.bodyType) {
      case 'json': return 'json';
      case 'xml': return 'xml';
      case 'text': return 'text';
      default: return 'json'; // default raw format
    }
  });

  function selectRadio(opt: RadioOption) {
    switch (opt) {
      case 'none': request.bodyType = 'none'; break;
      case 'form-data': request.bodyType = 'form-data'; break;
      case 'x-www-form-urlencoded': request.bodyType = 'form-urlencoded'; break;
      case 'raw': request.bodyType = 'json'; break; // default raw to JSON
      case 'binary': request.bodyType = 'binary'; break;
      case 'graphql': request.bodyType = 'graphql'; break;
    }
  }

  function selectRawFormat(fmt: RawFormat) {
    switch (fmt) {
      case 'json': request.bodyType = 'json'; break;
      case 'xml': request.bodyType = 'xml'; break;
      case 'text': request.bodyType = 'text'; break;
      case 'html': request.bodyType = 'text'; break; // treat html as text with label
      case 'javascript': request.bodyType = 'text'; break; // treat js as text with label
    }
  }

  let useCodeMirror = $state(false);

  $effect(() => {
    const needsCM = request.bodyType === 'json' || request.bodyType === 'xml' || request.bodyType === 'graphql';
    if (needsCM && !useCodeMirror) {
      const t = setTimeout(() => { useCodeMirror = true; }, 100);
      return () => clearTimeout(t);
    }
    if (!needsCM) useCodeMirror = false;
  });

  const radioOptions: { id: RadioOption; label: string }[] = [
    { id: 'none', label: 'none' },
    { id: 'form-data', label: 'form-data' },
    { id: 'x-www-form-urlencoded', label: 'x-www-form-urlencoded' },
    { id: 'raw', label: 'raw' },
    { id: 'binary', label: 'binary' },
    { id: 'graphql', label: 'GraphQL' },
  ];
</script>

<div class="flex flex-col h-full">
  <!-- Body type selector - Postman style radio buttons -->
  <div class="flex items-center gap-4 px-4 py-2.5 border-b border-border-subtle shrink-0">
    <div class="flex items-center gap-3">
      {#each radioOptions as opt (opt.id)}
        <label class="flex items-center gap-1.5 cursor-pointer group">
          <input
            type="radio"
            name="body-type"
            checked={activeRadio === opt.id}
            onchange={() => selectRadio(opt.id)}
            class="size-3.5 accent-primary cursor-pointer"
          />
          <span class="text-xs {activeRadio === opt.id ? 'text-text font-medium' : 'text-text-subtle'} group-hover:text-text whitespace-nowrap">
            {opt.label}
          </span>
        </label>
      {/each}
    </div>

    <!-- Raw format dropdown (only when raw is selected) -->
    {#if activeRadio === 'raw'}
      <div class="ml-auto flex items-center gap-2">
        <select
          value={rawFormat}
          onchange={(e) => selectRawFormat((e.currentTarget as HTMLSelectElement).value as RawFormat)}
          class="h-7 px-2 rounded border border-border bg-surface-inset text-xs text-primary font-medium
            focus:outline-none focus:border-border-focus cursor-pointer"
        >
          {#each rawFormats as fmt (fmt.id)}
            <option value={fmt.id}>{fmt.label}</option>
          {/each}
        </select>
      </div>
    {/if}
  </div>

  <!-- Body content -->
  <div class="flex-1 overflow-hidden">
    {#if request.bodyType === 'none'}
      <div class="flex flex-col items-center justify-center h-full gap-2 text-text-subtlest text-sm">
        <span class="text-text-subtle">This request does not have a body</span>
      </div>
    {:else if request.bodyType === 'form-data'}
      <FormDataEditor pairs={request.formData} onchange={(p) => { request.formData = p; }} />
    {:else if request.bodyType === 'form-urlencoded'}
      <PairEditor
        pairs={request.formData}
        namePlaceholder="Key"
        valuePlaceholder="Value"
        onchange={(p) => { request.formData = p; }}
      />
    {:else if request.bodyType === 'binary'}
      <div class="flex h-full items-center justify-center p-6">
        <label class="flex w-full max-w-lg cursor-pointer flex-col items-center gap-3 rounded-xl border border-dashed border-border bg-surface-inset px-6 py-10 text-center hover:border-border-focus hover:bg-surface-highlight">
          <span class="text-sm font-medium text-text">Select File</span>
          <span class="text-xs text-text-subtlest">The file will be sent as the request body</span>
          <span class="max-w-full truncate rounded bg-surface-active px-3 py-1.5 font-mono text-xs text-text-subtle">{request.body || 'No file selected'}</span>
          <input type="file" class="sr-only" onchange={(e) => { const file = e.currentTarget.files?.[0] as (File & { path?: string }) | undefined; if (file) request.body = file.path ?? file.name; }} />
        </label>
      </div>
    {:else if useCodeMirror}
      {#await import('./core/CodeEditor.svelte') then { default: CodeEditor }}
        <CodeEditor
          value={request.body}
          onchange={(v) => { request.body = v; }}
          language={request.bodyType === 'xml' ? 'xml' : request.bodyType === 'graphql' ? 'javascript' : 'json'}
          placeholder={request.bodyType === 'json' ? '{\n  "key": "value"\n}' : request.bodyType === 'graphql' ? 'query {\n  \n}' : '<root>\n  \n</root>'}
        />
      {/await}
    {:else}
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
