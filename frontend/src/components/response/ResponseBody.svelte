<!--
  ResponseBody - Response viewer with line numbers and JSON syntax highlighting.
  Matches Postman's dark editor look with proper line numbers.
-->
<script lang="ts">
  import { isJSON, isHTML, isXML, languageFromContentType } from '../../lib/utils';
  import { formatJSONAsync } from '../../lib/formatter';

  interface Props {
    body: string;
    contentType: string;
    pretty: boolean;
    searchQuery?: string;
  }

  let { body, contentType, pretty, searchQuery = '' }: Props = $props();

  let formattedBody = $state('');
  let formatting = $state(false);

  const language = $derived(languageFromContentType(contentType));
  const shouldHighlight = $derived(isJSON(contentType) || isXML(contentType) || isHTML(contentType));

  // Format body off-main-thread for JSON
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

  // Search match count
  const matchCount = $derived(() => {
    if (!searchQuery || !formattedBody) return 0;
    const regex = new RegExp(searchQuery.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    return (formattedBody.match(regex) || []).length;
  });

  // Split into lines for line-numbered display
  const lines = $derived(() => {
    if (!formattedBody) return [];
    return formattedBody.split('\n');
  });

  // Generate highlighted HTML per line
  function highlightLine(line: string): string {
    let html: string;
    if (shouldHighlight && !searchQuery) {
      html = syntaxHighlight(escapeHtml(line), language);
    } else {
      html = escapeHtml(line);
    }

    // Apply search highlighting
    if (searchQuery) {
      const escapedQuery = searchQuery.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      const regex = new RegExp(`(${escapedQuery})`, 'gi');
      html = html.replace(regex, '<mark class="bg-warning/40 rounded-sm px-0.5">$1</mark>');
    }

    return html || '&nbsp;';
  }
</script>

<div class="relative w-full h-full overflow-auto response-viewer">
  {#if !body}
    <div class="flex items-center justify-center h-full text-sm text-text-subtlest">
      Empty response
    </div>
  {:else if formatting}
    <div class="flex items-center justify-center h-full text-xs text-text-subtlest">
      Formatting...
    </div>
  {:else}
    <!-- Search info -->
    {#if searchQuery}
      <div class="sticky top-0 z-10 px-3 py-1 bg-surface-active/90 text-[10px] text-text-subtlest border-b border-border-subtle backdrop-blur-sm">
        {matchCount()} match{matchCount() !== 1 ? 'es' : ''}
      </div>
    {/if}

    <!-- Line-numbered body display -->
    <div class="flex min-h-full">
      <!-- Line numbers gutter -->
      <div class="shrink-0 select-none sticky left-0 bg-surface-inset border-r border-border-subtle z-[1]">
        {#each lines() as _, i}
          <div class="px-3 py-0 text-right text-[11px] leading-[1.65] font-mono text-text-subtlest/60 h-[18px]">
            {i + 1}
          </div>
        {/each}
      </div>

      <!-- Code content -->
      <pre class="response-body flex-1 m-0 p-0 text-xs leading-[1.65] whitespace-pre-wrap break-words
        font-mono text-text selection:bg-primary/20">{#each lines() as line, i}<div class="px-4 h-[18px]">{@html highlightLine(line)}</div>{/each}</pre>
    </div>

    <!-- Language badge -->
    <span class="absolute top-2 right-4 px-1.5 py-0.5 text-[9px] uppercase tracking-wider
      text-text-subtlest bg-surface-active rounded font-semibold pointer-events-none select-none">
      {language}
    </span>
  {/if}
</div>

<script module lang="ts">
  /** Lightweight JSON/XML syntax highlighter (works on already-escaped HTML) */
  function syntaxHighlight(escaped: string, lang: string): string {
    if (lang === 'json') {
      return escaped
        // Keys
        .replace(/(&quot;)((?:[^&]|&(?!quot;))*)(&quot;)\s*:/g,
          '<span class="syn-key">$1$2$3</span>:')
        // String values after colon
        .replace(/:\s*(&quot;)((?:[^&]|&(?!quot;))*)(&quot;)/g,
          ': <span class="syn-str">$1$2$3</span>')
        // Strings in arrays
        .replace(/(\[|,)\s*(&quot;)((?:[^&]|&(?!quot;))*)(&quot;)/g,
          '$1 <span class="syn-str">$2$3$4</span>')
        // Numbers
        .replace(/(:\s*)(-?\d+\.?\d*(?:[eE][+-]?\d+)?)/g,
          '$1<span class="syn-num">$2</span>')
        // Booleans and null
        .replace(/(:\s*)(true|false|null)/g,
          '$1<span class="syn-bool">$2</span>');
    }

    if (lang === 'xml' || lang === 'html') {
      return escaped
        .replace(/(&lt;\/?)([\w-]+)/g, '$1<span class="syn-tag">$2</span>')
        .replace(/([\w-]+)(=)/g, '<span class="syn-attr">$1</span>$2')
        .replace(/(=)(&quot;[^&]*&quot;)/g, '$1<span class="syn-str">$2</span>');
    }

    return escaped;
  }

  function escapeHtml(str: string): string {
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }
</script>

<style>
  .response-viewer {
    background: var(--color-surface-inset);
  }
  :global(.response-body .syn-key) { color: #c4b5fd; }
  :global(.response-body .syn-str) { color: #6ee7b7; }
  :global(.response-body .syn-num) { color: #fbbf24; }
  :global(.response-body .syn-bool) { color: #f87171; }
  :global(.response-body .syn-tag) { color: #f87171; }
  :global(.response-body .syn-attr) { color: #fbbf24; }
</style>
