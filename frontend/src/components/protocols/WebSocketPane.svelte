<!--
  WebSocketPane - Live WebSocket console. A connection bar (ws:// address +
  Connect/Disconnect) sits above a color-coded, timestamped message log and a
  composer. Frames stream in via the Go WebSocketService over the "ws:event"
  channel. Keyed by request.id in App, so each tab owns its connection.
-->
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { request } from '../../stores/app.svelte';
  import { wsConnect, wsSend, wsClose, wsOnEvent } from '../../lib/http';
  import type { WsMessage } from '../../lib/types';
  import Icon from '../core/Icon.svelte';

  const connID = request.id;
  const accent = 'var(--color-protocol-websocket)';

  let status = $state<'idle' | 'connecting' | 'open' | 'closed' | 'error'>('idle');
  let messages = $state.raw<WsMessage[]>([]);
  let draft = $state('');
  let unsub: (() => void) | null = null;
  let logEl = $state<HTMLDivElement | null>(null);

  const connected = $derived(status === 'open');
  const sentCount = $derived(messages.filter(m => m.direction === 'sent').length);
  const recvCount = $derived(messages.filter(m => m.direction === 'received').length);

  function push(m: WsMessage) {
    messages = [...messages, m];
  }
  function system(data: string) {
    push({ direction: 'system', data, timestamp: Date.now() });
  }

  // Auto-scroll to the newest frame.
  $effect(() => {
    messages.length;
    if (logEl) logEl.scrollTop = logEl.scrollHeight;
  });

  onMount(async () => {
    unsub = await wsOnEvent((ev) => {
      if (!ev || ev.connID !== connID) return;
      switch (ev.type) {
        case 'open': status = 'open'; system('Connected'); break;
        case 'message': push({ direction: 'received', data: ev.data ?? '', timestamp: Date.now() }); break;
        case 'sent': push({ direction: 'sent', data: ev.data ?? '', timestamp: Date.now() }); break;
        case 'close':
          status = 'closed';
          system(ev.code ? `Disconnected (${ev.code})` : 'Disconnected');
          break;
        case 'error':
          status = 'error';
          system(`Error: ${ev.data ?? 'connection failed'}`);
          break;
      }
    });
  });

  onDestroy(() => {
    unsub?.();
    if (status === 'open' || status === 'connecting') void wsClose(connID);
  });

  async function connect() {
    if (!request.url) return;
    status = 'connecting';
    const headers = request.headers.filter(h => h.name !== '').map(h => ({ id: h.id, name: h.name, value: h.value, enabled: h.enabled }));
    try {
      await wsConnect(connID, request.url, headers, request.websocket?.protocols ?? '');
    } catch (err: any) {
      status = 'error';
      system(err?.message ?? 'Failed to connect');
    }
  }

  async function disconnect() {
    await wsClose(connID);
    status = 'closed';
  }

  function toggle(e: Event) {
    e.preventDefault();
    if (connected || status === 'connecting') disconnect();
    else connect();
  }

  async function send() {
    const msg = draft;
    if (!connected || !msg.trim()) return;
    draft = '';
    try {
      await wsSend(connID, msg);
    } catch (err: any) {
      system(err?.message ?? 'Send failed');
    }
  }

  function onComposerKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send(); }
  }

  function fmtTime(ts: number): string {
    const d = new Date(ts);
    const p = (n: number, l = 2) => String(n).padStart(l, '0');
    return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}.${p(d.getMilliseconds(), 3)}`;
  }

  const statusMeta = $derived(
    status === 'open' ? { color: 'var(--color-success)', label: 'Connected' }
    : status === 'connecting' ? { color: accent, label: 'Connecting…' }
    : status === 'error' ? { color: 'var(--color-danger)', label: 'Error' }
    : status === 'closed' ? { color: 'var(--color-text-subtlest)', label: 'Disconnected' }
    : { color: 'var(--color-text-subtlest)', label: 'Not connected' }
  );
</script>

<div class="flex flex-col h-full bg-surface rounded-lg border border-border-subtle overflow-hidden">
  <!-- Connection bar -->
  <form onsubmit={toggle} class="flex items-stretch gap-2 px-3 py-2.5 shrink-0">
    <div
      class="flex items-stretch flex-1 min-w-0 h-9 rounded-lg border border-border
        bg-surface-inset focus-within:border-border-focus transition-colors overflow-hidden"
    >
      <div class="flex items-center gap-1.5 px-3 border-r border-border select-none" style="color: {accent}" title="WebSocket">
        <span class="size-2 rounded-full shrink-0" style="background: {statusMeta.color}"></span>
        <span class="text-[11px] font-bold font-mono tracking-tight">WS</span>
      </div>
      <input
        type="text"
        value={request.url}
        oninput={(e) => { request.url = (e.target as HTMLInputElement).value; }}
        placeholder="wss://echo.websocket.org"
        disabled={connected || status === 'connecting'}
        class="flex-1 min-w-0 bg-transparent px-3 text-xs text-text font-mono
          placeholder:text-placeholder border-0 focus:outline-none disabled:opacity-70"
        spellcheck="false"
        autocomplete="off"
      />
    </div>

    {#if connected || status === 'connecting'}
      <button
        type="submit"
        class="flex items-center gap-1.5 h-9 px-4 rounded-lg text-xs font-semibold shrink-0
          bg-danger/10 text-danger hover:bg-danger/20 transition-colors"
      >
        {#if status === 'connecting'}
          <span class="size-3.5 border-2 border-danger/40 border-t-danger rounded-full spinner"></span>
        {:else}
          <Icon name="wifi-off" size={14} />
        {/if}
        Disconnect
      </button>
    {:else}
      <button
        type="submit"
        disabled={!request.url}
        class="flex items-center gap-1.5 h-9 px-4 rounded-lg text-xs font-semibold shrink-0
          text-white transition-opacity disabled:opacity-40 disabled:cursor-not-allowed hover:opacity-90"
        style="background: {accent}"
      >
        <Icon name="plug-connected" size={14} />
        Connect
      </button>
    {/if}
  </form>

  <!-- Subprotocols row -->
  <div class="flex items-center gap-2 px-3 pb-2.5 border-b border-border-subtle shrink-0">
    <Icon name="bolt" size={13} class="text-text-subtlest shrink-0" />
    <input
      type="text"
      value={request.websocket?.protocols ?? ''}
      oninput={(e) => { if (request.websocket) request.websocket.protocols = (e.target as HTMLInputElement).value; }}
      placeholder="subprotocols (comma-separated, optional)"
      disabled={connected || status === 'connecting'}
      class="flex-1 min-w-0 bg-transparent text-xs text-text font-mono placeholder:text-placeholder border-0 focus:outline-none disabled:opacity-70"
      spellcheck="false"
    />
  </div>

  <!-- Log toolbar -->
  <div class="flex items-center gap-3 px-3 h-8 border-b border-border-subtle shrink-0 text-[11px]">
    <span class="flex items-center gap-1.5" style="color: {statusMeta.color}">
      <span class="size-1.5 rounded-full" style="background: {statusMeta.color}"></span>
      {statusMeta.label}
    </span>
    <div class="flex-1"></div>
    <span class="flex items-center gap-1 text-text-subtle" title="Sent">
      <Icon name="arrow-up" size={12} class="text-protocol-websocket" />{sentCount}
    </span>
    <span class="flex items-center gap-1 text-text-subtle" title="Received">
      <Icon name="arrow-down" size={12} class="text-success" />{recvCount}
    </span>
    <button
      type="button"
      onclick={() => { messages = []; }}
      disabled={!messages.length}
      class="flex items-center gap-1 text-text-subtlest hover:text-text disabled:opacity-40 transition-colors"
      title="Clear log"
    >
      <Icon name="trash" size={12} /> Clear
    </button>
  </div>

  <!-- Message log -->
  <div bind:this={logEl} class="flex-1 overflow-y-auto p-2 font-mono text-xs">
    {#if !messages.length}
      <div class="flex flex-col items-center justify-center h-full gap-2 text-center select-none">
        <Icon name="plug-connected" size={28} class="text-protocol-websocket/40" />
        <p class="text-xs text-text-subtle font-sans">
          {connected ? 'Connected — send a message below' : 'Connect to start streaming frames'}
        </p>
      </div>
    {:else}
      <div class="flex flex-col gap-px">
        {#each messages as m, i (i)}
          {#if m.direction === 'system'}
            <div class="flex items-center gap-2 px-2 py-1 text-text-subtlest">
              <Icon name="info-circle" size={12} class="shrink-0" />
              <span class="font-sans">{m.data}</span>
              <span class="ml-auto text-[10px] tabular-nums opacity-60">{fmtTime(m.timestamp)}</span>
            </div>
          {:else}
            {@const sent = m.direction === 'sent'}
            <div class="flex items-start gap-2 px-2 py-1 rounded hover:bg-surface-highlight/40 group">
              <Icon
                name={sent ? 'arrow-up' : 'arrow-down'}
                size={12}
                class="mt-0.5 shrink-0 {sent ? 'text-protocol-websocket' : 'text-success'}"
              />
              <pre class="flex-1 min-w-0 whitespace-pre-wrap break-all text-text leading-relaxed">{m.data}</pre>
              <span class="text-[10px] tabular-nums text-text-subtlest opacity-0 group-hover:opacity-100 transition-opacity shrink-0 mt-0.5">{fmtTime(m.timestamp)}</span>
            </div>
          {/if}
        {/each}
      </div>
    {/if}
  </div>

  <!-- Composer -->
  <div class="flex items-end gap-2 p-2 border-t border-border-subtle shrink-0">
    <textarea
      value={draft}
      oninput={(e) => { draft = (e.target as HTMLTextAreaElement).value; }}
      onkeydown={onComposerKeydown}
      placeholder={connected ? 'Message… (Enter to send, Shift+Enter for newline)' : 'Connect to send messages'}
      disabled={!connected}
      rows="1"
      class="flex-1 min-w-0 max-h-28 resize-none bg-surface-inset text-xs text-text font-mono
        px-3 py-2 rounded-lg border border-border placeholder:text-placeholder
        focus:outline-none focus:border-border-focus disabled:opacity-60"
      spellcheck="false"
    ></textarea>
    <button
      type="button"
      onclick={send}
      disabled={!connected || !draft.trim()}
      class="flex items-center gap-1.5 h-9 px-4 rounded-lg text-xs font-semibold shrink-0
        text-white transition-opacity disabled:opacity-40 disabled:cursor-not-allowed hover:opacity-90"
      style="background: {accent}"
      title="Send (Enter)"
    >
      <Icon name="send" size={14} />
      Send
    </button>
  </div>
</div>
