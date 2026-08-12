<!--
  CodeEditor - A simple textarea-based code editor with syntax highlighting placeholder.
  Used for request body (JSON, text, etc.) and response body display.
-->
<script lang="ts">
  interface Props {
    value: string;
    onchange?: (value: string) => void;
    readonly?: boolean;
    language?: string;
    placeholder?: string;
    class?: string;
  }

  let { value, onchange, readonly = false, language = 'json', placeholder = '', class: className = '' }: Props = $props();

  function handleInput(e: Event) {
    const target = e.target as HTMLTextAreaElement;
    onchange?.(target.value);
  }

  function handleKeydown(e: KeyboardEvent) {
    // Tab support in textarea
    if (e.key === 'Tab') {
      e.preventDefault();
      const target = e.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;
      const newValue = value.substring(0, start) + '  ' + value.substring(end);
      onchange?.(newValue);
      // Restore cursor position after Svelte updates the value
      requestAnimationFrame(() => {
        target.selectionStart = target.selectionEnd = start + 2;
      });
    }
  }
</script>

<div class="relative w-full h-full min-h-[120px] {className}">
  <textarea
    class="w-full h-full bg-surface-0 text-text-primary text-xs font-mono
      p-3 resize-none border-0 focus:outline-none
      placeholder:text-text-muted/40 leading-relaxed
      {readonly ? 'cursor-default' : ''}"
    {value}
    {readonly}
    {placeholder}
    spellcheck="false"
    autocomplete="off"
    autocorrect="off"
    autocapitalize="off"
    oninput={handleInput}
    onkeydown={handleKeydown}
  ></textarea>

  <!-- Language badge -->
  {#if language}
    <span class="absolute top-2 right-2 px-1.5 py-0.5 text-[9px] uppercase tracking-wider
      text-text-muted bg-surface-2 rounded font-medium">
      {language}
    </span>
  {/if}
</div>
