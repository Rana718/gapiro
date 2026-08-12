<!--
  ResponseTimeline - Visual waterfall chart showing request timing breakdown.
  Matches Yaak's HttpResponseTimeline visual mode.
-->
<script lang="ts">
  import { formatDuration } from '../../lib/utils';

  interface Props {
    dns: number;
    connect: number;
    tls: number;
    ttfb: number;
    total: number;
  }

  let { dns, connect, tls, ttfb, total }: Props = $props();

  const download = $derived(Math.max(0, total - ttfb));

  const segments = $derived([
    { label: 'DNS Lookup', value: dns, color: 'bg-cyan-500', icon: '🌐' },
    { label: 'TCP Connection', value: connect, color: 'bg-green-500', icon: '🔗' },
    { label: 'TLS Handshake', value: tls, color: 'bg-purple-500', icon: '🔒' },
    { label: 'Time to First Byte', value: ttfb, color: 'bg-amber-500', icon: '⏱️' },
    { label: 'Content Download', value: download, color: 'bg-blue-500', icon: '⬇️' },
  ]);
</script>

<div class="p-4 flex flex-col gap-4 overflow-y-auto h-full">
  <!-- Waterfall chart -->
  {#each segments as seg (seg.label)}
    <div class="flex flex-col gap-1.5">
      <div class="flex items-center justify-between">
        <span class="flex items-center gap-1.5 text-xs text-text-subtle">
          <span class="text-[11px]">{seg.icon}</span>
          {seg.label}
        </span>
        <span class="text-xs font-mono text-text font-medium">{formatDuration(seg.value)}</span>
      </div>
      <div class="h-2 bg-surface-active rounded-full overflow-hidden">
        <div
          class="h-full rounded-full {seg.color} transition-all duration-300 ease-out"
          style="width: {total > 0 ? Math.max(seg.value > 0 ? 2 : 0, (seg.value / total) * 100) : 0}%"
        ></div>
      </div>
    </div>
  {/each}

  <!-- Total -->
  <div class="flex items-center justify-between pt-2 mt-2 border-t border-border-subtle">
    <span class="text-xs font-semibold text-text">Total</span>
    <span class="text-xs font-mono font-bold text-text">{formatDuration(total)}</span>
  </div>
</div>
