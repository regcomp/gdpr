class UserDataService {
  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl;
    this.abortController = null;
  }

  /**
   * Load users with cursor-based pagination
   * @param {string|null} cursor - Timestamp cursor for pagination
   * @param {number} limit - Number of records to fetch
   * @returns {Promise<{data: Array, pagination: Object}>}
   */
  async loadUsers(cursor = null, limit = 20) {
    // Cancel any existing request
    if (this.abortController) {
      this.abortController.abort();
    }

    this.abortController = new AbortController();

    try {
      const params = new URLSearchParams({ limit: limit.toString() });
      if (cursor) {
        params.append('after', cursor);
      }

      const response = await fetch(`${this.baseUrl}/users?${params}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        signal: this.abortController.signal
      });

      if (!response.ok) {
        throw new Error(`Failed to load users: ${response.status} ${response.statusText}`);
      }

      const data = await response.json();

      // Validate response structure
      if (!data.data || !Array.isArray(data.data)) {
        throw new Error('Invalid response format: missing data array');
      }

      if (!data.pagination) {
        throw new Error('Invalid response format: missing pagination info');
      }

      return {
        data: data.data,
        pagination: {
          hasMore: data.pagination.hasMore || false,
          nextCursor: data.pagination.nextCursor || null,
          total: data.pagination.total || null
        }
      };

    } catch (error) {
      if (error.name === 'AbortError') {
        throw new Error('Request was cancelled');
      }

      // Network or other errors
      if (!navigator.onLine) {
        throw new Error('No internet connection');
      }

      throw error;

    } finally {
      this.abortController = null;
    }
  }

  /**
   * Delete a user by ID
   * @param {string} userId - UUID of user to delete
   * @returns {Promise<void>}
   */
  async deleteUser(userId) {
    if (!userId) {
      throw new Error('User ID is required');
    }

    try {
      const response = await fetch(`${this.baseUrl}/users/${userId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('User not found');
        }
        if (response.status === 403) {
          throw new Error('Not authorized to delete this user');
        }
        throw new Error(`Failed to delete user: ${response.status} ${response.statusText}`);
      }

      // Some APIs return 204 No Content, others return 200 with data
      if (response.status !== 204) {
        const result = await response.json();
        return result;
      }

    } catch (error) {
      // Network or other errors
      if (!navigator.onLine) {
        throw new Error('No internet connection');
      }

      throw error;
    }
  }

  /**
   * Cancel any ongoing requests
   */
  cancelRequests() {
    if (this.abortController) {
      this.abortController.abort();
      this.abortController = null;
    }
  }

  /**
   * Search users (optional enhancement)
   * @param {string} query - Search query
   * @param {string|null} cursor - Pagination cursor
   * @param {number} limit - Number of results
   * @returns {Promise<{data: Array, pagination: Object}>}
   */
  async searchUsers(query, cursor = null, limit = 20) {
    if (!query || query.trim().length === 0) {
      return this.loadUsers(cursor, limit);
    }

    const params = new URLSearchParams({
      q: query.trim(),
      limit: limit.toString()
    });

    if (cursor) {
      params.append('after', cursor);
    }

    try {
      const response = await fetch(`${this.baseUrl}/users/search?${params}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (!response.ok) {
        throw new Error(`Search failed: ${response.status} ${response.statusText}`);
      }

      const data = await response.json();

      return {
        data: data.data || [],
        pagination: {
          hasMore: data.pagination?.hasMore || false,
          nextCursor: data.pagination?.nextCursor || null,
          total: data.pagination?.total || null
        }
      };

    } catch (error) {
      if (!navigator.onLine) {
        throw new Error('No internet connection');
      }
      throw error;
    }
  }
}

// Create singleton instance
const userDataService = new UserDataService();

// Export for use in Alpine components
window.userDataService = userDataService;
