<!--
  ResponseBody - Renders response body with CodeMirror syntax highlighting.
-->
<script lang="ts">
  import { isJSON, prettyJSON, languageFromContentType } from '../../lib/utils';
  import CodeEditor from '../core/CodeEditor.svelte';

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

<div class="w-full h-full overflow-hidden">
  {#if !body}
    <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
      Empty response
    </div>
  {:else}
    <CodeEditor
      value={displayBody()}
      readonly
      {language}
      lineNums={true}
    />
  {/if}
</div>
