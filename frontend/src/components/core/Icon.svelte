<!--
  Icon - Lightweight wrapper for a curated set of Tabler-style line icons.
  Single source of truth for all iconography in hand-written components.
  Svelte 5 compatible (runes). Usage: <Icon name="send" size={16} />
-->
<script module lang="ts">
  // Icon paths use a 24×24 viewBox, no fill, currentColor stroke.
  const icons: Record<string, string> = {
    // ── Actions ──────────────────────────────────────────────
    'plus': '<path d="M12 5v14"/><path d="M5 12h14"/>',
    'x': '<path d="M18 6L6 18"/><path d="M6 6l12 12"/>',
    'check': '<path d="M5 12l5 5L20 7"/>',
    'search': '<circle cx="10" cy="10" r="7"/><path d="M21 21l-6-6"/>',
    'copy': '<rect x="8" y="8" width="12" height="12" rx="2"/><path d="M16 8V6a2 2 0 00-2-2H6a2 2 0 00-2 2v8a2 2 0 002 2h2"/>',
    'copy-check': '<path d="M7 9.667A2.667 2.667 0 019.667 7h8.666A2.667 2.667 0 0121 9.667v8.666A2.667 2.667 0 0118.333 21H9.667A2.667 2.667 0 017 18.333z"/><path d="M4.012 16.737A2 2 0 013 15V5a2 2 0 012-2h10c.75 0 1.158.385 1.5 1"/><path d="M11 14l2 2 4-4"/>',
    'eye': '<path d="M10 12a2 2 0 104 0 2 2 0 10-4 0"/><path d="M21 12c-2.4 4-5.4 6-9 6-3.6 0-6.6-2-9-6 2.4-4 5.4-6 9-6 3.6 0 6.6 2 9 6"/>',
    'trash': '<path d="M4 7h16"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M5 7l1 12a2 2 0 002 2h8a2 2 0 002-2l1-12"/><path d="M9 7V4a1 1 0 011-1h4a1 1 0 011 1v3"/>',
    'save': '<path d="M6 4h10l4 4v10a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2"/><path d="M12 14a2 2 0 100-4 2 2 0 000 4"/><path d="M14 4v4H8V4"/>',
    'settings': '<path d="M10.325 4.317a1.724 1.724 0 013.35 0 1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572 1.724 1.724 0 010 3.35 1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065 1.724 1.724 0 01-3.35 0 1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572 1.724 1.724 0 010-3.35 1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37 1 .608 2.296.07 2.572-1.065z"/><circle cx="12" cy="12" r="3"/>',
    'filter': '<path d="M5.5 5h13a1 1 0 01.5 1.5L14 12v7l-4-3v-4L4.5 6.5A1 1 0 015.5 5"/>',
    'copy-plus': '<path d="M12 5v14"/><path d="M5 12h14"/>',
    'refresh': '<path d="M20 11A8 8 0 006.3 5.7L3 9"/><path d="M4 13a8 8 0 0013.7 3.3L21 15"/><path d="M3 4v5h5"/><path d="M21 20v-5h-5"/>',
    'download': '<path d="M4 17v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M7 11l5 5 5-5"/><path d="M12 4v12"/>',
    'link': '<path d="M9 15l6-6"/><path d="M11 6l.5-.5a4 4 0 015.5 5.5l-.5.5"/><path d="M13 18l-.5.5a4 4 0 01-5.5-5.5l.5-.5"/>',
    'file': '<path d="M14 3v4a1 1 0 001 1h4"/><path d="M17 21H7a2 2 0 01-2-2V5a2 2 0 012-2h7l5 5v11a2 2 0 01-2 2z"/>',
    'file-code': '<path d="M14 3v4a1 1 0 001 1h4"/><path d="M17 21H7a2 2 0 01-2-2V5a2 2 0 012-2h7l5 5v11a2 2 0 01-2 2z"/><path d="M10 13l-2 2 2 2"/><path d="M14 13l2 2-2 2"/>',
    'book': '<path d="M3 5a2 2 0 012-2h6v16H5a2 2 0 00-2 2z"/><path d="M21 5a2 2 0 00-2-2h-6v16h6a2 2 0 012 2z"/>',
    'edit': '<path d="M4 20h4l10.5-10.5a2.121 2.121 0 00-3-3L5 17v3z"/><path d="M13.5 6.5l3 3"/>',
    'bolt': '<path d="M13 3L4 14h7l-1 7 9-11h-7l1-7z"/>',

    // ── Navigation / layout ──────────────────────────────────
    'chevron-down': '<path d="M6 9l6 6 6-6"/>',
    'chevron-right': '<path d="M9 6l6 6-6 6"/>',
    'chevron-left': '<path d="M15 6l-6 6 6 6"/>',
    'chevron-up': '<path d="M6 15l6-6 6 6"/>',
    'arrow-right': '<path d="M5 12h14"/><path d="M13 5l7 7-7 7"/>',
    'dots-vertical': '<circle cx="12" cy="12" r="1"/><circle cx="12" cy="5" r="1"/><circle cx="12" cy="19" r="1"/>',
    'layout-sidebar-left-collapse': '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18"/><path d="M14 9l-2 3 2 3"/>',
    'layout-sidebar-left-expand': '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18"/><path d="M12 9l2 3-2 3"/>',
    'layout-sidebar-right': '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M15 3v18"/>',
    'layout-bottombar': '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 15h18"/>',
    'maximize': '<path d="M4 8V6a2 2 0 012-2h2"/><path d="M4 16v2a2 2 0 002 2h2"/><path d="M16 4h2a2 2 0 012 2v2"/><path d="M16 20h2a2 2 0 002-2v-2"/>',
    'minimize': '<path d="M8 4v2a2 2 0 01-2 2H4"/><path d="M4 16h2a2 2 0 012 2v2"/><path d="M16 4v2a2 2 0 002 2h2"/><path d="M20 16h-2a2 2 0 00-2 2v2"/>',

    // ── Folders / collection ─────────────────────────────────
    'folder': '<path d="M5 4h4l3 3h7a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2"/>',
    'folder-plus': '<path d="M12 19H5a2 2 0 01-2-2V6a2 2 0 012-2h4l3 3h7a2 2 0 012 2v3.5"/><path d="M16 19h6"/><path d="M19 16v6"/>',

    // ── Protocols ────────────────────────────────────────────
    'world': '<circle cx="12" cy="12" r="9"/><path d="M3.6 9h16.8"/><path d="M3.6 15h16.8"/><path d="M12 3a15 15 0 014 9 15 15 0 01-4 9 15 15 0 01-4-9 15 15 0 014-9"/>',
    'brand-graphql': '<path d="M5.308 7.265l5.385-3.118M13.308 4.147l5.384 3.118M19 9v6M18.692 16.735l-5.384 3.118M10.692 19.853l-5.384-3.118M5 15V9"/><circle cx="12" cy="3.5" r="1.5"/><circle cx="19.5" cy="8" r="1.5"/><circle cx="19.5" cy="16" r="1.5"/><circle cx="12" cy="20.5" r="1.5"/><circle cx="4.5" cy="16" r="1.5"/><circle cx="4.5" cy="8" r="1.5"/>',
    'server': '<rect x="3" y="4" width="18" height="8" rx="2"/><rect x="3" y="12" width="18" height="8" rx="2"/><path d="M7 8h.01M7 16h.01"/>',
    'plug-connected': '<path d="M7 12l5 5-1.5 1.5a3.5 3.5 0 11-5-5L7 12"/><path d="M17 12l-5-5 1.5-1.5a3.5 3.5 0 115 5L17 12"/><path d="M14.5 17.5l2 2"/><path d="M7.5 6.5l2 2"/>',

    // ── Status / misc ────────────────────────────────────────
    'code': '<path d="M7 8l-4 4 4 4"/><path d="M17 8l4 4-4 4"/>',
    'clock': '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 3"/>',
    'send': '<path d="M10 14L21 3"/><path d="M21 3l-6.5 18a.55.55 0 01-1 0L10 14l-7-3.5a.55.55 0 010-1L21 3"/>',
    'play': '<path d="M7 4v16l13-8z"/>',
    'circle-dot': '<circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3" fill="currentColor" stroke="none"/>',
    'alert-circle': '<circle cx="12" cy="12" r="9"/><path d="M12 8v4"/><path d="M12 16h.01"/>',
    'info-circle': '<circle cx="12" cy="12" r="9"/><path d="M12 8h.01"/><path d="M11 12h1v4h1"/>',
    'arrow-up': '<path d="M12 5v14"/><path d="M5 12l7-7 7 7"/>',
    'arrow-down': '<path d="M12 5v14"/><path d="M5 12l7 7 7-7"/>',
    'wifi-off': '<path d="M3 3l18 18"/><path d="M9 17l3-3 3 3"/><path d="M6.7 10.7a8 8 0 013.3-1.6M14 9.1a8 8 0 013.3 1.6"/>',
  };

  export function getIconPath(name: string): string {
    return icons[name] ?? icons['x'] ?? '';
  }
</script>

<script lang="ts">
  interface Props {
    name: string;
    size?: number;
    class?: string;
    strokeWidth?: number;
  }

  let { name, size = 16, class: className = '', strokeWidth = 1.75 }: Props = $props();
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width={strokeWidth}
  stroke-linecap="round"
  stroke-linejoin="round"
  class="tabler-icon shrink-0 {className}"
  aria-hidden="true"
>
  {@html getIconPath(name)}
</svg>
