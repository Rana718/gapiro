<!--
  TitleBar - Custom window title bar for borderless Wails window.
  Draggable for window movement, shows app name.
-->
<script lang="ts">
  interface Props {
    title?: string;
  }

  let { title = 'Gapiro' }: Props = $props();
</script>

<div
  class="flex items-center h-9 px-3 bg-surface-0 border-b border-border-subtle
    select-none shrink-0"
  style="--wails-draggable: drag"
>
  <!-- App icon + name -->
  <div class="flex items-center gap-2">
    <svg class="w-4 h-4 text-accent" viewBox="0 0 24 24" fill="currentColor">
      <path d="M13 3L4 14h7l-2 7 9-11h-7l2-7z"/>
    </svg>
    <span class="text-xs font-semibold text-text-secondary">{title}</span>
  </div>

  <!-- Spacer (draggable area) -->
  <div class="flex-1"></div>
</div>
