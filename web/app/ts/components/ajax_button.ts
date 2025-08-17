interface AjaxButtonData {
  $el: HTMLElement;
  fetchAJAXContent(): void;
}

// The Alpine component factory function
const createAjaxButtonComponent = (): AjaxButtonData => ({
  $el: null as any, // This will be set by Alpine.js at runtime

  fetchAJAXContent(): void {
    const endpoint: string | undefined = this.$el.dataset.endpoint;

    if (!endpoint) {
      console.error('No endpoint specified in data-endpoint attribute');
      return;
    }

    fetch(endpoint)
      .then((response: Response) => response.text())
      .then((html: string) => {
        const target: HTMLElement | null = document.getElementById('content-inner');

        if (!target) {
          console.error('Target element #content-inner not found');
          return;
        }

        // Simply replace the content
        target.innerHTML = html;

        // Initialize Alpine on the new content
        window.Alpine.initTree(target);
      })
      .catch((error: Error) => {
        console.error('AJAX request failed:', error);
      });
  }
});

// Register the Alpine component when Alpine initializes
document.addEventListener('alpine:init', () => {
  window.Alpine.data('ajaxButton', createAjaxButtonComponent);
});
