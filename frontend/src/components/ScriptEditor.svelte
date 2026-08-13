<script lang="ts">
  import { request } from '../stores/app.svelte';
  let mode = $state<'pre' | 'post'>('pre');
</script>

<div class="flex h-full flex-col">
  <div class="flex items-center gap-1 border-b border-border-subtle bg-surface-inset px-3 py-2">
    <button type="button" onclick={() => mode = 'pre'} class="rounded px-2.5 py-1 text-xs font-medium {mode === 'pre' ? 'bg-primary/15 text-primary' : 'text-text-subtlest hover:bg-surface-highlight'}">Pre-request</button>
    <button type="button" onclick={() => mode = 'post'} class="rounded px-2.5 py-1 text-xs font-medium {mode === 'post' ? 'bg-primary/15 text-primary' : 'text-text-subtlest hover:bg-surface-highlight'}">Post-response</button>
    <span class="ml-auto text-[10px] text-text-subtlest">Saved with this request</span>
  </div>
  <textarea value={mode === 'pre' ? request.preRequestScript : request.postResponseScript} oninput={(e) => { if (mode === 'pre') request.preRequestScript = e.currentTarget.value; else request.postResponseScript = e.currentTarget.value; }} placeholder={mode === 'pre' ? '// Prepare variables before sending' : '// Inspect the response and add tests'} class="flex-1 resize-none border-0 bg-surface-inset p-4 font-mono text-xs leading-relaxed text-text placeholder:text-placeholder focus:outline-none" spellcheck="false"></textarea>
</div>
