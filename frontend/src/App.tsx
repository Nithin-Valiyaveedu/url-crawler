import { useState, useEffect } from "react";

// Components
import {
  UrlInput,
  CrawlResultsTable,
  CrawlResultDetail,
  Loader,
  ErrorNotice,
  UrlStatusbar,
} from "./components";

// Types
import type { CrawlResult, ViewMode, CrawlFilters } from "./types";

// API
import { apiService, handleApiError } from "./apis/index";

function App() {
  const [currentView, setCurrentView] = useState<ViewMode>("dashboard");
  const [selectedResultId, setSelectedResultId] = useState<string | null>(null);
  const [crawlResults, setCrawlResults] = useState<CrawlResult[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isAddingUrl, setIsAddingUrl] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0,
  });

  // Loading initial data
  useEffect(() => {
    loadCrawlResults();
  }, []);

  const loadCrawlResults = async (filters?: Partial<CrawlFilters>) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiService.getCrawlResults({
        page: 1,
        pageSize: 20,
        ...filters,
      });
      console.log("Crawl results loaded", response);
      setCrawlResults(response.results);
      setPagination({
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
      });
    } catch (error) {
      console.log("Error loading crawl results", error);
      const errorMessage = handleApiError(error);
      setError(errorMessage);
      console.error("Failed to load crawl results:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Polling every 5 seconds for status updates
  useEffect(() => {
    const interval = setInterval(() => {
      // Only poll if there are running or queued items
      const hasActiveItems = crawlResults.some(
        (result) => result.status === "running" || result.status === "queued"
      );

      if (hasActiveItems) {
        loadCrawlResults();
      }
    }, 5000);

    return () => clearInterval(interval);
  }, [crawlResults]);

  const handleAddUrl = async (url: string): Promise<void> => {
    setIsAddingUrl(true);
    setError(null);
    try {
      await apiService.submitCrawlRequest(url);

      await loadCrawlResults();
    } catch (error) {
      const errorMessage = handleApiError(error);
      setError(errorMessage);
      console.error("Failed to add URL:", error);
      throw error;
    } finally {
      setIsAddingUrl(false);
    }
  };

  const handleViewDetails = async (id: string) => {
    try {
      // Fetch the latest data for the selected result
      const result = await apiService.getCrawlResult(id);

      setCrawlResults((prev) => prev.map((r) => (r.id === id ? result : r)));

      setSelectedResultId(id);
      setCurrentView("detail");
    } catch (error) {
      const errorMessage = handleApiError(error);
      setError(errorMessage);
      console.error("Failed to load crawl result details:", error);
    }
  };

  const handleBackToDashboard = () => {
    setCurrentView("dashboard");
    setSelectedResultId(null);
  };

  const handleDeleteResults = async (ids: string[]) => {
    try {
      await apiService.deleteCrawlResults(ids);

      setCrawlResults((prev) =>
        prev.filter((result) => !ids.includes(result.id))
      );

      // Update pagination
      setPagination((prev) => ({
        ...prev,
        total: prev.total - ids.length,
      }));
    } catch (error) {
      const errorMessage = handleApiError(error);
      setError(errorMessage);
      console.error("Failed to delete results:", error);
      throw error;
    }
  };

  const handleRerunResults = async (ids: string[]) => {
    try {
      await apiService.rerunCrawlResults(ids);

      // Update local state to show queued status
      setCrawlResults((prev) =>
        prev.map((result) =>
          ids.includes(result.id)
            ? { ...result, status: "queued" as const, updatedAt: new Date() }
            : result
        )
      );
    } catch (error) {
      const errorMessage = handleApiError(error);
      setError(errorMessage);
      console.error("Failed to rerun results:", error);
      throw error;
    }
  };

  const selectedResult = selectedResultId
    ? crawlResults.find((r) => r.id === selectedResultId)
    : null;

  if (isLoading && crawlResults.length === 0) {
    return <Loader />;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col lg:flex-row lg:justify-between lg:items-center py-4 lg:h-16">
            <div className="flex items-center mb-4 lg:mb-0">
              <h1 className="text-xl lg:text-2xl font-bold text-gray-900">
                URL Data Dashboard
              </h1>
              {currentView === "detail" && (
                <span className="ml-4 text-sm text-gray-500 hidden sm:inline">
                  › Detailed Analysis
                </span>
              )}
            </div>
            <UrlStatusbar crawlResults={crawlResults} />
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {error && <ErrorNotice error={error} setError={setError} />}

        {currentView === "dashboard" ? (
          <>
            <UrlInput onAddUrl={handleAddUrl} isLoading={isAddingUrl} />
            <CrawlResultsTable
              results={crawlResults}
              onViewDetails={handleViewDetails}
              onDeleteResults={handleDeleteResults}
              onRerunResults={handleRerunResults}
              isLoading={isLoading}
            />
          </>
        ) : (
          selectedResult && (
            <CrawlResultDetail
              result={selectedResult}
              onBack={handleBackToDashboard}
            />
          )
        )}
      </main>
    </div>
  );
}

export default App;
