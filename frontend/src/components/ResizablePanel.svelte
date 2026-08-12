<!--
  ResizablePanel - GPU-accelerated resizable split layout.
  Used for sidebar/content and request/response splits.
-->
<script lang="ts">
  import type { Snippet } from 'svelte';
  import { clamp } from '../lib/utils';

  interface Props {
    direction?: 'horizontal' | 'vertical';
    initialRatio?: number;
    minFirst?: number;
    minSecond?: number;
    storageKey?: string;
    first: Snippet;
    second: Snippet;
  }

  let {
    direction = 'horizontal',
    initialRatio = 0.5,
    minFirst = 200,
    minSecond = 200,
    storageKey,
    first,
    second,
  }: Props = $props();

  let container: HTMLElement;
  let dragging = $state(false);
  let ratio = $state(initialRatio);

  // Restore from localStorage if key provided
  if (storageKey) {
    const saved = localStorage.getItem(`gapiro:resize:${storageKey}`);
    if (saved) ratio = parseFloat(saved);
  }

  function onPointerDown(e: PointerEvent) {
    e.preventDefault();
    dragging = true;
    const target = e.currentTarget as HTMLElement;
    target.setPointerCapture(e.pointerId);
  }

  function onPointerMove(e: PointerEvent) {
    if (!dragging || !container) return;
    const rect = container.getBoundingClientRect();
    const total = direction === 'horizontal' ? rect.width : rect.height;
    const pos = direction === 'horizontal'
      ? e.clientX - rect.left
      : e.clientY - rect.top;

    const minRatio = minFirst / total;
    const maxRatio = 1 - (minSecond / total);
    ratio = clamp(pos / total, minRatio, maxRatio);
  }

  function onPointerUp(e: PointerEvent) {
    dragging = false;
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
    if (storageKey) {
      localStorage.setItem(`gapiro:resize:${storageKey}`, ratio.toString());
    }
  }
</script>

<div
  bind:this={container}
  class="flex w-full h-full overflow-hidden gpu-layer {direction === 'horizontal' ? 'flex-row' : 'flex-col'}"
>
  <!-- First panel -->
  <div
    class="overflow-hidden gpu-layer"
    style="{direction === 'horizontal' ? 'width' : 'height'}: {ratio * 100}%"
  >
    {@render first()}
  </div>

  <!-- Resize handle -->
  <div
    class="flex-shrink-0 flex items-center justify-center group z-10
      {direction === 'horizontal' ? 'w-[5px] cursor-col-resize' : 'h-[5px] cursor-row-resize'}
      {dragging ? 'bg-accent/20' : 'hover:bg-accent/10'}
      transition-colors duration-100"
    onpointerdown={onPointerDown}
    onpointermove={onPointerMove}
    onpointerup={onPointerUp}
    role="separator"
    aria-orientation={direction}
  >
    <div class="
      rounded-full bg-border-default group-hover:bg-accent transition-colors duration-100
      {direction === 'horizontal' ? 'w-[2px] h-8' : 'h-[2px] w-8'}
    "></div>
  </div>

  <!-- Second panel -->
  <div
    class="overflow-hidden gpu-layer flex-1"
  >
    {@render second()}
  </div>
</div>
