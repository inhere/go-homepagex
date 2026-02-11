<script>
  export let navs = [];
  export let currentPath = '/';
</script>

<nav class="navbar">
  <div class="nav-list">
    {#each navs as nav}
      <a
        href={nav.url}
        class="nav-item"
        class:active={currentPath === nav.url}
      >
        <i class={nav.icon}></i>
        <span>{nav.name}</span>
      </a>
    {/each}
  </div>
</nav>

<style>
  .navbar {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    padding: 8px;
    margin-bottom: 20px;
  }

  .nav-list {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    color: rgba(255, 255, 255, 0.7);
    text-decoration: none;
    border-radius: 8px;
    transition: all 0.3s ease;
    font-size: 0.95rem;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #ffffff;
  }

  .nav-item.active {
    background: var(--theme-primary-rgba);
    color: var(--theme-primary);
  }

  .nav-item i {
    font-size: 1rem;
  }

  @media (max-width: 768px) {
    .nav-list {
      justify-content: center;
    }

    .nav-item span {
      display: none;
    }

    .nav-item {
      padding: 10px 14px;
    }
  }
</style>
