// Import for local use
import { crawlApis } from "./crawl";

export { crawlApis } from "./crawl";

export type {
  PaginatedResponse,
  CrawlRequestResponse,
  StatsResponse,
} from "./crawl";

// Main API service object for backward compatibility
export const apiService = {
  ...crawlApis,
};

// Error handling utility
export const handleApiError = (error: unknown): string => {
  if (error && typeof error === "object" && "response" in error) {
    const axiosError = error as any;
    if (axiosError.response?.data?.error) {
      return axiosError.response.data.error;
    }
    if (axiosError.response?.data?.message) {
      return axiosError.response.data.message;
    }
    if (axiosError.response?.statusText) {
      return axiosError.response.statusText;
    }
    if (axiosError.message) {
      return axiosError.message;
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "An unexpected error occurred";
};
