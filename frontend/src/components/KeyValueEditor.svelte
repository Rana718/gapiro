<!--
  KeyValueEditor - Reusable key-value pair editor with add/remove rows.
  Used for headers, query params, form data, etc.
-->
<script lang="ts">
  import type { KeyValue } from '../lib/types';
  import { emptyKV } from '../lib/utils';

  interface Props {
    items: KeyValue[];
    keyPlaceholder?: string;
    valuePlaceholder?: string;
    onchange: (items: KeyValue[]) => void;
  }

  let { items, keyPlaceholder = 'Key', valuePlaceholder = 'Value', onchange }: Props = $props();

  function toggle(index: number) {
    items[index].enabled = !items[index].enabled;
    onchange(items);
  }

  function updateKey(index: number, value: string) {
    items[index].key = value;
    ensureEmptyRow();
    onchange(items);
  }

  function updateValue(index: number, value: string) {
    items[index].value = value;
    ensureEmptyRow();
    onchange(items);
  }

  function removeRow(index: number) {
    if (items.length <= 1) return;
    items.splice(index, 1);
    onchange(items);
  }

  function ensureEmptyRow() {
    const last = items[items.length - 1];
    if (last && (last.key !== '' || last.value !== '')) {
      items.push(emptyKV());
    }
  }
</script>

<div class="flex flex-col gap-0 w-full">
  <!-- Header -->
  <div class="grid grid-cols-[28px_1fr_1fr_28px] gap-1 px-2 py-1 text-[10px] uppercase tracking-wider text-text-muted font-semibold">
    <span></span>
    <span>{keyPlaceholder}</span>
    <span>{valuePlaceholder}</span>
    <span></span>
  </div>

  <!-- Rows -->
  {#each items as item, i (item.id)}
    <div class="grid grid-cols-[28px_1fr_1fr_28px] gap-1 px-2 py-0.5 group items-center
      hover:bg-surface-2/50 transition-colors duration-75">
      <!-- Toggle checkbox -->
      <label class="flex items-center justify-center cursor-pointer">
        <input
          type="checkbox"
          checked={item.enabled}
          onchange={() => toggle(i)}
          class="w-3.5 h-3.5 rounded border-border-default bg-surface-2
            checked:bg-accent checked:border-accent cursor-pointer
            focus:ring-1 focus:ring-accent/50"
        />
      </label>

      <!-- Key input -->
      <input
        type="text"
        value={item.key}
        placeholder={keyPlaceholder}
        oninput={(e) => updateKey(i, (e.target as HTMLInputElement).value)}
        class="w-full bg-transparent border-0 px-2 py-1 text-xs text-text-primary
          placeholder:text-text-muted/50 focus:outline-none focus:bg-surface-2 rounded
          {!item.enabled ? 'opacity-40' : ''}"
      />

      <!-- Value input -->
      <input
        type="text"
        value={item.value}
        placeholder={valuePlaceholder}
        oninput={(e) => updateValue(i, (e.target as HTMLInputElement).value)}
        class="w-full bg-transparent border-0 px-2 py-1 text-xs text-text-primary
          placeholder:text-text-muted/50 focus:outline-none focus:bg-surface-2 rounded
          {!item.enabled ? 'opacity-40' : ''}"
      />

      <!-- Delete button -->
      <button
        onclick={() => removeRow(i)}
        class="flex items-center justify-center w-5 h-5 rounded text-text-muted
          opacity-0 group-hover:opacity-100 hover:text-error hover:bg-error/10
          transition-all duration-75 cursor-pointer"
        aria-label="Remove row"
      >
        <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>
  {/each}
</div>
