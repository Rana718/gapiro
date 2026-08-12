<!--
  ResponseBody - Renders response body with syntax detection.
-->
<script lang="ts">
  import { isJSON, isHTML, isXML, prettyJSON, languageFromContentType } from '../../lib/utils';

  interface Props {
    body: string;
    contentType: string;
    pretty: boolean;
  }

  let { body, contentType, pretty }: Props = $props();

  const displayBody = $derived(() => {
    if (!body) return '';
    if (pretty && isJSON(contentType)) return prettyJSON(body);
    return body;
  });

  const language = $derived(languageFromContentType(contentType));
</script>

<div class="relative w-full h-full overflow-auto">
  {#if !body}
    <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
      Empty response
    </div>
  {:else}
    <pre class="w-full h-full p-3 m-0 text-xs font-mono text-text leading-relaxed
      whitespace-pre-wrap break-words overflow-auto bg-surface-inset">{displayBody()}</pre>
    <!-- Language badge -->
    <span class="absolute top-2 right-4 px-1.5 py-0.5 text-[9px] uppercase tracking-wider
      text-text-subtlest bg-surface-highlight rounded font-semibold pointer-events-none">
      {language}
    </span>
  {/if}
</div>
