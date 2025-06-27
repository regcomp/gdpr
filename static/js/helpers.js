document.addEventListener('alpine:init', () => {
  Alpine.store('api', {
    async fetch(url, options = {}) {
      const makeRequest = async () => {
        const headers = { ...options.headers };

        headers['Content-Type'] = options.contentType || 'application/json'

        const config = {
          ...options,
          credentials: 'include',
          headers
        };

        const response = await fetch(url, config);
        return response;
      };

      try {
        let response = await makeRequest();

        if (response.status === 401 && response.headers.has('X-Token-Retry')) {
          const refreshed = await this.refreshToken();
          if (refreshed) {
            response = await makeRequest();
          }
        }

        return response;
      } catch (error) {
        console.error('API request failed:', error);
        throw error;
      }
    },

    async refreshToken() {
      try {
        const response = await fetch('/auth/refresh', {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' }
        });

        return response.ok;
      } catch (error) {
        console.error('Token refresh failed:', error);
        return false;
      }
    }
  });
});
