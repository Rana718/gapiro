<!--
  SidebarLayout - Fixed-width sidebar with resizable handle.
  Floats on narrow viewports (<600px).
-->
<script lang="ts">
  import type { Snippet } from 'svelte';
  import { clamp } from '../../lib/utils';

  interface Props {
    width: number;
    onWidthChange: (w: number) => void;
    hidden?: boolean;
    onHiddenChange?: (h: boolean) => void;
    minWidth?: number;
    defaultWidth?: number;
    sidebar: Snippet;
    children: Snippet;
    class?: string;
  }

  let {
    width,
    onWidthChange,
    hidden = false,
    onHiddenChange,
    minWidth = 50,
    defaultWidth = 250,
    sidebar,
    children,
    class: className = '',
  }: Props = $props();

  let containerEl: HTMLElement;
  let containerWidth = $state(0);
  let dragging = $state(false);
  let startX = $state(0);
  let startWidth = $state(0);

  const floating = $derived(containerWidth > 0 && containerWidth <= 600);
  const sideWidth = $derived(hidden ? 0 : width);

  const gridStyle = $derived(
    `grid-template-columns: ${sideWidth}px 0px 1fr; grid-template-rows: 1fr;`
  );

  function onPointerDown(e: PointerEvent) {
    e.preventDefault();
    dragging = true;
    startX = e.clientX;
    startWidth = width;
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
  }

  function onPointerMove(e: PointerEvent) {
    if (!dragging) return;
    const newWidth = startWidth + (e.clientX - startX);
    if (newWidth < minWidth) {
      onHiddenChange?.(true);
      onWidthChange(defaultWidth);
    } else {
      if (hidden) onHiddenChange?.(false);
      onWidthChange(Math.min(newWidth, containerWidth * 0.5));
    }
  }

  function onPointerUp(e: PointerEvent) {
    dragging = false;
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
  }

  function onDblClick() {
    onWidthChange(defaultWidth);
    onHiddenChange?.(false);
  }

  function observeContainer(node: HTMLElement) {
    containerEl = node;
    const ro = new ResizeObserver(entries => {
      containerWidth = entries[0]?.contentRect.width ?? 0;
    });
    ro.observe(node);
    return { destroy() { ro.disconnect(); } };
  }
</script>

{#if floating}
  <div use:observeContainer class="relative w-full h-full {className}">
    <!-- Floating sidebar overlay -->
    {#if !hidden}
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="absolute inset-0 z-20 bg-black/40"
        onclick={() => onHiddenChange?.(true)}
        onkeydown={(e) => { if (e.key === 'Escape') onHiddenChange?.(true); }}
      ></div>
      <div
        class="absolute top-0 left-0 bottom-0 z-30 w-[320px] gpu
          animate-[slideIn_150ms_ease-out]"
      >
        {@render sidebar()}
      </div>
    {/if}
    {@render children()}
  </div>
{:else}
  <div
    use:observeContainer
    class="grid w-full h-full {className} {!dragging ? 'transition-grid' : ''}"
    style={gridStyle}
  >
    <div class="overflow-hidden min-w-0">
      {@render sidebar()}
    </div>
    <div
      class="z-10 w-[6px] -mx-[3px] cursor-col-resize flex items-center justify-center
        hover:bg-primary/10 {dragging ? 'bg-primary/15' : ''}
        transition-colors duration-75 group select-none"
      onpointerdown={onPointerDown}
      onpointermove={onPointerMove}
      onpointerup={onPointerUp}
      ondblclick={onDblClick}
      role="separator"
      aria-orientation="vertical"
    >
      {#if dragging}
        <div class="fixed inset-0 cursor-col-resize"></div>
      {/if}
    </div>
    <div class="overflow-hidden min-w-0">
      {@render children()}
    </div>
  </div>
{/if}

<style>
  @keyframes slideIn {
    from { transform: translateX(-100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
</style>
