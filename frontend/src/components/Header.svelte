<script>
  import { userInfo } from '../stores.js';

  export let title = 'Home Dashboard';
  export let subtitle = '';
  export let logo = '';

  function logout() {
    sessionStorage.setItem('loggedOut', '1');
    window.location.href = '/';
  }

  function handleLoginClick() {
    sessionStorage.removeItem('loggedOut');
  }

  $: loginUrl = '/api/auth?return=' + encodeURIComponent(
    typeof window !== 'undefined' ? window.location.pathname : '/'
  );
</script>

<header class="header">
  <div class="header-content">
    {#if logo}
      <img src={logo} alt="Logo" class="logo" />
    {:else}
      <div class="logo-placeholder">
        <i class="fas fa-home"></i>
      </div>
    {/if}

    <div class="title-section">
      <h1 class="title">{title}</h1>
      {#if subtitle}
        <p class="subtitle">{subtitle}</p>
      {/if}
    </div>

    <div class="user-section">
      {#if $userInfo}
        <span class="user-info">
          <i class="fas fa-user-circle"></i>
          <span class="username">{$userInfo.username}</span>
        </span>
        <button class="btn-auth" on:click={logout} title="退出登录">
          <i class="fas fa-sign-out-alt"></i>
          <span>退出</span>
        </button>
      {:else}
        <a class="btn-auth login" href={loginUrl} on:click={handleLoginClick} title="登录">
          <i class="fas fa-sign-in-alt"></i>
          <span>登录</span>
        </a>
      {/if}
    </div>
  </div>
</header>

<style>
  .header {
    margin-bottom: 30px;
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: 20px;
    padding: 20px 30px;
    backdrop-filter: blur(10px);
  }

  .logo {
    width: 64px;
    height: 64px;
    object-fit: contain;
    border-radius: 12px;
    flex-shrink: 0;
  }

  .logo-placeholder {
    width: 64px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--theme-primary), var(--theme-secondary));
    border-radius: 12px;
    font-size: 2rem;
    color: white;
    flex-shrink: 0;
  }

  .title-section {
    flex: 1;
    min-width: 0;
  }

  .title {
    font-size: 2.5rem;
    font-weight: 700;
    color: #ffffff;
    margin: 0;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .subtitle {
    font-size: 1.1rem;
    color: rgba(255, 255, 255, 0.7);
    margin: 8px 0 0 0;
  }

  /* 用户信息区域 */
  .user-section {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-shrink: 0;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 14px;
    background: rgba(255, 255, 255, 0.08);
    border-radius: 20px;
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.9rem;
  }

  .user-info i {
    font-size: 1.1rem;
    color: var(--theme-primary);
  }

  .btn-auth {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    text-decoration: none;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.85);
  }

  .btn-auth:hover {
    background: rgba(255, 255, 255, 0.15);
    color: #ffffff;
    border-color: rgba(255, 255, 255, 0.4);
  }

  .btn-auth.login {
    border-color: var(--theme-primary);
    color: var(--theme-primary);
  }

  .btn-auth.login:hover {
    background: var(--theme-primary-rgba);
  }

  @media (max-width: 768px) {
    .header-content {
      flex-wrap: wrap;
      padding: 16px 20px;
      gap: 12px;
    }

    .title-section {
      flex: 1 1 auto;
    }

    .title {
      font-size: 1.8rem;
    }

    .logo, .logo-placeholder {
      width: 48px;
      height: 48px;
    }

    .user-section {
      width: 100%;
      justify-content: flex-end;
    }

    .btn-auth span {
      display: none;
    }

    .btn-auth {
      padding: 8px 12px;
    }
  }
</style>
