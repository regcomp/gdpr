document.addEventListener('alpine:init', () => {
  Alpine.data('userListComponent', () => ({
    // State
    users: [],
    loading: false,
    error: null,
    hasMore: true,
    nextCursor: null,

    // Computed properties for template
    get loadingMore() {
      return this.loading && this.users.length > 0;
    },

    get endOfData() {
      return !this.hasMore && this.users.length > 0;
    },

    get isEmpty() {
      return this.users.length === 0 && !this.loading;
    },

    // Methods
    async loadUsers(append = false) {
      this.loading = true;
      this.error = null;

      try {
        const response = await userDataService.loadUsers(
          append ? this.nextCursor : null
        );

        if (append) {
          this.users.push(...response.data);
        } else {
          this.users = response.data;
        }

        this.hasMore = response.pagination.hasMore;
        this.nextCursor = response.pagination.nextCursor;

      } catch (err) {
        this.error = 'Failed to load users';
      } finally {
        this.loading = false;
      }
    },

    refreshUsers() {
      this.loadUsers(false);
    },

    checkScroll(event) {
      const element = event.target;
      const threshold = 200;

      if (element.scrollTop + element.clientHeight >= element.scrollHeight - threshold) {
        if (this.hasMore && !this.loading) {
          this.loadUsers(true);
        }
      }
    },

    deleteUser(event) {
      const userId = event.target.dataset.userId;
      this.performDelete(userId);
    },

    async performDelete(userId) {
      const userIndex = this.users.findIndex(u => u.id === userId);
      const removedUser = this.users[userIndex];

      this.users.splice(userIndex, 1);

      try {
        await userDataService.deleteUser(userId);
      } catch (err) {
        this.users.splice(userIndex, 0, removedUser);
        this.error = 'Failed to delete user';
      }
    }
  }));
});
