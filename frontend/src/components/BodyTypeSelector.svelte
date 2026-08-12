<!--
  BodyTypeSelector - Segmented control for selecting request body type.
-->
<script lang="ts">
  import type { BodyType } from '../lib/types';

  interface Props {
    bodyType: BodyType;
    onchange: (type: BodyType) => void;
  }

  let { bodyType, onchange }: Props = $props();

  const options: { id: BodyType; label: string }[] = [
    { id: 'none', label: 'None' },
    { id: 'json', label: 'JSON' },
    { id: 'text', label: 'Text' },
    { id: 'form-urlencoded', label: 'Form URL' },
    { id: 'form-data', label: 'Multipart' },
    { id: 'binary', label: 'Binary' },
  ];
</script>

<div class="flex items-center gap-1 px-3 py-2 border-b border-border-subtle">
  {#each options as opt (opt.id)}
    <button
      onclick={() => onchange(opt.id)}
      class="px-2.5 py-1 text-[11px] font-medium rounded cursor-pointer
        transition-colors duration-75
        {bodyType === opt.id
          ? 'bg-accent/15 text-accent border border-accent/30'
          : 'text-text-muted hover:text-text-secondary hover:bg-surface-2 border border-transparent'}"
    >
      {opt.label}
    </button>
  {/each}
</div>
