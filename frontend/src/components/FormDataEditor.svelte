<script lang="ts">
  import type { Pair } from '../lib/types';
  import { emptyPair } from '../lib/utils';
  import Icon from './core/Icon.svelte';

  interface Props { pairs: Pair[]; onchange: (pairs: Pair[]) => void; }
  let { pairs, onchange }: Props = $props();

  function change(i: number, patch: Partial<Pair>) {
    Object.assign(pairs[i], patch);
    const last = pairs[pairs.length - 1];
    if (last && (last.name || last.value)) pairs.push(emptyPair());
    onchange(pairs);
  }
  function pick(i: number, files: FileList | null) {
    const file = files?.[0];
    if (!file) return;
    change(i, { value: (file as File & { path?: string }).path ?? file.name, fileName: file.name, valueType: 'file' });
  }
  function remove(i: number) {
    if (pairs.length === 1) Object.assign(pairs[0], emptyPair());
    else pairs.splice(i, 1);
    onchange(pairs);
  }
</script>

<div class="flex h-full flex-col overflow-y-auto">
  <div class="grid grid-cols-[28px_120px_1fr_1.4fr_28px] gap-1 border-b border-border-subtle px-3 py-2 text-[10px] font-semibold uppercase tracking-wider text-text-subtlest">
    <span></span><span>Type</span><span>Key</span><span>Value</span><span></span>
  </div>
  {#each pairs as pair, i (pair.id)}
    <div class="group grid grid-cols-[28px_120px_1fr_1.4fr_28px] items-center gap-1 border-b border-border-subtle/60 px-3 py-1.5">
      <input type="checkbox" checked={pair.enabled} onchange={() => change(i, { enabled: !pair.enabled })} class="size-3 accent-primary" />
      <select value={pair.valueType ?? 'text'} onchange={(e) => change(i, { valueType: (e.currentTarget as HTMLSelectElement).value as 'text' | 'file', value: '', fileName: '' })} class="h-7 rounded border border-border bg-surface-inset px-2 text-xs text-text focus:outline-none focus:border-border-focus">
        <option value="text">Text</option><option value="file">File</option>
      </select>
      <input value={pair.name} oninput={(e) => change(i, { name: e.currentTarget.value })} placeholder="Key" class="h-7 rounded bg-transparent px-2 font-mono text-xs text-text placeholder:text-placeholder focus:bg-surface-highlight focus:outline-none" />
      {#if pair.valueType === 'file'}
        <label class="flex h-7 cursor-pointer items-center gap-2 rounded border border-dashed border-border px-2 text-xs text-text-subtle hover:bg-surface-highlight">
          <Icon name="file" size={13} /><span class="truncate">{pair.fileName || 'Select file'}</span>
          <input type="file" class="sr-only" onchange={(e) => pick(i, e.currentTarget.files)} />
        </label>
      {:else}
        <input value={pair.value} oninput={(e) => change(i, { value: e.currentTarget.value })} placeholder="Value" class="h-7 rounded bg-transparent px-2 font-mono text-xs text-text placeholder:text-placeholder focus:bg-surface-highlight focus:outline-none" />
      {/if}
      <button type="button" onclick={() => remove(i)} class="flex size-6 items-center justify-center rounded text-text-subtlest opacity-60 hover:bg-danger/10 hover:text-danger" aria-label="Remove row"><Icon name="x" size={12} /></button>
    </div>
  {/each}
</div>
