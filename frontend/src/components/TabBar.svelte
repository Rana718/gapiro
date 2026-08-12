<!--
  TabBar - Reusable tab navigation component.
  Used for request tabs (params, headers, body, auth, settings)
  and response tabs (body, headers, timing, info).
-->
<script lang="ts">
  interface Tab {
    id: string;
    label: string;
    badge?: number | string;
  }

  interface Props {
    tabs: Tab[];
    active: string;
    onchange: (id: string) => void;
    size?: 'sm' | 'md';
  }

  let { tabs, active, onchange, size = 'sm' }: Props = $props();
</script>

<div class="flex items-center gap-0.5 border-b border-border-subtle px-2" role="tablist">
  {#each tabs as tab (tab.id)}
    <button
      role="tab"
      aria-selected={active === tab.id}
      class="relative px-3 py-1.5 text-{size === 'sm' ? 'xs' : 'sm'} font-medium rounded-t
        transition-colors duration-100 cursor-pointer
        {active === tab.id
          ? 'text-text-primary'
          : 'text-text-muted hover:text-text-secondary'}"
      onclick={() => onchange(tab.id)}
    >
      <span class="flex items-center gap-1.5">
        {tab.label}
        {#if tab.badge}
          <span class="inline-flex items-center justify-center min-w-[16px] h-4 px-1
            text-[10px] font-semibold rounded-full bg-accent/20 text-accent">
            {tab.badge}
          </span>
        {/if}
      </span>
      {#if active === tab.id}
        <div class="absolute bottom-0 left-1 right-1 h-[2px] bg-accent rounded-t"></div>
      {/if}
    </button>
  {/each}
</div>
