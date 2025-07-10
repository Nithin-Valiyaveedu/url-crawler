import { useState } from "react";
import { Plus, Loader } from "lucide-react";

interface UrlInputProps {
  onAddUrl: (url: string) => Promise<void>;
  isLoading?: boolean;
}

export const UrlInput: React.FC<UrlInputProps> = ({
  onAddUrl,
  isLoading = false,
}) => {
  const [url, setUrl] = useState("");
  const [error, setError] = useState("");

  const isValidUrl = (string: string) => {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!url.trim()) {
      setError("Please enter a URL");
      return;
    }

    if (!isValidUrl(url)) {
      setError("Please enter a valid URL");
      return;
    }

    try {
      await onAddUrl(url);
      setUrl("");
    } catch (err) {
      setError("Failed to add URL. Please try again.");
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-6">
      <h2 className="text-lg sm:text-xl font-semibold text-gray-800 mb-4">
        Add URL to Crawl
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="url"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Website URL
          </label>
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-2">
            <input
              type="url"
              id="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://example.com"
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm sm:text-base"
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !url.trim()}
              className="w-full sm:w-auto px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 text-sm sm:text-base font-medium"
            >
              {isLoading ? (
                <Loader className="w-4 h-4 animate-spin" />
              ) : (
                <Plus className="w-4 h-4" />
              )}
              {isLoading ? "Adding..." : "Add URL"}
            </button>
          </div>
          {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
        </div>
      </form>

      <div className="mt-4 text-xs sm:text-sm text-gray-500">
        <p>Enter a valid URL to start crawling and analyzing the website.</p>
      </div>
    </div>
  );
};
