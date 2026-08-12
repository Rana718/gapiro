<!--
  AuthEditor - Authentication configuration panel.
  Supports: None, Basic Auth, Bearer Token, API Key.
-->
<script lang="ts">
  import { request } from '../stores/app.svelte';
  import type { AuthType } from '../lib/types';

  const authTypes: { id: AuthType; label: string }[] = [
    { id: 'none', label: 'No Auth' },
    { id: 'basic', label: 'Basic Auth' },
    { id: 'bearer', label: 'Bearer Token' },
    { id: 'api-key', label: 'API Key' },
  ];

  function setAuthType(t: AuthType) {
    request.auth.type = t;
    if (t === 'basic' && !request.auth.basic) {
      request.auth.basic = { username: '', password: '' };
    }
    if (t === 'bearer' && !request.auth.bearer) {
      request.auth.bearer = { token: '', prefix: 'Bearer' };
    }
    if (t === 'api-key' && !request.auth.apiKey) {
      request.auth.apiKey = { key: '', value: '', addTo: 'header' };
    }
  }
</script>

<div class="flex flex-col h-full overflow-y-auto">
  <!-- Auth type selector -->
  <div class="flex items-center gap-1 px-3 py-2 border-b border-border-subtle">
    {#each authTypes as at (at.id)}
      <button
        onclick={() => setAuthType(at.id)}
        class="px-2.5 py-1 text-[11px] font-medium rounded
          transition-colors duration-50
          {request.auth.type === at.id
            ? 'bg-primary/15 text-primary border border-primary/30'
            : 'text-text-subtlest hover:text-text-subtle hover:bg-surface-highlight border border-transparent'}"
      >
        {at.label}
      </button>
    {/each}
  </div>

  <!-- Auth config -->
  <div class="p-4 flex flex-col gap-3">
    {#if request.auth.type === 'none'}
      <p class="text-sm text-text-subtlest text-center py-8">No authentication</p>
    {:else if request.auth.type === 'basic'}
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Username</span>
        <input
          type="text"
          value={request.auth.basic?.username ?? ''}
          oninput={(e) => { if (request.auth.basic) request.auth.basic.username = (e.target as HTMLInputElement).value; }}
          placeholder="username"
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </label>
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Password</span>
        <input
          type="password"
          value={request.auth.basic?.password ?? ''}
          oninput={(e) => { if (request.auth.basic) request.auth.basic.password = (e.target as HTMLInputElement).value; }}
          placeholder="password"
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </label>
    {:else if request.auth.type === 'bearer'}
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Prefix</span>
        <input
          type="text"
          value={request.auth.bearer?.prefix ?? 'Bearer'}
          oninput={(e) => { if (request.auth.bearer) request.auth.bearer.prefix = (e.target as HTMLInputElement).value; }}
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </label>
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Token</span>
        <textarea
          value={request.auth.bearer?.token ?? ''}
          oninput={(e) => { if (request.auth.bearer) request.auth.bearer.token = (e.target as HTMLTextAreaElement).value; }}
          placeholder="your-token-here"
          rows="3"
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus resize-none"
        ></textarea>
      </label>
    {:else if request.auth.type === 'api-key'}
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Key</span>
        <input
          type="text"
          value={request.auth.apiKey?.key ?? ''}
          oninput={(e) => { if (request.auth.apiKey) request.auth.apiKey.key = (e.target as HTMLInputElement).value; }}
          placeholder="X-API-Key"
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </label>
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Value</span>
        <input
          type="password"
          value={request.auth.apiKey?.value ?? ''}
          oninput={(e) => { if (request.auth.apiKey) request.auth.apiKey.value = (e.target as HTMLInputElement).value; }}
          placeholder="your-api-key"
          class="px-3 py-1.5 text-xs font-mono bg-surface-inset border border-border rounded-md
            text-text placeholder:text-placeholder focus:outline-none focus:border-border-focus"
        />
      </label>
      <label class="flex flex-col gap-1">
        <span class="text-[11px] font-medium text-text-subtle">Add to</span>
        <select
          value={request.auth.apiKey?.addTo ?? 'header'}
          onchange={(e) => { if (request.auth.apiKey) request.auth.apiKey.addTo = (e.target as HTMLSelectElement).value as 'header' | 'query'; }}
          class="px-3 py-1.5 text-xs bg-surface-inset border border-border rounded-md
            text-text focus:outline-none focus:border-border-focus"
        >
          <option value="header">Header</option>
          <option value="query">Query Parameter</option>
        </select>
      </label>
    {/if}
  </div>
</div>
