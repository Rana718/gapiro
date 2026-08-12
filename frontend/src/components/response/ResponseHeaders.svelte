<!--
  ResponseHeaders - Display response (and request) headers in a table.
-->
<script lang="ts">
  interface Props {
    requestHeaders: Record<string, string>;
    responseHeaders: Record<string, string>;
  }

  let { requestHeaders, responseHeaders }: Props = $props();

  const sortedResponseHeaders = $derived(
    Object.entries(responseHeaders).sort(([a], [b]) => a.localeCompare(b))
  );
  const sortedRequestHeaders = $derived(
    Object.entries(requestHeaders).sort(([a], [b]) => a.localeCompare(b))
  );
</script>

<div class="overflow-y-auto h-full">
  <!-- Response Headers -->
  {#if sortedResponseHeaders.length > 0}
    <div class="px-3 pt-3 pb-1">
      <h4 class="text-[10px] uppercase tracking-wider text-text-subtlest font-semibold mb-1">
        Response Headers ({sortedResponseHeaders.length})
      </h4>
      <div class="border border-border-subtle rounded-md overflow-hidden">
        {#each sortedResponseHeaders as [key, value], i (key)}
          <div class="grid grid-cols-[minmax(120px,0.4fr)_1fr] gap-2 px-3 py-1.5
            {i > 0 ? 'border-t border-border-subtle/50' : ''}
            hover:bg-surface-highlight/50 transition-colors duration-50">
            <span class="text-xs font-medium text-primary font-mono truncate">{key}</span>
            <span class="text-xs text-text font-mono break-all">{value}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Request Headers -->
  {#if sortedRequestHeaders.length > 0}
    <div class="px-3 pt-3 pb-3">
      <h4 class="text-[10px] uppercase tracking-wider text-text-subtlest font-semibold mb-1">
        Request Headers ({sortedRequestHeaders.length})
      </h4>
      <div class="border border-border-subtle rounded-md overflow-hidden">
        {#each sortedRequestHeaders as [key, value], i (key)}
          <div class="grid grid-cols-[minmax(120px,0.4fr)_1fr] gap-2 px-3 py-1.5
            {i > 0 ? 'border-t border-border-subtle/50' : ''}
            hover:bg-surface-highlight/50 transition-colors duration-50">
            <span class="text-xs font-medium text-notice font-mono truncate">{key}</span>
            <span class="text-xs text-text font-mono break-all">{value}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
