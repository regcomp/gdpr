import {
  Configuration,
  RecordsApi,
  Record as ApiRecord,
  GetRecordsWithFilters200Response,
} from "../../../generated/openapi";

interface DisplayRecord extends ApiRecord {
  formattedRequestedAt: string;
}

interface RecordsListComponentData {
  records: DisplayRecord[];
  loading: boolean;
  error: string;
  showError: boolean;
  hasMore: boolean;
  isEmpty: boolean;
  lastScrollTime: number;
  nextCursor: Date;
  recordsApi: RecordsApi

  init(): Promise<void>;
  loadRecords(): Promise<void>;
  handleScroll(event: Event): void;
  recomputeDerived(): void;
}

const createRecordsListComponent = (): RecordsListComponentData => {
  const recordsApi = new RecordsApi(new Configuration({
    basePath: "https://localhost:8080/app/api" // TODO: Softcode this string
  }))

  return {
    records: [] as DisplayRecord[],
    loading: false as boolean,
    error: "" as string,
    showError: false as boolean,
    hasMore: true as boolean,
    isEmpty: true as boolean,
    lastScrollTime: 0 as number,
    nextCursor: new Date(0),
    recordsApi: recordsApi,

    async init(): Promise<void> {
      console.log('records manager loading...');
      this.loadRecords();
      console.log('...records manager loaded!');
    },

    async loadRecords(): Promise<void> {
      console.log('loading records...');
      this.loading = true;
      this.error = "";

      if (!this.hasMore) {
        return;
      }

      let response: GetRecordsWithFilters200Response;
      try {
        response = await this.recordsApi.getRecordsWithFilters({
          after: this.nextCursor,
          limit: 20,
        });
        const incoming = response.data.map<DisplayRecord>((resp) => {
          return {
            ...resp,
            formattedRequestedAt: resp.requestedAt.toLocaleString()
          };
        });
        this.records.push(...incoming);
        this.hasMore = response.pagination.hasMore;
        this.nextCursor = response.pagination.nextCursor;
      } catch (err) {
        this.error = getErrorMessage(err);
      } finally {
        this.recomputeDerived();
        this.loading = false;
      }
    },

    handleScroll(event: Event): void {
      const now = Date.now();

      if (now - this.lastScrollTime < 250) return;
      this.lastScrollTime = now;

      const el = event.currentTarget as HTMLElement;
      const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 100;
      if (nearBottom && !this.loading && this.hasMore) {
        console.log('handling scroll...')
        this.loadRecords();
      }
    },

    recomputeDerived(): void {
      const notLoading = !this.loading;
      const noLength = this.records.length === 0;
      const noError = !this.error;
      this.isEmpty = notLoading && noLength && noError;

      this.showError = this.error != "" && this.hasMore;
    }
  }
};

document.addEventListener('alpine:init', () => {
  window.Alpine.data('recordsListComponent', createRecordsListComponent);
});

function getErrorMessage(e: unknown): string {
  if (e instanceof Error) return e.message;
  if (typeof e === "string") return e;
  return JSON.stringify(e);
}

