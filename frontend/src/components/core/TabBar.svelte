<!--
  TabBar - Horizontal tab bar with count badges and dropdown options.
  Reused for both request and response tab navigation.
-->
<script lang="ts">
  interface Tab {
    id: string;
    label: string;
    badge?: number | string;
    options?: { id: string; label: string }[];
  }

  interface Props {
    tabs: Tab[];
    active: string;
    onchange: (id: string) => void;
    class?: string;
  }

  let { tabs, active, onchange, class: className = '' }: Props = $props();
</script>

<div class="flex items-end gap-0 border-b border-border-subtle shrink-0 pl-3 {className}" role="tablist">
  {#each tabs as tab (tab.id)}
    <button
      role="tab"
      aria-selected={active === tab.id}
      class="relative px-3 py-1.5 text-xs font-medium transition-colors duration-75
        {active === tab.id
          ? 'text-text'
          : 'text-text-subtlest hover:text-text-subtle'}"
      onclick={() => onchange(tab.id)}
    >
      <span class="flex items-center gap-1.5">
        {tab.label}
        {#if tab.badge !== undefined && tab.badge !== 0}
          <span class="inline-flex items-center justify-center min-w-[15px] h-[15px] px-[4px]
            text-[9px] font-bold rounded-full bg-primary/15 text-primary">
            {tab.badge}
          </span>
        {/if}
      </span>
      <!-- Active indicator -->
      {#if active === tab.id}
        <div class="absolute bottom-[-1px] left-2 right-2 h-[2px] bg-primary rounded-t-full"></div>
      {/if}
    </button>
  {/each}
</div>
