import axiosInstance from "./axiosInstance";
import { ENDPOINTS } from "./config";
import type { CrawlResult, CrawlFilters, CrawlStats } from "../types";

// Response types that match the backend
export interface PaginatedResponse<T> {
  results: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface CrawlRequestResponse {
  id: string;
  url: string;
  status: string;
  message: string;
}

export interface StatsResponse {
  database: CrawlStats;
  queue: {
    queue_length: number;
    active_tasks: number;
    workers: number;
    running: boolean;
  };
  timestamp: string;
}

// Health check
const healthCheck = async (): Promise<{
  status: string;
  timestamp: string;
}> => {
  const response = await axiosInstance.get(ENDPOINTS.HEALTH);
  return response.data as { status: string; timestamp: string };
};

// Submit URL for crawling
const submitCrawlRequest = async (
  url: string
): Promise<CrawlRequestResponse> => {
  const response = await axiosInstance.post(ENDPOINTS.CRAWL, { url });
  return response.data as CrawlRequestResponse;
};

// Get crawl results with pagination and filtering
const getCrawlResults = async (
  filters?: Partial<CrawlFilters>
): Promise<PaginatedResponse<CrawlResult>> => {
  const params: Record<string, any> = {};

  if (filters) {
    if (filters.page) params.page = filters.page;
    if (filters.pageSize) params.pageSize = filters.pageSize;
    if (filters.search) params.search = filters.search;
    if (filters.status) params.status = filters.status;
    if (filters.sortBy) params.sortBy = filters.sortBy;
    if (filters.sortDir) params.sortDir = filters.sortDir;
  }

  const response = await axiosInstance.get(ENDPOINTS.CRAWL, { params });
  const data = response.data as PaginatedResponse<CrawlResult>;
  // Transform date strings to Date objects
  if (data.results !== null) {
    data.results = data.results.map((result) => ({
      ...result,
      createdAt: new Date(result.createdAt),
      updatedAt: new Date(result.updatedAt),
    }));
    return data;
  }
  return {
    results: [],
    total: 0,
    page: 1,
    pageSize: 10,
    totalPages: 1,
  };
};

// Get single crawl result
const getCrawlResult = async (id: string): Promise<CrawlResult> => {
  const response = await axiosInstance.get(`${ENDPOINTS.CRAWL}/${id}`);
  const result = response.data as CrawlResult;

  // Transform date strings to Date objects
  return {
    ...result,
    createdAt: new Date(result.createdAt),
    updatedAt: new Date(result.updatedAt),
  };
};

// Get crawl status
const getCrawlStatus = async (
  id: string
): Promise<{
  id: string;
  status: string;
  url: string;
  created_at: string;
  updated_at: string;
  error_message?: string;
}> => {
  const response = await axiosInstance.get(`${ENDPOINTS.CRAWL}/${id}/status`);
  return response.data as {
    id: string;
    status: string;
    url: string;
    created_at: string;
    updated_at: string;
    error_message?: string;
  };
};

// Delete crawl results
const deleteCrawlResults = async (
  ids: string[]
): Promise<{
  message: string;
  deleted_count: number;
}> => {
  const response = await axiosInstance.request({
    method: "DELETE",
    url: ENDPOINTS.CRAWL,
    data: { ids },
  });
  return response.data as unknown as { message: string; deleted_count: number };
};

// Re-run crawl analysis
const rerunCrawlResults = async (
  ids: string[]
): Promise<{
  message: string;
  success_count: number;
  total_requested: number;
  errors?: string[];
}> => {
  const response = await axiosInstance.post(ENDPOINTS.CRAWL_RERUN, { ids });
  return response.data as {
    message: string;
    success_count: number;
    total_requested: number;
    errors?: string[];
  };
};

// Get crawl statistics
const getCrawlStats = async (): Promise<StatsResponse> => {
  const response = await axiosInstance.get(ENDPOINTS.CRAWL_STATS);
  return response.data as StatsResponse;
};

export const crawlApis = {
  healthCheck,
  submitCrawlRequest,
  getCrawlResults,
  getCrawlResult,
  getCrawlStatus,
  deleteCrawlResults,
  rerunCrawlResults,
  getCrawlStats,
};
