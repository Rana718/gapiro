<!--
  StatusBadge - Displays HTTP status code with appropriate color.
-->
<script lang="ts">
  import { statusColor } from '../lib/utils';

  interface Props {
    status: number;
    statusText: string;
  }

  let { status, statusText }: Props = $props();

  const label = $derived(() => {
    const parts = statusText.split(' ');
    // Remove the leading "200" from "200 OK" since we show it separately
    return parts.length > 1 ? parts.slice(1).join(' ') : statusText;
  });
</script>

<div class="flex items-center gap-1.5">
  <span class="px-1.5 py-0.5 text-[11px] font-bold rounded {statusColor(status)} bg-current/10">
    {status}
  </span>
  <span class="text-[11px] text-text-secondary">{label()}</span>
</div>
