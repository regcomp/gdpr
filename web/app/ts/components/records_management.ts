import { Configuration, RecordsApi, Record, GetRecordsWithFilters200Response } from "../../../generated/openapi";

interface RecordsListComponentData {
  records: Record[];
  rawResults: string;
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  nextCursor: Date;
  loadRecords(append?: boolean): Promise<void>;
  recordsApi: RecordsApi
}

const createRecordsListComponent = (): RecordsListComponentData => {
  const recordsApi = new RecordsApi(new Configuration({
    basePath: "https://localhost:8080/app/api" // TODO: Softcode this string
  }))

  return {
    records: [] as Record[],
    rawResults: "" as string, // currently just for dosplay purposes
    loading: false as boolean,
    error: null as string | null,
    hasMore: true as boolean,
    nextCursor: new Date(0),
    recordsApi: recordsApi,

    async loadRecords(): Promise<void> {
      this.loading = true;
      this.error = null; // Clear previous errors

      if (!this.hasMore) {
        return
      }

      let response: GetRecordsWithFilters200Response;
      try {
        response = await this.recordsApi.getRecordsWithFilters({
          after: this.nextCursor,
          limit: 20,
        })
      } catch (error) {
        // TODO:
        return
      }

      this.records.push(...response.data)
      this.rawResults = JSON.stringify(this.records, null, 2);

      this.hasMore = response.pagination.hasMore
      this.nextCursor = response.pagination.nextCursor
      this.loading = false
    }
  }
};

document.addEventListener('alpine:init', () => {
  window.Alpine.data('recordsListComponent', createRecordsListComponent);
});
