<script>
  export let tags = [];
  export let selectedTag = '';
  export let onSelectTag = () => {};

  function getTagCount(tag) {
    return tag.count || 0;
  }
</script>

<div class="tag-filter">
  <div class="tag-header">
    <i class="fas fa-tags"></i>
    <span>标签筛选</span>
  </div>
  <div class="tag-list">
    <button
      class="tag-btn"
      class:active={selectedTag === ''}
      on:click={() => onSelectTag('')}
    >
      <span class="tag-name">全部</span>
      <span class="tag-count"></span>
    </button>
    {#each tags as tag}
      <button
        class="tag-btn"
        class:active={selectedTag === tag.name}
        on:click={() => onSelectTag(selectedTag === tag.name ? '' : tag.name)}
      >
        <span class="tag-name">{tag.name}</span>
        <span class="tag-count">{tag.count}</span>
      </button>
    {/each}
  </div>
</div>

<style>
  .tag-filter {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .tag-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.8);
  }

  .tag-header i {
    color: var(--theme-primary);
  }

  .tag-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .tag-btn {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 8px 12px;
    background: transparent;
    border: 1px solid transparent;
    border-radius: 8px;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.85rem;
  }

  .tag-btn:hover {
    background: rgba(255, 255, 255, 0.08);
    color: #ffffff;
  }

  .tag-btn.active {
    background: var(--theme-primary-rgba);
    border-color: var(--theme-primary);
    color: var(--theme-primary);
  }

  .tag-name {
    flex: 1;
    text-align: left;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    font-weight: 500;
  }

  .tag-count {
    font-size: 0.75rem;
    padding: 2px 8px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 10px;
    color: rgba(255, 255, 255, 0.5);
    font-weight: 500;
    min-width: 24px;
    text-align: center;
  }

  .tag-btn.active .tag-count {
    background: var(--theme-primary);
    color: #ffffff;
  }

  @media (max-width: 768px) {
    .tag-filter {
      display: none;
    }
  }
</style>
