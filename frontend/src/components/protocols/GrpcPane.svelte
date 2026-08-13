<!--
  GrpcPane - gRPC unary editor. Command bar (server address) + proto loader +
  method picker + request message (seeded from the method's input template) +
  metadata. Invokes via executeGrpc() against the Go GrpcService.
-->
<script lang="ts">
  import { request, response } from '../../stores/app.svelte';
  import { parseProtoFile } from '../../lib/http';
  import type { GrpcConfig } from '../../lib/types';
  import CommandBar from '../CommandBar.svelte';
  import TabBar from '../core/TabBar.svelte';
  import PairEditor from '../core/PairEditor.svelte';
  import Icon from '../core/Icon.svelte';

  interface Props {
    onSend: () => void;
    onCancel: () => void;
  }

  let { onSend, onCancel }: Props = $props();

  const grpc = $derived<GrpcConfig>(request.grpc ?? { protoFile: '', importPaths: [], fullMethod: '', message: '{}', metadata: [] });

  let activeTab = $state<'message' | 'metadata'>('message');
  let services = $state<any[]>([]);
  let loading = $state(false);
  let loadError = $state<string | null>(null);
  let methodOpen = $state(false);

  const metaCount = $derived((request.grpc?.metadata ?? []).filter(m => m.name !== '').length);
  const allMethods = $derived(services.flatMap((s: any) => (s.methods ?? []).map((m: any) => ({ ...m, service: s.name }))));
  const selectedMethod = $derived(allMethods.find((m: any) => m.fullName === grpc.fullMethod) ?? null);
  const canInvoke = $derived(!!request.url && !!grpc.fullMethod);

  const tabs = $derived([
    { id: 'message', label: 'Message' },
    { id: 'metadata', label: 'Metadata', badge: metaCount || undefined },
  ]);

  async function loadMethods() {
    if (!request.grpc?.protoFile) { loadError = 'Enter a .proto file path first'; return; }
    loading = true;
    loadError = null;
    try {
      const result = await parseProtoFile(request.grpc.protoFile, request.grpc.importPaths ?? []);
      services = result ?? [];
      if (!services.length) loadError = 'No services found in proto file';
    } catch (err: any) {
      loadError = err?.message ?? 'Failed to parse proto file';
      services = [];
    } finally {
      loading = false;
    }
  }

  function selectMethod(m: any) {
    if (!request.grpc) return;
    request.grpc.fullMethod = m.fullName;
    // Seed the message with the input template when empty or still a bare object.
    const cur = request.grpc.message?.trim();
    if (!cur || cur === '{}' || cur === '{\n  \n}') {
      request.grpc.message = m.inputTemplate || '{}';
    }
    methodOpen = false;
  }

  function setImportPaths(v: string) {
    if (request.grpc) request.grpc.importPaths = v.split(',').map(s => s.trim()).filter(Boolean);
  }
  function setMessage(v: string) { if (request.grpc) request.grpc.message = v; }
</script>

<svelte:document onclick={() => { if (methodOpen) methodOpen = false; }} />

<div class="flex flex-col h-full bg-surface rounded-lg border border-border-subtle overflow-hidden">
  <CommandBar
    protocol="grpc"
    url={request.url}
    loading={response.loading}
    placeholder="localhost:50051"
    disabled={!canInvoke}
    onUrlChange={(u) => { request.url = u; }}
    {onSend}
    {onCancel}
  />

  <!-- Proto + method config -->
  <div class="flex flex-col gap-2 px-3 pb-2.5 border-b border-border-subtle shrink-0">
    <div class="flex items-stretch gap-2">
      <div class="flex items-center flex-1 min-w-0 h-8 rounded-md border border-border bg-surface-inset px-2.5 gap-2 focus-within:border-border-focus">
        <Icon name="file-code" size={13} class="text-text-subtlest" />
        <input
          type="text"
          value={grpc.protoFile}
          oninput={(e) => { if (request.grpc) request.grpc.protoFile = (e.target as HTMLInputElement).value; }}
          placeholder="/path/to/service.proto"
          class="flex-1 min-w-0 bg-transparent text-xs text-text font-mono placeholder:text-placeholder border-0 focus:outline-none"
          spellcheck="false"
        />
      </div>
      <button
        type="button"
        onclick={loadMethods}
        disabled={loading}
        class="flex items-center gap-1.5 h-8 px-3 rounded-md text-[11px] font-medium shrink-0
          bg-protocol-grpc/10 text-protocol-grpc hover:bg-protocol-grpc/20 disabled:opacity-50 transition-colors"
      >
        <Icon name={loading ? 'refresh' : 'download'} size={13} class={loading ? 'spinner' : ''} />
        Load methods
      </button>
    </div>

    <!-- Method selector + import paths -->
    <div class="flex items-stretch gap-2">
      <div class="relative flex-1 min-w-0">
        <button
          type="button"
          onclick={(e) => { e.stopPropagation(); if (allMethods.length) methodOpen = !methodOpen; }}
          disabled={!allMethods.length}
          class="w-full flex items-center gap-2 h-8 px-2.5 rounded-md border border-border bg-surface-inset
            text-xs hover:border-border-focus disabled:opacity-50 transition-colors"
        >
          <Icon name="server" size={13} class="text-protocol-grpc shrink-0" />
          {#if selectedMethod}
            <span class="font-mono text-text truncate">{selectedMethod.service}<span class="text-text-subtlest">/</span>{selectedMethod.name}</span>
          {:else}
            <span class="text-text-subtlest">{allMethods.length ? 'Select a method…' : 'Load a proto file to list methods'}</span>
          {/if}
          <div class="flex-1"></div>
          <Icon name="chevron-down" size={13} class="text-text-subtlest shrink-0" />
        </button>

        {#if methodOpen}
          <div
            class="absolute top-full left-0 right-0 mt-1 z-50 max-h-64 overflow-y-auto py-1
              bg-surface-active border border-border rounded-lg shadow-xl"
            onclick={(e) => e.stopPropagation()}
            role="listbox"
            tabindex="-1"
          >
            {#each services as svc (svc.name)}
              <div class="px-3 py-1 text-[10px] font-semibold uppercase tracking-wider text-text-subtlest">{svc.name}</div>
              {#each svc.methods ?? [] as m (m.fullName)}
                <button
                  type="button"
                  onclick={() => selectMethod({ ...m, service: svc.name })}
                  class="w-full flex items-center gap-2 px-3 py-1.5 text-left text-xs hover:bg-surface-highlight
                    {m.fullName === grpc.fullMethod ? 'bg-surface-highlight' : ''}"
                  role="option"
                  aria-selected={m.fullName === grpc.fullMethod}
                >
                  <span class="font-mono text-text truncate flex-1">{m.name}</span>
                  {#if m.isClientStream || m.isServerStream}
                    <span class="text-[9px] font-bold text-protocol-grpc uppercase">stream</span>
                  {/if}
                </button>
              {/each}
            {/each}
          </div>
        {/if}
      </div>

      <div class="flex items-center w-[38%] min-w-0 h-8 rounded-md border border-border bg-surface-inset px-2.5 gap-2 focus-within:border-border-focus">
        <Icon name="folder" size={13} class="text-text-subtlest" />
        <input
          type="text"
          value={(grpc.importPaths ?? []).join(', ')}
          oninput={(e) => setImportPaths((e.target as HTMLInputElement).value)}
          placeholder="import paths (optional)"
          class="flex-1 min-w-0 bg-transparent text-xs text-text font-mono placeholder:text-placeholder border-0 focus:outline-none"
          spellcheck="false"
        />
      </div>
    </div>

    {#if loadError}
      <div class="flex items-center gap-1.5 text-[11px] text-danger">
        <Icon name="alert-circle" size={12} />
        {loadError}
      </div>
    {/if}
  </div>

  <TabBar {tabs} active={activeTab} onchange={(id) => { activeTab = id as any; }} />

  <div class="flex-1 overflow-hidden">
    {#if activeTab === 'message'}
      {#await import('../core/CodeEditor.svelte') then { default: CodeEditor }}
        <CodeEditor value={grpc.message} onchange={setMessage} language="json" placeholder={'{\n  \n}'} />
      {/await}
    {:else if activeTab === 'metadata'}
      <PairEditor
        pairs={request.grpc?.metadata ?? []}
        namePlaceholder="Metadata key"
        valuePlaceholder="Value"
        onchange={(p) => { if (request.grpc) request.grpc.metadata = p; }}
      />
    {/if}
  </div>
</div>
