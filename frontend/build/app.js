// Home Dashboard Frontend
// Vanilla JavaScript implementation

class HomeDashboard {
  constructor() {
    this.pageConfig = {};
    this.currentRoute = '/';
    this.viewStyle = 'cards';
    this.app = document.getElementById('app');
    
    this.init();
  }

  async init() {
    this.showLoading();
    
    try {
      await this.loadConfig();
      this.render();
    } catch (error) {
      this.showError(error.message);
    }
  }

  showLoading() {
    this.app.innerHTML = `
      <div class="loading">
        <i class="fas fa-spinner fa-spin"></i>
        <span>Loading...</span>
      </div>
    `;
  }

  showError(message) {
    this.app.innerHTML = `
      <div class="error">
        <i class="fas fa-exclamation-circle"></i>
        <span>${message}</span>
      </div>
    `;
  }

  async loadConfig() {
    const route = window.location.pathname;
    this.currentRoute = route === '/' ? '/' : route;
    
    const response = await fetch(`/api/page${this.currentRoute === '/' ? '' : this.currentRoute}`);
    const result = await response.json();
    
    if (!result.success) {
      throw new Error(result.error || 'Failed to load config');
    }
    
    this.pageConfig = result.data;
    this.viewStyle = this.pageConfig.style || 'cards';
  }

  toggleStyle() {
    this.viewStyle = this.viewStyle === 'cards' ? 'list' : 'cards';
    this.render();
  }

  render() {
    document.title = this.pageConfig.title || 'Home Dashboard';
    
    this.app.className = `app ${this.viewStyle}`;
    
    this.app.innerHTML = `
      ${this.renderHeader()}
      ${this.renderToolbar()}
      ${this.renderServices()}
      ${this.renderFooter()}
    `;
    
    // Attach event listeners
    const toggleBtn = this.app.querySelector('.style-toggle');
    if (toggleBtn) {
      toggleBtn.addEventListener('click', () => this.toggleStyle());
    }
    
    // Attach click handlers to service items
    const items = this.app.querySelectorAll('.service-item');
    items.forEach(item => {
      item.addEventListener('click', (e) => {
        const url = item.dataset.url;
        const target = item.dataset.target;
        if (url) {
          if (target === '_blank') {
            window.open(url, '_blank');
          } else {
            window.location.href = url;
          }
        }
      });
    });
  }

  renderHeader() {
    const { title, subtitle, logo } = this.pageConfig;
    
    return `
      <header class="header">
        <div class="header-content">
          ${logo 
            ? `<img src="${logo}" alt="Logo" class="logo" />`
            : `<div class="logo-placeholder"><i class="fas fa-home"></i></div>`
          }
          <div class="title-section">
            <h1 class="title">${title || 'Home Dashboard'}</h1>
            ${subtitle ? `<p class="subtitle">${subtitle}</p>` : ''}
          </div>
        </div>
      </header>
    `;
  }

  renderToolbar() {
    return `
      <div class="toolbar">
        <button class="style-toggle">
          <i class="fas ${this.viewStyle === 'cards' ? 'fa-list' : 'fa-th-large'}"></i>
          <span>${this.viewStyle === 'cards' ? 'List View' : 'Card View'}</span>
        </button>
      </div>
    `;
  }

  renderServices() {
    const { services = [], columns = '3' } = this.pageConfig;
    
    return `
      <div class="services-container" style="--columns: ${columns}">
        ${services.map(service => this.renderServiceGroup(service)).join('')}
      </div>
    `;
  }

  renderServiceGroup(service) {
    return `
      <div class="service-group">
        <div class="group-header">
          <i class="${service.icon || 'fas fa-folder'}"></i>
          <h2>${service.name}</h2>
        </div>
        <div class="items-container">
          ${service.items.map(item => this.renderServiceItem(item)).join('')}
        </div>
      </div>
    `;
  }

  renderServiceItem(item) {
    return `
      <div class="service-item ${this.viewStyle}" 
           data-url="${item.url}" 
           data-target="${item.target || '_self'}"
           role="button" 
           tabindex="0">
        <div class="item-logo">
          ${item.logo 
            ? `<img src="${item.logo}" alt="${item.name}" />`
            : `<div class="logo-fallback"><i class="fas fa-link"></i></div>`
          }
        </div>
        <div class="item-content">
          ${this.viewStyle === 'cards' 
            ? `
              <div class="item-header">
                <h3>${item.name}</h3>
                ${item.tag ? `<span class="tag">${item.tag}</span>` : ''}
              </div>
            `
            : `<h3>${item.name}</h3>`
          }
          ${item.subtitle ? `<p class="item-subtitle">${item.subtitle}</p>` : ''}
        </div>
        ${this.viewStyle === 'list' && item.tag ? `<span class="tag">${item.tag}</span>` : ''}
      </div>
    `;
  }

  renderFooter() {
    const { footer } = this.pageConfig;
    
    if (!footer) return '';
    
    return `<footer class="footer">${footer}</footer>`;
  }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new HomeDashboard();
});

// Handle browser back/forward
window.addEventListener('popstate', () => {
  new HomeDashboard();
});
