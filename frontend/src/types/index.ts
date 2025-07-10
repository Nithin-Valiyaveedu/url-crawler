export type ViewMode = "dashboard" | "detail";

export interface CrawlResult {
  id: string;
  url: string;
  title: string;
  htmlVersion: string;
  internalLinksCount: number;
  externalLinksCount: number;
  inaccessibleLinksCount: number;
  hasLoginForm: boolean;
  headingCounts: {
    h1: number;
    h2: number;
    h3: number;
    h4: number;
    h5: number;
    h6: number;
  };
  brokenLinks: BrokenLink[];
  externalLinks?: string[];
  status: CrawlStatus;
  createdAt: Date;
  updatedAt: Date;
}

export interface BrokenLink {
  url: string;
  statusCode: number;
  statusText: string;
}

export type CrawlStatus = "queued" | "running" | "completed" | "error";

export interface CrawlRequest {
  url: string;
}

export interface TableFilter {
  column: string;
  value: string;
}

export interface SortConfig {
  column: keyof CrawlResult;
  direction: "asc" | "desc";
}

export interface PaginationConfig {
  page: number;
  pageSize: number;
  total: number;
}

// API-related types
export interface CrawlFilters {
  status?: CrawlStatus;
  search?: string;
  page: number;
  pageSize: number;
  sortBy?: string;
  sortDir?: "asc" | "desc";
}

export interface CrawlStats {
  total: number;
  queued: number;
  running: number;
  completed: number;
  error: number;
}
