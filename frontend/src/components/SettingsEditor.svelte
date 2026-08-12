<!--
  SettingsEditor - Request-level settings (timeout, redirects, SSL).
-->
<script lang="ts">
  import { request } from '../stores/app.svelte';
</script>

<div class="p-4 flex flex-col gap-4 overflow-y-auto">
  <!-- Timeout -->
  <label class="flex items-center justify-between">
    <div class="flex flex-col">
      <span class="text-xs font-medium text-text">Request Timeout</span>
      <span class="text-[10px] text-text-subtlest">Time in seconds (0 = no timeout)</span>
    </div>
    <input
      type="number"
      min="0"
      max="600"
      value={request.settings.timeout}
      oninput={(e) => { request.settings.timeout = parseInt((e.target as HTMLInputElement).value) || 0; }}
      class="w-20 px-2.5 py-1.5 text-xs text-center font-mono
        bg-surface-inset border border-border rounded-md
        text-text focus:outline-none focus:border-border-focus"
    />
  </label>

  <!-- Follow Redirects -->
  <label class="flex items-center justify-between cursor-pointer">
    <div class="flex flex-col">
      <span class="text-xs font-medium text-text">Follow Redirects</span>
      <span class="text-[10px] text-text-subtlest">Automatically follow HTTP redirects</span>
    </div>
    <input
      type="checkbox"
      checked={request.settings.followRedirects}
      onchange={() => { request.settings.followRedirects = !request.settings.followRedirects; }}
      class="w-4 h-4 rounded border-border bg-surface-inset
        checked:bg-primary checked:border-primary accent-primary cursor-pointer"
    />
  </label>

  <!-- Max Redirects -->
  {#if request.settings.followRedirects}
    <label class="flex items-center justify-between">
      <div class="flex flex-col">
        <span class="text-xs font-medium text-text">Max Redirects</span>
        <span class="text-[10px] text-text-subtlest">Maximum number of redirects to follow</span>
      </div>
      <input
        type="number"
        min="1"
        max="50"
        value={request.settings.maxRedirects}
        oninput={(e) => { request.settings.maxRedirects = parseInt((e.target as HTMLInputElement).value) || 10; }}
        class="w-20 px-2.5 py-1.5 text-xs text-center font-mono
          bg-surface-inset border border-border rounded-md
          text-text focus:outline-none focus:border-border-focus"
      />
    </label>
  {/if}

  <!-- Verify SSL -->
  <label class="flex items-center justify-between cursor-pointer">
    <div class="flex flex-col">
      <span class="text-xs font-medium text-text">Verify SSL Certificate</span>
      <span class="text-[10px] text-text-subtlest">Reject self-signed or invalid certificates</span>
    </div>
    <input
      type="checkbox"
      checked={request.settings.verifySSL}
      onchange={() => { request.settings.verifySSL = !request.settings.verifySSL; }}
      class="w-4 h-4 rounded border-border bg-surface-inset
        checked:bg-primary checked:border-primary accent-primary cursor-pointer"
    />
  </label>
</div>
