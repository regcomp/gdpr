Alpine.data('recordsListComponent', () => ({
  records: [],
  rawResults: "",
  loading: false,
  error: null,
  hasMore: true,
  nextCursor: null,

  async loadRecords(append = false) {
    const url = new URL("/app/api/records", window.location.origin)
    url.searchParams.set('limit', '20')
    const request = new Request(url);

    try {
      const response = await fetch(request);

      if (append) {
        this.records.push(...response.data);
      } else {
        this.records = response.data;
      }

      this.hasMore = response.pagination.hasMore;
      this.nextCursor = response.pagination.nextCursor;

      this.rawResults = JSON.stringify(response.data, null, 2)

    } catch (err) {
      this.error = 'Failed to load users';
    } finally {
      this.loading = false;
    }
  },

}));
