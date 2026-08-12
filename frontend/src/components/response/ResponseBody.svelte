<!--
  ResponseBody - High-performance response body viewer.
  Uses Web Worker for JSON formatting (non-blocking).
  Plain <pre> for rendering (no CodeMirror overhead).
-->
<script lang="ts">
  import { isJSON, languageFromContentType } from '../../lib/utils';
  import { formatJSONAsync } from '../../lib/formatter';

  interface Props {
    body: string;
    contentType: string;
    pretty: boolean;
  }

  let { body, contentType, pretty }: Props = $props();

  let formattedBody = $state('');
  let formatting = $state(false);

  const language = $derived(languageFromContentType(contentType));

  // Format body off-main-thread when props change
  $effect(() => {
    if (!body) {
      formattedBody = '';
      return;
    }

    if (pretty && isJSON(contentType)) {
      formatting = true;
      formatJSONAsync(body).then(result => {
        formattedBody = result;
        formatting = false;
      });
    } else {
      formattedBody = body;
      formatting = false;
    }
  });

  // Truncate extremely large bodies for display perf
  const MAX_DISPLAY = 1_000_000; // 1MB
  const displayText = $derived(() => {
    if (formattedBody.length > MAX_DISPLAY) {
      return formattedBody.slice(0, MAX_DISPLAY) + `\n\n─── Truncated (${Math.round(formattedBody.length / 1024)} KB total) ───`;
    }
    return formattedBody;
  });
</script>

<div class="relative w-full h-full overflow-auto bg-surface-inset">
  {#if !body}
    <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
      Empty response
    </div>
  {:else if formatting}
    <div class="flex items-center justify-center h-full text-xs text-text-subtlest">
      Formatting...
    </div>
  {:else}
    <pre class="w-full min-h-full p-3 m-0 text-xs leading-relaxed whitespace-pre-wrap break-words
      font-mono text-text selection:bg-primary/20">{displayText()}</pre>
    <span class="absolute top-2 right-4 px-1.5 py-0.5 text-[9px] uppercase tracking-wider
      text-text-subtlest bg-surface-highlight rounded font-semibold pointer-events-none select-none">
      {language}
    </span>
  {/if}
</div>
