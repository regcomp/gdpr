function contentManager() {
  return {
    currentComponent: 'empty',

    loadComponent(component) {
      this.currentComponent = component;
    },

    loadDashboard() {
      this.loadComponent('dashboard');
    },

    loadAnalytics() {
      this.loadComponent('analytics');
    },

    isActive(component) {
      return this.currentComponent === component;
    }
  }
}

document.addEventListener('alpine:init', () => {
  Alpine.data('contentManager', contentManager);
});
