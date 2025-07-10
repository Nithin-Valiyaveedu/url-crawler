import axios from "axios";
import { BASE_URL, API_KEY } from "./config";

// Creating an Axios instance
const axiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
    Authorization: `Bearer ${API_KEY}`,
  },
});

// Request interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    // Log request for debugging
    console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error("[API] Request error:", error);
    return Promise.reject(error);
  }
);

// Response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    console.log(`[API] ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error(
      "[API] Response error:",
      error.response?.status,
      error.message
    );

    if (error.response) {
      const { status, data } = error.response;
      const errorMessage = data?.error || data?.message || "An error occurred";

      switch (status) {
        case 400:
          console.error("Bad Request:", errorMessage);
          break;

        case 401:
          console.error("Unauthorized:", errorMessage);
          break;

        case 403:
          console.error("Forbidden:", errorMessage);
          break;

        case 404:
          console.error("Not Found:", errorMessage);
          break;

        case 429:
          console.error("Rate Limited:", errorMessage);
          break;

        case 500:
          console.error("Internal Server Error:", errorMessage);
          break;

        default:
          console.error(`HTTP ${status}:`, errorMessage);
          break;
      }
    } else if (error.request) {
      console.error("Network error - no response received");
    } else {
      console.error("Request setup error:", error.message);
    }

    return Promise.reject(error);
  }
);

export default axiosInstance;
