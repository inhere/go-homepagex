<script>
  import { onMount } from 'svelte';
  import Header from './components/Header.svelte';
  import Navbar from './components/Navbar.svelte';
  import Toolbar from './components/Toolbar.svelte';
  import ServiceGroup from './components/ServiceGroup.svelte';
  import { pageConfig, currentRoute, viewStyle } from './stores.js';

  let loading = true;
  let error = null;
  let searchQuery = '';
  let filteredServices = [];

  $: {
    if (searchQuery.trim()) {
      filteredServices = ($pageConfig.services || []).map(service => ({
        ...service,
        items: (service.items || []).filter(item =>
          item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          item.subtitle?.toLowerCase().includes(searchQuery.toLowerCase()) ||
          item.tag?.toLowerCase().includes(searchQuery.toLowerCase())
        )
      })).filter(service => service.items.length > 0);
    } else {
      filteredServices = $pageConfig.services || [];
    }
  }

  // 获取当前路由
  function getRoute() {
    const path = window.location.pathname;
    return path === '/' ? '/' : path;
  }

  // 加载页面配置
  async function loadConfig() {
    loading = true;
    error = null;

    try {
      const route = getRoute();
      currentRoute.set(route);

      const response = await fetch(`/api/page${route === '/' ? '' : route}`);
      const result = await response.json();

      if (!result.success) {
        throw new Error(result.error || 'Failed to load config');
      }

      pageConfig.set(result.data);
      viewStyle.set(result.data.style || 'cards');
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadConfig();

    // 监听浏览器前进后退
    window.addEventListener('popstate', loadConfig);

    return () => {
      window.removeEventListener('popstate', loadConfig);
    };
  });

  function handleSearch(query) {
    searchQuery = query;
  }
</script>

<svelte:head>
  <title>{$pageConfig?.title || 'Home Dashboard'}</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
</svelte:head>

<main class="app {$viewStyle}">
  {#if loading}
    <div class="loading">
      <i class="fas fa-spinner fa-spin"></i>
      <span>Loading...</span>
    </div>
  {:else if error}
    <div class="error">
      <i class="fas fa-exclamation-circle"></i>
      <span>{error}</span>
    </div>
  {:else}
    <Header
      title={$pageConfig.title}
      subtitle={$pageConfig.subtitle}
      logo={$pageConfig.logo}
    />

    {#if $pageConfig.navs && $pageConfig.navs.length > 0}
      <Navbar navs={$pageConfig.navs} currentPath={$currentRoute} />
    {/if}

    <Toolbar
      services={$pageConfig.services || []}
      bind:viewStyle={$viewStyle}
      onSearch={handleSearch}
    />

    <div class="services-container" style="--columns: {$pageConfig.columns || '3'}">
      {#each filteredServices as service}
        <ServiceGroup {service} style={$viewStyle} />
      {/each}
    </div>

    {#if $pageConfig.footer}
      <footer class="footer">
        {$pageConfig.footer}
      </footer>
    {/if}
  {/if}
</main>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
    min-height: 100vh;
    color: #e4e4e4;
  }

  .app {
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
    min-height: 100vh;
  }

  .loading, .error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    gap: 16px;
    font-size: 1.2rem;
  }

  .loading i, .error i {
    font-size: 3rem;
  }

  .error {
    color: #e74c3c;
  }

  .services-container {
    display: grid;
    grid-template-columns: repeat(var(--columns, 3), 1fr);
    gap: 24px;
  }

  /* Responsive */
  @media (max-width: 1200px) {
    .services-container {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 768px) {
    .services-container {
      grid-template-columns: 1fr;
    }
  }

  /* List view styles */
  .list .services-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .footer {
    text-align: center;
    padding: 40px 20px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
  }
</style>
