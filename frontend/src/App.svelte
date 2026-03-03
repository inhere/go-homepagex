<script>
  import { onMount } from 'svelte';
  import Header from './components/Header.svelte';
  import Navbar from './components/Navbar.svelte';
  import Toolbar from './components/Toolbar.svelte';
  import TagFilter from './components/TagFilter.svelte';
  import ServiceGroup from './components/ServiceGroup.svelte';
  import { pageConfig, currentRoute, viewStyle, currentTheme, getThemeColors, userInfo } from './stores.js';

  let loading = true;
  let error = null;
  let searchQuery = '';
  let selectedTag = '';
  let filteredServices = [];
  let allTags = [];

  $: themeColors = getThemeColors($currentTheme);
  $: themeVars = `
    --theme-primary: ${themeColors.primary};
    --theme-secondary: ${themeColors.secondary};
    --theme-accent: ${themeColors.accent};
    --theme-background: ${themeColors.background};
    --theme-primary-rgba: ${themeColors.primary}33;
    --theme-secondary-rgba: ${themeColors.secondary}33;
    --theme-accent-rgba: ${themeColors.accent}33;
    --theme-background-rgba: ${themeColors.background}dd;
  `;

  function matchesAllKeywords(text, keywords) {
    if (!keywords.length) return true;
    const lowerText = text.toLowerCase();
    return keywords.every(keyword => lowerText.includes(keyword));
  }

  function collectTags(services) {
    const tagCount = {};
    (services || []).forEach(service => {
      (service.items || []).forEach(item => {
        (item.tags || []).forEach(tag => {
          tagCount[tag] = (tagCount[tag] || 0) + 1;
        });
      });
    });
    return Object.entries(tagCount)
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count);
  }

  function filterByTag(services, tag) {
    if (!tag) return services;
    return services.map(service => ({
      ...service,
      items: (service.items || []).filter(item =>
        (item.tags || []).includes(tag)
      )
    })).filter(service => service.items.length > 0);
  }

  $: {
    const baseServices = $pageConfig.services || [];
    allTags = collectTags(baseServices);

    let result = baseServices;

    if (searchQuery.length > 1) {
      const keywords = searchQuery.split(/\s+/).filter(k => k.length > 0);
      result = result.map(service => ({
        ...service,
        items: (service.items || []).filter(item => {
          const searchFields = [
            item.name,
            item.subtitle,
            ...(item.tags || []),
            service.name
          ].filter(Boolean).join(' ');
          return matchesAllKeywords(searchFields, keywords);
        })
      })).filter(service => service.items.length > 0);
    }

    if (selectedTag) {
      result = filterByTag(result, selectedTag);
    }

    filteredServices = result;
  }

  // 获取当前路由
  function getRoute() {
    const path = window.location.pathname;
    return path === '/' ? '/' : path;
  }

  // 检查并处理 URL 中的 refresh 参数
  function getRefreshParam() {
    const urlParams = new URLSearchParams(window.location.search);
    const refresh = urlParams.get('refresh');
    if (refresh === 'true') {
      // 移除 URL 中的 refresh 参数
      urlParams.delete('refresh');
      const newUrl = window.location.pathname + (urlParams.toString() ? '?' + urlParams.toString() : '');
      history.replaceState({}, '', newUrl);
      return true;
    }
    return false;
  }

  // 加载页面配置
  async function loadConfig(route) {
    const currentRoutePath = route || getRoute();
    loading = true;
    error = null;
    searchQuery = '';
    selectedTag = '';

    const shouldRefresh = getRefreshParam();
    const apiPath = shouldRefresh
      ? `/api/page${currentRoutePath === '/' ? '' : currentRoutePath}?refresh=true`
      : `/api/page${currentRoutePath === '/' ? '' : currentRoutePath}`;

    try {
      currentRoute.set(currentRoutePath);

      const fetchOptions = {};
      if (sessionStorage.getItem('loggedOut') === '1') {
        // 退出登录后发送空凭据，覆盖浏览器缓存的凭据，以游客身份访问
        fetchOptions.headers = { 'Authorization': 'Basic Og==' };
      }
      const response = await fetch(apiPath, fetchOptions);

      if (response.status === 401) {
        // 需要认证：清除 loggedOut 标记，跳转到认证端点，认证成功后重定向回当前页面
        sessionStorage.removeItem('loggedOut');
        window.location.href = '/api/auth?return=' + encodeURIComponent(window.location.pathname);
        return;
      }

      const result = await response.json();

      if (!result.success) {
        throw new Error(result.error || 'Failed to load config');
      }

      pageConfig.set(result.data);
      userInfo.set(result.data.user_info || null);
      if (!localStorage.getItem('viewStyle')) {
        viewStyle.set(result.data.style || 'cards');
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  // 处理导航（单页应用）
  function handleNavigate(route) {
    if (route !== getRoute()) {
      history.pushState({}, '', route);
      loadConfig(route);
    }
  }

  onMount(() => {
    loadConfig();

    // 监听浏览器前进后退
    window.addEventListener('popstate', () => loadConfig());

    return () => {
      window.removeEventListener('popstate', loadConfig);
    };
  });

  function handleSearch(query) {
    searchQuery = query;
  }

  function handleSelectTag(tag) {
    selectedTag = tag;
  }
</script>

<svelte:head>
  <title>{$pageConfig?.title || 'Home Dashboard'}</title>
</svelte:head>

<div class="theme-wrapper" style={themeVars}>
  <main class="app {$viewStyle} theme-{$currentTheme}">
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
        <Navbar navs={$pageConfig.navs} currentPath={$currentRoute} onNavigate={handleNavigate} />
      {/if}

      <Toolbar onSearch={handleSearch} />

      <div class="main-content">
        <aside class="sidebar">
          <TagFilter
            tags={allTags}
            {selectedTag}
            onSelectTag={handleSelectTag}
          />
        </aside>

        <div class="services-wrapper">
          {#if filteredServices.length === 0}
            <div class="no-results">
              <i class="fas fa-search"></i>
              <p>没有找到匹配的服务</p>
            </div>
          {:else}
            <div class="services-container" style="--columns: {$pageConfig.columns || '3'}">
              {#each filteredServices as service}
                <ServiceGroup {service} style={$viewStyle} />
              {/each}
            </div>
          {/if}
        </div>
      </div>

      {#if $pageConfig.footer}
        <footer class="footer">
          {$pageConfig.footer}
        </footer>
      {/if}
    {/if}
  </main>
</div>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: var(--theme-primary);
    min-height: 100vh;
    color: var(--theme-background);
    transition: background 0.5s ease, color 0.5s ease;
  }

  :global(:root) {
    --theme-primary: #1a2332;
    --theme-secondary: #2d8b8b;
    --theme-accent: #a8dadc;
    --theme-background: #f1faee;
    --theme-primary-rgba: rgba(26, 35, 50, 0.2);
    --theme-secondary-rgba: rgba(45, 139, 139, 0.2);
    --theme-accent-rgba: rgba(168, 218, 220, 0.2);
    --theme-background-rgba: rgba(241, 250, 238, 0.87);
  }

  .theme-wrapper {
    min-height: 100vh;
    background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-secondary) 50%, var(--theme-accent) 100%);
    transition: background 0.5s ease;
  }

  .app {
    max-width: 1600px;
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

  .main-content {
    display: flex;
    gap: 20px;
  }

  .sidebar {
    width: 200px;
    flex-shrink: 0;
  }

  .services-wrapper {
    flex: 1;
    min-width: 0;
  }

  .no-results {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    color: rgba(255, 255, 255, 0.4);
    gap: 16px;
  }

  .no-results i {
    font-size: 3rem;
    opacity: 0.5;
  }

  .no-results p {
    font-size: 1.1rem;
  }

  .services-container {
    display: grid;
    grid-template-columns: repeat(var(--columns, 3), minmax(280px, 1fr));
    gap: 24px;
  }

  /* Responsive */
  @media (max-width: 1400px) {
    .services-container {
      grid-template-columns: repeat(2, minmax(280px, 1fr));
    }
  }

  @media (max-width: 900px) {
    .main-content {
      flex-direction: column;
    }

    .sidebar {
      width: 100%;
    }

    .services-container {
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
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
