<!--
  ResponseBody - Response viewer with JSON syntax highlighting.
  Uses lightweight regex-based highlighting (no CodeMirror for read-only).
  Supports search with match highlighting.
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

  // Truncate for display
  const MAX_DISPLAY = 800_000;
  const displayText = $derived(() => {
    if (formattedBody.length > MAX_DISPLAY) {
      return formattedBody.slice(0, MAX_DISPLAY) + `\n\n--- Truncated (${Math.round(formattedBody.length / 1024)} KB total) ---`;
    }
    return formattedBody;
  });

  // Generate highlighted HTML
  const highlightedHtml = $derived(() => {
    const text = displayText();
    if (!text) return '';

    let html: string;
    if (shouldHighlight && !searchQuery) {
      html = syntaxHighlight(text, language);
    } else {
      html = escapeHtml(text);
    }

    // Apply search highlighting on top
    if (searchQuery) {
      const escapedQuery = searchQuery.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      const regex = new RegExp(`(${escapeHtml(escapedQuery)})`, 'gi');
      html = html.replace(regex, '<mark class="bg-warning/40 rounded-sm">$1</mark>');
    }

    return html;
  });
</script>

<div class="relative w-full h-full overflow-auto bg-background">
  {#if !body}
    <div class="flex items-center justify-center h-full text-sm text-muted-foreground">
      Empty response
    </div>
  {:else if formatting}
    <div class="flex items-center justify-center h-full text-xs text-muted-foreground">
      Formatting...
    </div>
  {:else}
    <!-- Search info -->
    {#if searchQuery}
      <div class="sticky top-0 z-10 px-3 py-1 bg-muted/90 text-[10px] text-muted-foreground border-b border-border">
        {matchCount()} match{matchCount() !== 1 ? 'es' : ''}
      </div>
    {/if}

    <pre class="response-body w-full min-h-full p-3 m-0 text-xs leading-[1.65] whitespace-pre-wrap break-words
      font-mono text-foreground selection:bg-primary/20">{@html highlightedHtml()}</pre>

    <!-- Language badge -->
    <span class="absolute top-2 right-4 px-1.5 py-0.5 text-[9px] uppercase tracking-wider
      text-muted-foreground bg-muted rounded font-semibold pointer-events-none select-none">
      {language}
    </span>
  {/if}
</div>

<script context="module" lang="ts">
  /** Lightweight JSON/XML syntax highlighter */
  function syntaxHighlight(text: string, lang: string): string {
    const escaped = escapeHtml(text);

    if (lang === 'json') {
      return escaped
        // Strings (keys and values)
        .replace(/(&quot;)((?:[^&]|&(?!quot;))*)(&quot;)\s*:/g,
          '<span class="syn-key">$1$2$3</span>:')
        .replace(/:\s*(&quot;)((?:[^&]|&(?!quot;))*)(&quot;)/g,
          ': <span class="syn-str">$1$2$3</span>')
        // Standalone strings in arrays
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
        // Tags
        .replace(/(&lt;\/?)([\w-]+)/g, '$1<span class="syn-tag">$2</span>')
        // Attributes
        .replace(/([\w-]+)(=)/g, '<span class="syn-attr">$1</span>$2')
        // Attribute values
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
  :global(.response-body .syn-key) { color: #7c6ef6; }
  :global(.response-body .syn-str) { color: #30a46c; }
  :global(.response-body .syn-num) { color: #f5a623; }
  :global(.response-body .syn-bool) { color: #e5484d; }
  :global(.response-body .syn-tag) { color: #e5484d; }
  :global(.response-body .syn-attr) { color: #f5a623; }
</style>
