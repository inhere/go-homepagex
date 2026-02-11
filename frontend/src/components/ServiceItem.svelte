<script>
  export let item = { name: '', url: '', logo: '', subtitle: '', tag: '' };
  export let style = 'cards';

  function handleClick() {
    if (item.target === '_blank') {
      window.open(item.url, '_blank');
    } else {
      window.location.href = item.url;
    }
  }

  function handleKeydown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleClick();
    }
  }
</script>

{#if style === 'cards'}
  <div class="service-item card" on:click={handleClick} on:keydown={handleKeydown} role="button" tabindex="0">
    <div class="item-logo">
      {#if item.logo}
        <img src={item.logo} alt={item.name} />
      {:else}
        <div class="logo-fallback">
          <i class="fas fa-link"></i>
        </div>
      {/if}
    </div>
    <div class="item-content">
      <div class="item-header">
        <h3>{item.name}</h3>
        {#if item.tag}
          <span class="tag">{item.tag}</span>
        {/if}
      </div>
      {#if item.subtitle}
        <p class="subtitle">{item.subtitle}</p>
      {/if}
    </div>
  </div>
{:else}
  <div class="service-item list" on:click={handleClick} on:keydown={handleKeydown} role="button" tabindex="0">
    <div class="item-logo">
      {#if item.logo}
        <img src={item.logo} alt={item.name} />
      {:else}
        <div class="logo-fallback">
          <i class="fas fa-link"></i>
        </div>
      {/if}
    </div>
    <div class="item-content">
      <h3>{item.name}</h3>
      {#if item.subtitle}
        <p class="item-subtitle">{item.subtitle}</p>
      {/if}
    </div>
    {#if item.tag}
      <span class="tag">{item.tag}</span>
    {/if}
  </div>
{/if}

<style>
  .service-item {
    cursor: pointer;
    transition: all 0.3s ease;
  }

  /* Card style */
  .service-item.card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }

  .service-item.card:hover {
    background: rgba(255, 255, 255, 0.1);
    border-color: var(--theme-primary-rgba);
    transform: translateX(4px);
  }

  /* List style */
  .service-item.list {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }

  .service-item.list:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: var(--theme-primary-rgba);
  }

  .item-logo {
    flex-shrink: 0;
  }

  .item-logo img {
    width: 48px;
    height: 48px;
    object-fit: contain;
    border-radius: 8px;
  }

  .logo-fallback {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--theme-primary), var(--theme-secondary));
    border-radius: 8px;
    color: white;
    font-size: 1.2rem;
  }

  .list .item-logo img,
  .list .logo-fallback {
    width: 36px;
    height: 36px;
  }

  .list .logo-fallback {
    font-size: 1rem;
  }

  .item-content {
    flex: 1;
    min-width: 0;
  }

  .item-header {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .item-content h3 {
    font-size: 1.1rem;
    font-weight: 600;
    color: #ffffff;
    margin: 0;
  }

  .list .item-content h3 {
    font-size: 1rem;
  }

  .item-subtitle {
    font-size: 0.85rem;
    color: rgba(255, 255, 255, 0.6);
    margin: 4px 0 0 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .list .item-subtitle {
    font-size: 0.8rem;
    margin-top: 2px;
  }

  .tag {
    display: inline-block;
    padding: 2px 8px;
    background: var(--theme-primary-rgba);
    color: var(--theme-primary);
    font-size: 0.7rem;
    font-weight: 500;
    border-radius: 4px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .list .tag {
    margin-left: auto;
  }

  /* Focus styles for accessibility */
  .service-item:focus {
    outline: 2px solid var(--theme-primary);
    outline-offset: 2px;
  }
</style>
