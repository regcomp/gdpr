interface RecordItem {
  id: string;
  name: string;
  // TODO: finalize the shape of this
}

interface RecordsApiResponse {
  data: RecordItem[];
  pagination: {
    hasMore: boolean;
    nextCursor: string | null;
  };
}

interface RecordsListComponentData {
  records: RecordItem[];
  rawResults: string;
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  nextCursor: string | null;
  loadRecords(append?: boolean): Promise<void>;
}

const createRecordsListComponent = (): RecordsListComponentData => ({
  records: [] as RecordItem[],
  rawResults: "" as string,
  loading: false as boolean,
  error: null as string | null,
  hasMore: true as boolean,
  nextCursor: null as string | null,

  async loadRecords(append: boolean = false): Promise<void> {
    this.loading = true;
    this.error = null; // Clear previous errors

    const url = new URL("/app/api/records", window.location.origin);
    url.searchParams.set("limit", "20");

    // Add cursor for pagination if we have one
    if (append && this.nextCursor) {
      url.searchParams.set("cursor", this.nextCursor);
    }

    try {
      const response = await fetch(url.toString());

      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }

      const json: RecordsApiResponse = await response.json();

      if (append) {
        this.records.push(...json.data);
      } else {
        this.records = json.data;
      }

      this.hasMore = json.pagination.hasMore;
      this.nextCursor = json.pagination.nextCursor;
      this.rawResults = JSON.stringify(json.data, null, 2);

    } catch (err) {
      console.error('Failed to load records:', err);
      this.error = err instanceof Error ? err.message : "Failed to load records";
    } finally {
      this.loading = false;
    }
  }
});

document.addEventListener('alpine:init', () => {
  window.Alpine.data('recordsListComponent', createRecordsListComponent);
});
