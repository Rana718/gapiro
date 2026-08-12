<!--
  SplitLayout - GPU-composited resizable split panel.
  Mimics Yaak's CSS Grid-based SplitLayout.
-->
<script lang="ts">
  import type { Snippet } from 'svelte';
  import { clamp } from '../../lib/utils';

  interface Props {
    layout?: 'horizontal' | 'vertical' | 'responsive';
    storageKey: string;
    defaultRatio?: number;
    minPx?: number;
    first: Snippet;
    second: Snippet;
    class?: string;
  }

  let {
    layout = 'responsive',
    storageKey,
    defaultRatio = 0.5,
    minPx = 100,
    first,
    second,
    class: className = '',
  }: Props = $props();

  let container: HTMLElement;
  let dragging = $state(false);
  let containerWidth = $state(0);
  let containerHeight = $state(0);

  // Load persisted ratio
  let ratio = $state(defaultRatio);
  const savedRatio = localStorage.getItem(`gapiro:split:${storageKey}`);
  if (savedRatio) ratio = parseFloat(savedRatio);

  // Determine vertical based on responsive breakpoint
  const vertical = $derived(
    layout === 'vertical' || (layout === 'responsive' && containerWidth > 0 && containerWidth < 600)
  );

  // Clamp ratio based on min pixel constraints
  const renderedRatio = $derived(() => {
    const total = vertical ? containerHeight : containerWidth;
    if (total <= 0) return ratio;
    const minRatio = Math.min(0.9, minPx / total);
    const maxRatio = 1 - minRatio;
    return clamp(ratio, minRatio, maxRatio);
  });

  // CSS Grid template
  const gridStyle = $derived(() => {
    const r = renderedRatio();
    if (vertical) {
      return `grid-template-rows: ${1 - r}fr 0px ${r}fr; grid-template-columns: 1fr;`;
    }
    return `grid-template-columns: ${1 - r}fr 0px ${r}fr; grid-template-rows: 1fr;`;
  });

  function onPointerDown(e: PointerEvent) {
    e.preventDefault();
    dragging = true;
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
  }

  function onPointerMove(e: PointerEvent) {
    if (!dragging || !container) return;
    const rect = container.getBoundingClientRect();
    const total = vertical ? rect.height : rect.width;
    const pos = vertical ? (e.clientY - rect.top) : (e.clientX - rect.left);
    const minR = minPx / total;
    const maxR = 1 - minR;
    // Ratio is "second panel" fraction, so we use (total - pos) / total
    ratio = clamp(1 - pos / total, minR, maxR);
  }

  function onPointerUp(e: PointerEvent) {
    dragging = false;
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
    localStorage.setItem(`gapiro:split:${storageKey}`, ratio.toString());
  }

  function onDblClick() {
    ratio = defaultRatio;
    localStorage.setItem(`gapiro:split:${storageKey}`, ratio.toString());
  }

  // Observe container size
  function observeSize(node: HTMLElement) {
    container = node;
    const ro = new ResizeObserver(entries => {
      for (const entry of entries) {
        containerWidth = entry.contentRect.width;
        containerHeight = entry.contentRect.height;
      }
    });
    ro.observe(node);
    return { destroy() { ro.disconnect(); } };
  }
</script>

<div
  use:observeSize
  class="grid w-full h-full overflow-hidden gpu {className} {!dragging ? 'transition-grid' : ''}"
  style={gridStyle()}
>
  <!-- First slot -->
  <div class="overflow-hidden min-w-0 min-h-0">
    {@render first()}
  </div>

  <!-- Resize handle -->
  <div
    class="z-10 flex items-center justify-center group select-none
      {vertical ? 'cursor-row-resize h-[6px] -my-[3px]' : 'cursor-col-resize w-[6px] -mx-[3px]'}
      hover:bg-primary/10 {dragging ? 'bg-primary/15' : ''}
      transition-colors duration-75"
    onpointerdown={onPointerDown}
    onpointermove={onPointerMove}
    onpointerup={onPointerUp}
    ondblclick={onDblClick}
    role="separator"
    aria-orientation={vertical ? 'horizontal' : 'vertical'}
  >
    {#if dragging}
      <div class="fixed inset-0 {vertical ? 'cursor-row-resize' : 'cursor-col-resize'}"></div>
    {/if}
  </div>

  <!-- Second slot -->
  <div class="overflow-hidden min-w-0 min-h-0">
    {@render second()}
  </div>
</div>
