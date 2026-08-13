<!--
  GraphQLPane - GraphQL editor. Command bar (POST implied) + tabbed
  Query / Variables (CodeMirror) / Headers / Docs (live introspection).
  Executes via executeGraphQL() against the Go GraphQLService.
-->
<script lang="ts">
  import { request, response } from '../../stores/app.svelte';
  import { introspectGraphQL } from '../../lib/http';
  import { emptyPair } from '../../lib/utils';
  import type { GraphQLConfig } from '../../lib/types';
  import CommandBar from '../CommandBar.svelte';
  import TabBar from '../core/TabBar.svelte';
  import PairEditor from '../core/PairEditor.svelte';
  import Icon from '../core/Icon.svelte';

  interface Props {
    onSend: () => void;
    onCancel: () => void;
  }

  let { onSend, onCancel }: Props = $props();

  let activeTab = $state<'query' | 'variables' | 'headers' | 'docs'>('query');

  // GraphQL config always exists for a graphql tab, but guard for safety.
  const gql = $derived<GraphQLConfig>(request.graphql ?? { query: '', variables: '{}' });
  const headerCount = $derived(request.headers.filter(h => h.name !== '').length);

  const tabs = $derived([
    { id: 'query', label: 'Query' },
    { id: 'variables', label: 'Variables' },
    { id: 'headers', label: 'Headers', badge: headerCount || undefined },
    { id: 'docs', label: 'Docs' },
  ]);

  function setQuery(v: string) { if (request.graphql) request.graphql.query = v; }
  function setVariables(v: string) { if (request.graphql) request.graphql.variables = v; }

  // ── Introspection (Docs) ──────────────────────────────────────────────
  let schema = $state<any>(null);
  let docsLoading = $state(false);
  let docsError = $state<string | null>(null);
  let docsSearch = $state('');

  async function loadDocs() {
    if (!request.url) { docsError = 'Enter a GraphQL endpoint URL first'; return; }
    docsLoading = true;
    docsError = null;
    try {
      const headers = request.headers.filter(h => h.name !== '').map(h => ({ name: h.name, value: h.value, enabled: h.enabled }));
      const result = await introspectGraphQL(request.url, headers);
      if (!result) { docsError = 'Introspection requires the desktop app (run with `task dev`)'; return; }
      schema = result;
    } catch (err: any) {
      docsError = err?.message ?? 'Introspection failed';
    } finally {
      docsLoading = false;
    }
  }

  const filteredQueries = $derived(
    (schema?.queries ?? []).filter((f: any) => !docsSearch || f.name.toLowerCase().includes(docsSearch.toLowerCase()))
  );
  const filteredMutations = $derived(
    (schema?.mutations ?? []).filter((f: any) => !docsSearch || f.name.toLowerCase().includes(docsSearch.toLowerCase()))
  );
</script>

<div class="flex flex-col h-full bg-surface rounded-lg border border-border-subtle overflow-hidden">
  <CommandBar
    protocol="graphql"
    url={request.url}
    loading={response.loading}
    placeholder="https://api.example.com/graphql"
    onUrlChange={(u) => { request.url = u; }}
    {onSend}
    {onCancel}
  />

  <TabBar {tabs} active={activeTab} onchange={(id) => { activeTab = id as any; }} />

  <div class="flex-1 overflow-hidden">
    {#if activeTab === 'query'}
      {#await import('../core/CodeEditor.svelte') then { default: CodeEditor }}
        <CodeEditor value={gql.query} onchange={setQuery} language="graphql" placeholder={'query {\n  \n}'} />
      {/await}

    {:else if activeTab === 'variables'}
      {#await import('../core/CodeEditor.svelte') then { default: CodeEditor }}
        <CodeEditor value={gql.variables} onchange={setVariables} language="json" placeholder={'{\n  "id": "123"\n}'} />
      {/await}

    {:else if activeTab === 'headers'}
      <PairEditor
        pairs={request.headers}
        namePlaceholder="Header"
        valuePlaceholder="Value"
        onchange={(p) => { request.headers = p; }}
      />

    {:else if activeTab === 'docs'}
      <div class="flex flex-col h-full">
        <div class="flex items-center gap-2 px-3 py-2 border-b border-border-subtle shrink-0">
          <button
            type="button"
            onclick={loadDocs}
            disabled={docsLoading}
            class="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-medium
              bg-protocol-graphql/10 text-protocol-graphql hover:bg-protocol-graphql/20
              disabled:opacity-50 transition-colors"
          >
            <Icon name={docsLoading ? 'refresh' : 'book'} size={13} class={docsLoading ? 'spinner' : ''} />
            {schema ? 'Refresh schema' : 'Load schema'}
          </button>
          {#if schema}
            <div class="flex items-center gap-1.5 flex-1 min-w-0">
              <Icon name="search" size={12} class="text-text-subtlest" />
              <input
                type="text"
                bind:value={docsSearch}
                placeholder="Search fields..."
                class="flex-1 min-w-0 bg-transparent text-xs text-text placeholder:text-placeholder border-0 focus:outline-none"
              />
            </div>
          {/if}
        </div>

        <div class="flex-1 overflow-y-auto p-3">
          {#if docsError}
            <div class="flex items-center gap-2 text-xs text-danger">
              <Icon name="alert-circle" size={14} />
              {docsError}
            </div>
          {:else if docsLoading}
            <div class="flex items-center justify-center h-full text-text-subtlest text-xs gap-2">
              <span class="size-4 border-2 border-protocol-graphql/30 border-t-protocol-graphql rounded-full spinner"></span>
              Introspecting schema…
            </div>
          {:else if !schema}
            <div class="flex flex-col items-center justify-center h-full gap-2 text-center">
              <Icon name="brand-graphql" size={28} class="text-protocol-graphql/40" />
              <p class="text-xs text-text-subtle">Load the schema to explore queries & mutations</p>
            </div>
          {:else}
            {#if filteredQueries.length}
              <div class="mb-4">
                <div class="text-[10px] font-semibold uppercase tracking-wider text-text-subtlest mb-1.5">Queries</div>
                <div class="flex flex-col gap-0.5">
                  {#each filteredQueries as f (f.name)}
                    <div class="px-2 py-1.5 rounded-md hover:bg-surface-highlight/50">
                      <div class="flex items-baseline gap-2">
                        <span class="text-xs font-mono text-protocol-graphql">{f.name}</span>
                        <span class="text-[11px] font-mono text-text-subtlest truncate">{f.type}</span>
                      </div>
                      {#if f.description}<div class="text-[11px] text-text-subtle mt-0.5">{f.description}</div>{/if}
                    </div>
                  {/each}
                </div>
              </div>
            {/if}
            {#if filteredMutations.length}
              <div class="mb-4">
                <div class="text-[10px] font-semibold uppercase tracking-wider text-text-subtlest mb-1.5">Mutations</div>
                <div class="flex flex-col gap-0.5">
                  {#each filteredMutations as f (f.name)}
                    <div class="px-2 py-1.5 rounded-md hover:bg-surface-highlight/50">
                      <div class="flex items-baseline gap-2">
                        <span class="text-xs font-mono text-warning">{f.name}</span>
                        <span class="text-[11px] font-mono text-text-subtlest truncate">{f.type}</span>
                      </div>
                      {#if f.description}<div class="text-[11px] text-text-subtle mt-0.5">{f.description}</div>{/if}
                    </div>
                  {/each}
                </div>
              </div>
            {/if}
            {#if !filteredQueries.length && !filteredMutations.length}
              <div class="text-xs text-text-subtlest text-center py-6">No matching fields</div>
            {/if}
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>
