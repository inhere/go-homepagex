<script>
  export let title = 'Home Dashboard';
  export let subtitle = '';
  export let logo = '';
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
    padding: 30px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 16px;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .logo {
    width: 64px;
    height: 64px;
    object-fit: contain;
    border-radius: 12px;
  }

  .logo-placeholder {
    width: 64px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    font-size: 2rem;
    color: white;
  }

  .title-section {
    flex: 1;
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

  @media (max-width: 768px) {
    .header-content {
      flex-direction: column;
      text-align: center;
      padding: 20px;
    }

    .title {
      font-size: 1.8rem;
    }

    .logo, .logo-placeholder {
      width: 48px;
      height: 48px;
    }
  }
</style>
