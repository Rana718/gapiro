<!--
  PairEditor - Core name/value pair editor with enable/disable, add/remove.
  Used for headers, query params, form data.
-->
<script lang="ts">
  import type { Pair } from '../../lib/types';
  import { emptyPair } from '../../lib/utils';
  import Icon from './Icon.svelte';

  interface Props {
    pairs: Pair[];
    namePlaceholder?: string;
    valuePlaceholder?: string;
    onchange: (pairs: Pair[]) => void;
    class?: string;
  }

  let {
    pairs,
    namePlaceholder = 'Name',
    valuePlaceholder = 'Value',
    onchange,
    class: className = '',
  }: Props = $props();

  function toggle(i: number) {
    pairs[i].enabled = !pairs[i].enabled;
    onchange(pairs);
  }

  function updateName(i: number, val: string) {
    pairs[i].name = val;
    ensureEmpty();
    onchange(pairs);
  }

  function updateValue(i: number, val: string) {
    pairs[i].value = val;
    ensureEmpty();
    onchange(pairs);
  }

  function remove(i: number) {
    if (pairs.length <= 1) { pairs[0].name = ''; pairs[0].value = ''; onchange(pairs); return; }
    pairs.splice(i, 1);
    onchange(pairs);
  }

  function ensureEmpty() {
    const last = pairs[pairs.length - 1];
    if (last && (last.name !== '' || last.value !== '')) {
      pairs.push(emptyPair());
    }
  }
</script>

<div class="flex flex-col w-full overflow-y-auto {className}">
  <!-- Header row -->
  <div class="grid grid-cols-[26px_1fr_1fr_26px] gap-1 px-3 py-2
    text-[10px] uppercase tracking-wider text-text-subtlest font-semibold border-b border-border-subtle">
    <span></span>
    <span>{namePlaceholder}</span>
    <span>{valuePlaceholder}</span>
    <span></span>
  </div>

  <div class="px-3 py-2 text-[11px] text-text-subtlest border-b border-border-subtle bg-surface-inset">
    Enabled rows are sent with the request. Add values inline; empty rows are ignored.
  </div>

  <!-- Rows -->
  {#each pairs as pair, i (pair.id)}
    <div class="grid grid-cols-[26px_1fr_1fr_26px] gap-1 px-3 py-[3px] items-center group
      hover:bg-surface-highlight/50 transition-colors duration-50">
      <!-- Checkbox -->
      <label class="flex items-center justify-center">
        <input
          type="checkbox"
          checked={pair.enabled}
          onchange={() => toggle(i)}
          class="w-3 h-3 rounded-sm border border-border bg-surface-inset
            checked:bg-primary checked:border-primary accent-primary cursor-pointer"
        />
      </label>

      <!-- Name -->
      <input
        type="text"
        value={pair.name}
        placeholder={namePlaceholder}
        readonly={pair.readOnly}
        oninput={(e) => updateName(i, (e.target as HTMLInputElement).value)}
        class="w-full bg-transparent px-2 py-1 text-xs text-text rounded
          placeholder:text-placeholder focus:outline-none focus:bg-surface-highlight
          border-0 font-mono
          {!pair.enabled ? 'opacity-40' : ''}
          {pair.readOnly ? 'text-primary italic' : ''}"
        spellcheck="false"
      />

      <!-- Value -->
      <input
        type="text"
        value={pair.value}
        placeholder={valuePlaceholder}
        oninput={(e) => updateValue(i, (e.target as HTMLInputElement).value)}
        class="w-full bg-transparent px-2 py-1 text-xs text-text rounded
          placeholder:text-placeholder focus:outline-none focus:bg-surface-highlight
          border-0 font-mono
          {!pair.enabled ? 'opacity-40' : ''}"
        spellcheck="false"
      />

      <!-- Delete -->
      <button
        onclick={() => remove(i)}
        class="flex items-center justify-center w-5 h-5 rounded
          text-text-subtlest opacity-60 group-hover:opacity-100
          hover:text-danger hover:bg-danger/10 transition-all duration-50"
        aria-label="Remove"
      >
        <Icon name="x" size={12} />
      </button>
    </div>
  {/each}
</div>
