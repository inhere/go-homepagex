<script>
  import ServiceItem from './ServiceItem.svelte';
  
  export let service = { name: '', icon: '', items: [] };
  export let style = 'cards';
</script>

<div class="service-group {style}">
  <div class="group-header">
    <i class="{service.icon || 'fas fa-folder'}"></i>
    <h2>{service.name}</h2>
  </div>
  
  <div class="items-container">
    {#each service.items || [] as item}
      <ServiceItem {item} {style} />
    {/each}
  </div>
</div>

<style>
  .service-group {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 16px;
    padding: 20px;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }

  .service-group:hover {
    transform: translateY(-4px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  }

  .group-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .group-header i {
    font-size: 1.5rem;
    color: #667eea;
    width: 32px;
    text-align: center;
  }

  .group-header h2 {
    font-size: 1.3rem;
    font-weight: 600;
    color: #ffffff;
    margin: 0;
  }

  .items-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  /* List view specific styles */
  .list.service-group {
    padding: 16px;
  }

  .list .group-header {
    margin-bottom: 12px;
    padding-bottom: 8px;
  }

  .list .group-header i {
    font-size: 1.2rem;
  }

  .list .group-header h2 {
    font-size: 1.1rem;
  }

  .list .items-container {
    gap: 8px;
  }
</style>
