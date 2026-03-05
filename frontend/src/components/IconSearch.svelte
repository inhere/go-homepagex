<script>
  import { createEventDispatcher, onMount } from 'svelte';

  const dispatch = createEventDispatcher();

  let searchQuery = '';
  let icons = [];
  let filteredIcons = [];
  let loading = true;
  let error = null;

  const ICON_CDN = 'https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons@master';
  const METADATA_URL = `${ICON_CDN}/metadata.json`;

  onMount(async () => {
    try {
      const response = await fetch(METADATA_URL);
      if (!response.ok) throw new Error('Failed to load icons');
      
      const metadata = await response.json();
      icons = Object.entries(metadata).map(([name, data]) => ({
        name,
        url: `${ICON_CDN}/png/${name}.png`,
        aliases: data.aliases || []
      }));
      
      filteredIcons = icons.slice(0, 20);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  });

  function handleSearch(event) {
    searchQuery = event.target.value.toLowerCase();
    
    if (!searchQuery) {
      filteredIcons = icons.slice(0, 20);
      return;
    }

    const results = icons.filter(icon => {
      const nameMatch = icon.name.toLowerCase().includes(searchQuery);
      const aliasMatch = icon.aliases.some(alias => 
        alias.toLowerCase().includes(searchQuery)
      );
      return nameMatch || aliasMatch;
    });

    filteredIcons = results.slice(0, 50);
  }

  function selectIcon(icon) {
    dispatch('select', icon.url);
  }
</script>

<div class="icon-search">
  <div class="search-header">
    <div class="search-input-wrapper">
      <i class="fas fa-search"></i>
      <input
        type="text"
        class="search-input"
        placeholder="搜索图标..."
        value={searchQuery}
        on:input={handleSearch}
      />
    </div>
  </div>

  <div class="icon-grid-container">
    {#if loading}
      <div class="loading-state">
        <i class="fas fa-spinner fa-spin"></i>
        <span>加载图标中...</span>
      </div>
    {:else if error}
      <div class="error-state">
        <i class="fas fa-exclamation-triangle"></i>
        <span>{error}</span>
      </div>
    {:else if filteredIcons.length === 0}
      <div class="empty-state">
        <i class="fas fa-search"></i>
        <span>未找到匹配的图标</span>
      </div>
    {:else}
      <div class="icon-grid">
        {#each filteredIcons as icon}
          <button
            class="icon-card"
            on:click={() => selectIcon(icon)}
            title={icon.name}
          >
            <img src={icon.url} alt={icon.name} class="icon-image" />
            <span class="icon-name">{icon.name}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .icon-search {
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: 100%;
  }

  .search-header {
    flex-shrink: 0;
  }

  .search-input-wrapper {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 8px;
    transition: all 0.3s ease;
  }

  .search-input-wrapper:focus-within {
    background: rgba(255, 255, 255, 0.1);
    border-color: var(--theme-primary, #4a9eff);
  }

  .search-input-wrapper i {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #e4e4e4;
    font-size: 0.9rem;
  }

  .search-input::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .icon-grid-container {
    flex: 1;
    overflow-x: auto;
    overflow-y: hidden;
    min-height: 0;
  }

  .icon-grid {
    display: flex;
    gap: 12px;
    padding: 4px;
    white-space: nowrap;
  }

  .icon-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    padding: 10px 12px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    flex-shrink: 0;
    min-width: 80px;
  }

  .icon-card:hover {
    background: rgba(255, 255, 255, 0.1);
    border-color: var(--theme-primary, #4a9eff);
    transform: translateY(-2px);
  }

  .icon-image {
    width: 48px;
    height: 48px;
    object-fit: contain;
    filter: brightness(0) invert(1);
    opacity: 0.9;
  }

  .icon-name {
    font-size: 0.7rem;
    color: rgba(255, 255, 255, 0.7);
    text-align: center;
    word-break: break-word;
    line-height: 1.2;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 40px 20px;
    color: rgba(255, 255, 255, 0.5);
  }

  .loading-state i,
  .error-state i,
  .empty-state i {
    font-size: 2rem;
  }

  .error-state {
    color: #ff6b6b;
  }

  .icon-grid-container::-webkit-scrollbar {
    width: 6px;
  }

  .icon-grid-container::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 3px;
  }

  .icon-grid-container::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }

  .icon-grid-container::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
  }
</style>
