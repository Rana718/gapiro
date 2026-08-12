<!--
  CodeEditor - CodeMirror 6-based editor with syntax highlighting.
  Supports JSON, XML, HTML, JavaScript, GraphQL, and plain text.
  GPU-optimized, theme matches app dark mode.
-->
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EditorState, type Extension } from '@codemirror/state';
  import { EditorView, keymap, placeholder as cmPlaceholder, lineNumbers, drawSelection, highlightActiveLine, highlightActiveLineGutter } from '@codemirror/view';
  import { defaultKeymap, indentWithTab, history, historyKeymap } from '@codemirror/commands';
  import { syntaxHighlighting, defaultHighlightStyle, bracketMatching, foldGutter, indentOnInput } from '@codemirror/language';
  import { json } from '@codemirror/lang-json';
  import { xml } from '@codemirror/lang-xml';
  import { html } from '@codemirror/lang-html';
  import { javascript } from '@codemirror/lang-javascript';
  import { closeBrackets, closeBracketsKeymap, autocompletion } from '@codemirror/autocomplete';
  import { searchKeymap, highlightSelectionMatches } from '@codemirror/search';
  import { oneDark } from '@codemirror/theme-one-dark';

  interface Props {
    value: string;
    onchange?: (value: string) => void;
    language?: string;
    readonly?: boolean;
    placeholder?: string;
    lineNums?: boolean;
    class?: string;
  }

  let {
    value,
    onchange,
    language = 'text',
    readonly = false,
    placeholder = '',
    lineNums = true,
    class: className = '',
  }: Props = $props();

  let container: HTMLElement;
  let view: EditorView | undefined;
  let skipUpdate = false;

  function getLanguageExtension(lang: string): Extension[] {
    switch (lang) {
      case 'json': return [json()];
      case 'xml': return [xml()];
      case 'html': return [html()];
      case 'javascript': case 'js': return [javascript()];
      case 'graphql': return []; // Basic text mode for now
      default: return [];
    }
  }

  onMount(() => {
    const extensions: Extension[] = [
      oneDark,
      syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
      bracketMatching(),
      indentOnInput(),
      highlightActiveLine(),
      highlightActiveLineGutter(),
      highlightSelectionMatches(),
      drawSelection(),
      history(),
      closeBrackets(),
      autocompletion(),
      foldGutter(),
      keymap.of([
        ...defaultKeymap,
        ...historyKeymap,
        ...closeBracketsKeymap,
        ...searchKeymap,
        indentWithTab,
      ]),
      ...getLanguageExtension(language),
      // Custom dark theme overrides
      EditorView.theme({
        '&': {
          backgroundColor: 'var(--color-surface-inset)',
          color: 'var(--color-text)',
          fontSize: '12px',
          fontFamily: 'var(--font-mono)',
          height: '100%',
        },
        '.cm-content': {
          padding: '8px 0',
          caretColor: 'var(--color-primary)',
        },
        '.cm-gutters': {
          backgroundColor: 'var(--color-surface-inset)',
          color: 'var(--color-text-subtlest)',
          borderRight: '1px solid var(--color-border-subtle)',
          minWidth: '32px',
        },
        '.cm-activeLineGutter': {
          backgroundColor: 'var(--color-surface-highlight)',
        },
        '.cm-activeLine': {
          backgroundColor: 'rgba(99, 102, 241, 0.04)',
        },
        '.cm-cursor': {
          borderLeftColor: 'var(--color-primary)',
        },
        '.cm-selectionBackground': {
          backgroundColor: 'rgba(99, 102, 241, 0.2) !important',
        },
        '&.cm-focused .cm-selectionBackground': {
          backgroundColor: 'rgba(99, 102, 241, 0.25) !important',
        },
        '.cm-tooltip': {
          backgroundColor: 'var(--color-surface-active)',
          border: '1px solid var(--color-border)',
          borderRadius: '6px',
        },
        '.cm-tooltip-autocomplete': {
          '& > ul > li[aria-selected]': {
            backgroundColor: 'var(--color-primary)',
            color: 'white',
          },
        },
      }),
    ];

    if (lineNums) {
      extensions.push(lineNumbers());
    }

    if (placeholder) {
      extensions.push(cmPlaceholder(placeholder));
    }

    if (readonly) {
      extensions.push(EditorState.readOnly.of(true));
      extensions.push(EditorView.editable.of(false));
    } else {
      extensions.push(EditorView.updateListener.of((update) => {
        if (update.docChanged && !skipUpdate) {
          onchange?.(update.state.doc.toString());
        }
      }));
    }

    view = new EditorView({
      state: EditorState.create({
        doc: value,
        extensions,
      }),
      parent: container,
    });
  });

  onDestroy(() => {
    view?.destroy();
  });

  // Update editor when value prop changes externally
  $effect(() => {
    if (view && value !== view.state.doc.toString()) {
      skipUpdate = true;
      view.dispatch({
        changes: { from: 0, to: view.state.doc.length, insert: value },
      });
      skipUpdate = false;
    }
  });
</script>

<div bind:this={container} class="w-full h-full overflow-hidden {className}"></div>

<style>
  /* Ensure CodeMirror fills its container */
  :global(.cm-editor) {
    height: 100% !important;
  }
  :global(.cm-scroller) {
    overflow: auto !important;
  }
</style>
