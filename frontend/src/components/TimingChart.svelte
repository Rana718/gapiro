<!--
  TimingChart - Visual waterfall chart showing request timing breakdown.
-->
<script lang="ts">
  import { formatDuration } from '../lib/utils';

  interface Props {
    dns: number;
    connect: number;
    tls: number;
    ttfb: number;
    total: number;
  }

  let { dns, connect, tls, ttfb, total }: Props = $props();

  const segments = $derived(() => {
    const download = Math.max(0, total - ttfb);
    return [
      { label: 'DNS Lookup', value: dns, color: 'bg-cyan-500' },
      { label: 'TCP Connect', value: connect, color: 'bg-green-500' },
      { label: 'TLS Handshake', value: tls, color: 'bg-purple-500' },
      { label: 'Time to First Byte', value: ttfb, color: 'bg-amber-500' },
      { label: 'Content Download', value: download, color: 'bg-blue-500' },
      { label: 'Total', value: total, color: 'bg-text-primary' },
    ];
  });
</script>

<div class="p-4 flex flex-col gap-3">
  <!-- Waterfall bars -->
  {#each segments() as seg (seg.label)}
    <div class="flex flex-col gap-1">
      <div class="flex items-center justify-between">
        <span class="text-[11px] text-text-secondary">{seg.label}</span>
        <span class="text-[11px] font-mono text-text-primary">{formatDuration(seg.value)}</span>
      </div>
      <div class="h-2 bg-surface-3 rounded-full overflow-hidden">
        <div
          class="h-full rounded-full {seg.color} transition-all duration-300"
          style="width: {total > 0 ? Math.max((seg.value / total) * 100, seg.value > 0 ? 2 : 0) : 0}%"
        ></div>
      </div>
    </div>
  {/each}
</div>
