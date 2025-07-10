// API Configuration
export const BASE_URL = import.meta.env.VITE_API_BASE_URL;
export const API_KEY = import.meta.env.VITE_API_KEY;

// API Endpoints
export const ENDPOINTS = {
  HEALTH: "/api/health",
  CRAWL: "/api/crawl",
  CRAWL_STATS: "/api/crawl/stats",
  CRAWL_RERUN: "/api/crawl/rerun",
} as const;
