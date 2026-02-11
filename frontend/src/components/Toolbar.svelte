<script>
  export let services = [];
  export let viewStyle = 'cards';
  export let onSearch = () => {};

  let searchQuery = '';
  let showResults = false;
  let filteredItems = [];

  $: {
    if (searchQuery.trim()) {
      filteredItems = services.flatMap(service => 
        service.items.map(item => ({
          ...item,
          serviceName: service.name,
          serviceIcon: service.icon
        }))
      ).filter(item => 
        item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.subtitle?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.tag?.toLowerCase().includes(searchQuery.toLowerCase())
      );
      showResults = filteredItems.length > 0;
    } else {
      filteredItems = [];
      showResults = false;
    }
  }

  function handleSearch(event) {
    searchQuery = event.target.value;
    onSearch(searchQuery);
  }

  function clearSearch() {
    searchQuery = '';
    showResults = false;
    onSearch('');
  }

  function toggleStyle() {
    viewStyle = viewStyle === 'cards' ? 'list' : 'cards';
  }

  function navigateTo(url, target) {
    if (target === '_blank') {
      window.open(url, '_blank');
    } else {
      window.location.href = url;
    }
    showResults = false;
    searchQuery = '';
    onSearch('');
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      showResults = false;
      document.querySelector('.search-input').blur();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="toolbar">
  <div class="search-container">
    <div class="search-box">
      <i class="fas fa-search"></i>
      <input
        type="text"
        class="search-input"
        placeholder="Search services..."
        value={searchQuery}
        on:input={handleSearch}
        on:focus={() => { if (searchQuery) showResults = true; }}
      />
      {#if searchQuery}
        <button class="clear-btn" on:click={clearSearch}>
          <i class="fas fa-times"></i>
        </button>
      {/if}
    </div>

    {#if showResults && filteredItems.length > 0}
      <div class="search-results">
        {#each filteredItems as item}
          <button 
            class="search-result-item"
            on:click={() => navigateTo(item.url, item.target)}
          >
            <div class="result-icon">
              {#if item.logo}
                <img src={item.logo} alt={item.name} />
              {:else}
                <i class="fas fa-link"></i>
              {/if}
            </div>
            <div class="result-content">
              <div class="result-name">{item.name}</div>
              <div class="result-subtitle">{item.subtitle || item.serviceName}</div>
            </div>
            {#if item.tag}
              <span class="result-tag">{item.tag}</span>
            {/if}
          </button>
        {/each}
      </div>
    {/if}
  </div>

  <button class="style-toggle" on:click={toggleStyle}>
    <i class="fas {viewStyle === 'cards' ? 'fa-list' : 'fa-th-large'}"></i>
    <span>{viewStyle === 'cards' ? 'List View' : 'Card View'}</span>
  </button>
</div>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
  }

  .search-container {
    flex: 1;
    max-width: 400px;
    position: relative;
  }

  .search-box {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 16px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    transition: all 0.3s ease;
  }

  .search-box:focus-within {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(102, 126, 234, 0.5);
  }

  .search-box i {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #e4e4e4;
    font-size: 0.95rem;
  }

  .search-input::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .clear-btn {
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
  }

  .clear-btn:hover {
    color: #e4e4e4;
  }

  .search-results {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    background: rgba(30, 30, 50, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    overflow: hidden;
    z-index: 1000;
    max-height: 400px;
    overflow-y: auto;
    backdrop-filter: blur(10px);
  }

  .search-result-item {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 12px 16px;
    background: transparent;
    border: none;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
  }

  .search-result-item:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .result-icon {
    width: 36px;
    height: 36px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(102, 126, 234, 0.2);
    border-radius: 8px;
  }

  .result-icon img {
    width: 28px;
    height: 28px;
    object-fit: contain;
    border-radius: 4px;
  }

  .result-icon i {
    color: #667eea;
    font-size: 1rem;
  }

  .result-content {
    flex: 1;
    min-width: 0;
  }

  .result-name {
    font-size: 0.95rem;
    font-weight: 500;
    color: #ffffff;
  }

  .result-subtitle {
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.5);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .result-tag {
    padding: 2px 8px;
    background: rgba(102, 126, 234, 0.3);
    color: #a5b4fc;
    font-size: 0.7rem;
    font-weight: 500;
    border-radius: 4px;
    text-transform: uppercase;
  }

  .style-toggle {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    color: #e4e4e4;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.9rem;
    white-space: nowrap;
  }

  .style-toggle:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: translateY(-2px);
  }

  @media (max-width: 768px) {
    .toolbar {
      flex-direction: column;
    }

    .search-container {
      max-width: 100%;
      width: 100%;
    }

    .style-toggle span {
      display: none;
    }

    .style-toggle {
      padding: 10px 14px;
    }
  }
</style>
