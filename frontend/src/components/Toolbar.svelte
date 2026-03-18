<script>
  import { viewStyle, themes, currentTheme, getThemeColors, userInfo, currentRoute } from '../stores.js';

  export let onSearch = () => {};
  export let onOpenEditor = () => {};

  let searchQuery = '';
  let showThemeDropdown = false;

  $: themeColors = getThemeColors($currentTheme);
  $: currentThemeName = themes.find(t => t.id === $currentTheme)?.name || '默认主题';

  function handleSearch(event) {
    searchQuery = event.target.value;
    onSearch(searchQuery);
  }

  function clearSearch() {
    searchQuery = '';
    onSearch('');
  }

  function toggleStyle() {
    viewStyle.update(style => style === 'cards' ? 'list' : 'cards');
  }

  function selectTheme(themeId) {
    currentTheme.set(themeId);
    showThemeDropdown = false;
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      showThemeDropdown = false;
      document.querySelector('.search-input').blur();
    }
  }

  function handleClickOutside(event) {
    if (!event.target.closest('.theme-selector')) {
      showThemeDropdown = false;
    }
  }

  // 前端复用与后端类似的路径匹配逻辑，判断当前路由是否具备 rw 权限
  function matchPermForPath(pattern, perm, routePath) {
    if (!pattern) {
      return null;
    }

    let p = pattern;
    let per = perm || 'ro';
    let reqPath = routePath || '/';

    if (p !== '*' && !p.startsWith('/')) {
      p = '/' + p;
    }
    if (!reqPath.startsWith('/')) {
      reqPath = '/' + reqPath;
    }

    // 通配后缀：*, /*, /prefix*, /prefix/* → 前缀匹配
    if (p.endsWith('*')) {
      const prefix = p.slice(0, -1);
      return reqPath.startsWith(prefix) ? per : null;
    }

    // 精确匹配 或 子路径匹配
    if (reqPath === p || reqPath.startsWith(p + '/')) {
      return per;
    }

    return null;
  }

  function hasWritePermissionForRoute(perms, routePath) {
    if (!perms || !perms.length) return false;

    for (const rule of perms) {
      const pattern = rule.path || rule.Path;
      const perm = rule.perm || rule.Permission;
      const matchedPerm = matchPermForPath(pattern, perm, routePath);
      if (!matchedPerm) continue;

      // 后端解析时已将禁止规则排在前面，这里遇到 no 直接视为无权
      if (matchedPerm === 'no') {
        return false;
      }

      // 找到第一条匹配规则，且为 rw，则认为当前路由可编辑
      if (matchedPerm === 'rw') {
        return true;
      }

      // 匹配到 ro，则继续找后续是否有更具体的 rw 规则
    }

    return false;
  }

  $: canEditCurrentRoute =
    !!$userInfo && hasWritePermissionForRoute($userInfo.permissions || [], $currentRoute);
</script>

<svelte:window on:keydown={handleKeydown} on:click={handleClickOutside} />

<div class="toolbar">
  <div class="search-box">
    <i class="fas fa-search"></i>
    <input
      type="text"
      class="search-input"
      placeholder="搜索服务..."
      value={searchQuery}
      on:input={handleSearch}
    />
    {#if searchQuery}
      <button class="clear-btn" on:click={clearSearch} title="清除">
        <i class="fas fa-times"></i>
      </button>
    {/if}
  </div>

  <div class="toolbar-actions">
    <div class="theme-selector">
      <button
        class="theme-btn"
        on:click|stopPropagation={() => showThemeDropdown = !showThemeDropdown}
      >
        <div
          class="theme-indicator"
          style="background: linear-gradient(135deg, {themeColors.primary}, {themeColors.secondary})"
        ></div>
        <span class="theme-name">{currentThemeName}</span>
        <i class="fas fa-chevron-down" class:open={showThemeDropdown}></i>
      </button>

      {#if showThemeDropdown}
        <div class="theme-dropdown">
          {#each themes as theme}
            {@const colors = getThemeColors(theme.id)}
            <button
              class="theme-option"
              class:active={$currentTheme === theme.id}
              on:click={() => selectTheme(theme.id)}
            >
              <div
                class="theme-preview"
                style="background: linear-gradient(135deg, {colors.primary}, {colors.secondary})"
              ></div>
              <span>{theme.name}</span>
              {#if $currentTheme === theme.id}
                <i class="fas fa-check"></i>
              {/if}
            </button>
          {/each}
        </div>
      {/if}
    </div>

    <button class="style-toggle" on:click={toggleStyle} title="切换视图">
      <i class="fas {$viewStyle === 'cards' ? 'fa-list' : 'fa-th-large'}"></i>
      <span>{$viewStyle === 'cards' ? '列表' : '卡片'}</span>
    </button>

    {#if canEditCurrentRoute}
      <button class="edit-btn" type="button" on:click={onOpenEditor} title="编辑当前页面 YAML">
        <i class="fas fa-edit"></i>
        <span>编辑</span>
      </button>
    {/if}
  </div>
</div>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
  }

  .search-box {
    flex: 1;
    max-width: 400px;
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
    border-color: var(--theme-primary-rgba);
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

  .toolbar-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .theme-selector {
    position: relative;
  }

  .theme-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    color: #e4e4e4;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.9rem;
  }

  .theme-btn:hover {
    background: rgba(255, 255, 255, 0.15);
  }

  .theme-indicator {
    width: 20px;
    height: 20px;
    border-radius: 4px;
    flex-shrink: 0;
  }

  .theme-name {
    white-space: nowrap;
  }

  .theme-btn i {
    font-size: 0.75rem;
    transition: transform 0.2s;
  }

  .theme-btn i.open {
    transform: rotate(180deg);
  }

  .theme-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 8px;
    background: rgba(30, 30, 50, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    overflow: hidden;
    z-index: 1000;
    min-width: 160px;
    backdrop-filter: blur(10px);
  }

  .theme-option {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 10px 14px;
    background: transparent;
    border: none;
    cursor: pointer;
    text-align: left;
    color: #e4e4e4;
    font-size: 0.9rem;
    transition: background 0.2s;
  }

  .theme-option:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .theme-option.active {
    background: var(--theme-primary-rgba);
  }

  .theme-preview {
    width: 20px;
    height: 20px;
    border-radius: 4px;
    flex-shrink: 0;
  }

  .theme-option i {
    margin-left: auto;
    color: var(--theme-primary);
    font-size: 0.85rem;
  }

  .style-toggle {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    color: #e4e4e4;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.9rem;
  }

  .style-toggle:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: translateY(-2px);
  }

  .edit-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 9px 16px;
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.3);
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.95);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .edit-btn i {
    font-size: 0.9rem;
  }

  .edit-btn:hover {
    background: var(--theme-primary-rgba);
    border-color: var(--theme-primary);
    transform: translateY(-1px);
  }

  @media (max-width: 768px) {
    .toolbar {
      flex-direction: column;
    }

    .search-box {
      max-width: 100%;
      width: 100%;
    }

    .toolbar-actions {
      width: 100%;
      justify-content: space-between;
      flex-wrap: wrap;
    }

    .theme-name {
      display: none;
    }

    .theme-btn {
      padding: 10px 12px;
    }

    .style-toggle span {
      display: none;
    }

    .style-toggle {
      padding: 10px 12px;
    }

    .edit-btn span {
      display: none;
    }

    .edit-btn {
      padding: 9px 10px;
    }

    .theme-dropdown {
      right: 0;
      min-width: 140px;
    }
  }
</style>
