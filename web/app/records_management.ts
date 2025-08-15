import Alpine from "alpinejs";

// Define the shape of the API's JSON response
interface RecordItem {
  id: string;
  name: string;
  // add other record fields as needed
}

interface RecordsApiResponse {
  data: RecordItem[];
  pagination: {
    hasMore: boolean;
    nextCursor: string | null;
  };
}

// Register Alpine component
Alpine.data('recordsListComponent', () => ({
  records: [] as RecordItem[],
  rawResults: "" as string,
  loading: false as boolean,
  error: null as string | null,
  hasMore: true as boolean,
  nextCursor: null as string | null,

  async loadRecords(append: boolean = false) {
    this.loading = true;
    const url = new URL("/app/api/records", window.location.origin);
    url.searchParams.set("limit", "20");
    const request = new Request(url);

    try {
      const response = await fetch(request);
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
      console.error(err);
      this.error = "Failed to load records";
    } finally {
      this.loading = false;
    }
  }
}));
