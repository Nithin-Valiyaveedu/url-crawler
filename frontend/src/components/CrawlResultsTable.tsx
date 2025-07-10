import { useState, useMemo } from "react";
import {
  ChevronUp,
  ChevronDown,
  Search,
  Trash2,
  RefreshCw,
  Eye,
  CheckCircle,
  Clock,
  AlertCircle,
  Loader,
} from "lucide-react";
import type { CrawlResult, SortConfig } from "../types";

interface CrawlResultsTableProps {
  results: CrawlResult[];
  onViewDetails: (id: string) => void;
  onDeleteResults: (ids: string[]) => Promise<void>;
  onRerunResults: (ids: string[]) => Promise<void>;
  isLoading?: boolean;
}

const ITEMS_PER_PAGE = 10;

export const CrawlResultsTable: React.FC<CrawlResultsTableProps> = ({
  results,
  onViewDetails,
  onDeleteResults,
  onRerunResults,
  isLoading = false,
}) => {
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [sortConfig, setSortConfig] = useState<SortConfig>({
    column: "updatedAt",
    direction: "desc",
  });
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [currentPage, setCurrentPage] = useState(1);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isRerunning, setIsRerunning] = useState(false);

  const getStatusIcon = (status: CrawlResult["status"]) => {
    switch (status) {
      case "completed":
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case "running":
        return <Loader className="w-4 h-4 text-blue-500 animate-spin" />;
      case "queued":
        return <Clock className="w-4 h-4 text-yellow-500" />;
      case "error":
        return <AlertCircle className="w-4 h-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusBadge = (status: CrawlResult["status"]) => {
    const baseClasses =
      "inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium";
    switch (status) {
      case "completed":
        return `${baseClasses} bg-green-100 text-green-800`;
      case "running":
        return `${baseClasses} bg-blue-100 text-blue-800`;
      case "queued":
        return `${baseClasses} bg-yellow-100 text-yellow-800`;
      case "error":
        return `${baseClasses} bg-red-100 text-red-800`;
      default:
        return baseClasses;
    }
  };

  const filteredAndSortedResults = useMemo(() => {
    let filtered = results.filter((result) => {
      const matchesSearch =
        result.url.toLowerCase().includes(searchTerm.toLowerCase()) ||
        result.title.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesStatus =
        statusFilter === "all" || result.status === statusFilter;
      return matchesSearch && matchesStatus;
    });

    filtered.sort((a, b) => {
      const aValue = a[sortConfig.column];
      const bValue = b[sortConfig.column];

      // Handle null/undefined values
      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return sortConfig.direction === "asc" ? 1 : -1;
      if (bValue == null) return sortConfig.direction === "asc" ? -1 : 1;

      if (aValue < bValue) return sortConfig.direction === "asc" ? -1 : 1;
      if (aValue > bValue) return sortConfig.direction === "asc" ? 1 : -1;
      return 0;
    });

    return filtered;
  }, [results, searchTerm, statusFilter, sortConfig]);

  const paginatedResults = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredAndSortedResults.slice(
      startIndex,
      startIndex + ITEMS_PER_PAGE
    );
  }, [filteredAndSortedResults, currentPage]);

  const totalPages = Math.ceil(
    filteredAndSortedResults.length / ITEMS_PER_PAGE
  );

  const handleSort = (column: keyof CrawlResult) => {
    setSortConfig((prev) => ({
      column,
      direction:
        prev.column === column && prev.direction === "asc" ? "desc" : "asc",
    }));
  };

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedIds(new Set(paginatedResults.map((r) => r.id)));
    } else {
      setSelectedIds(new Set());
    }
  };

  const handleSelectItem = (id: string, checked: boolean) => {
    const newSelected = new Set(selectedIds);
    if (checked) {
      newSelected.add(id);
    } else {
      newSelected.delete(id);
    }
    setSelectedIds(newSelected);
  };

  const handleBulkDelete = async () => {
    if (selectedIds.size === 0) return;
    setIsDeleting(true);
    try {
      await onDeleteResults(Array.from(selectedIds));
      setSelectedIds(new Set());
    } finally {
      setIsDeleting(false);
    }
  };

  const handleBulkRerun = async () => {
    if (selectedIds.size === 0) return;
    setIsRerunning(true);
    try {
      await onRerunResults(Array.from(selectedIds));
      setSelectedIds(new Set());
    } finally {
      setIsRerunning(false);
    }
  };

  const SortButton: React.FC<{
    column: keyof CrawlResult;
    children: React.ReactNode;
  }> = ({ column, children }) => (
    <button
      onClick={() => handleSort(column)}
      className="flex items-center gap-1 hover:text-gray-900 text-left w-full"
    >
      {children}
      {sortConfig.column === column &&
        (sortConfig.direction === "asc" ? (
          <ChevronUp className="w-4 h-4" />
        ) : (
          <ChevronDown className="w-4 h-4" />
        ))}
    </button>
  );

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <div className="p-4 sm:p-6 border-b border-gray-200">
        <div className="flex flex-col gap-4">
          <h2 className="text-lg sm:text-xl font-semibold text-gray-800">
            Crawl Results
          </h2>

          <div className="flex flex-col sm:flex-row gap-3">
            <div className="relative flex-1">
              <Search className="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Search URLs or titles..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <div className="sm:w-48">
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="all">All Statuses</option>
                <option value="completed">Completed</option>
                <option value="running">Running</option>
                <option value="queued">Queued</option>
                <option value="error">Error</option>
              </select>
            </div>
          </div>
        </div>

        {/* Bulk actions */}
        {selectedIds.size > 0 && (
          <div className="mt-4 flex flex-col sm:flex-row sm:items-center gap-3">
            <span className="text-sm text-gray-600">
              {selectedIds.size} item{selectedIds.size > 1 ? "s" : ""} selected
            </span>
            <div className="flex gap-2">
              <button
                onClick={handleBulkRerun}
                disabled={isRerunning || isLoading}
                className="flex items-center gap-1 px-3 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 flex-1 sm:flex-initial justify-center"
              >
                {isRerunning ? (
                  <Loader className="w-3 h-3 animate-spin" />
                ) : (
                  <RefreshCw className="w-3 h-3" />
                )}
                Re-run
              </button>
              <button
                onClick={handleBulkDelete}
                disabled={isDeleting || isLoading}
                className="flex items-center gap-1 px-3 py-2 text-sm bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 flex-1 sm:flex-initial justify-center"
              >
                {isDeleting ? (
                  <Loader className="w-3 h-3 animate-spin" />
                ) : (
                  <Trash2 className="w-3 h-3" />
                )}
                Delete
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="px-4 sm:px-6 py-3 bg-gray-50 border-b border-gray-200">
        <div className="text-sm text-gray-700">
          Showing {paginatedResults.length} of {filteredAndSortedResults.length}{" "}
          results
        </div>
      </div>

      <div className="md:hidden">
        {paginatedResults.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500">No results found</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {paginatedResults.map((result) => (
              <div
                key={result.id}
                className="p-4 cursor-pointer hover:bg-gray-50 transition-colors"
                onClick={() => onViewDetails(result.id)}
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      checked={selectedIds.has(result.id)}
                      onChange={(e) =>
                        handleSelectItem(result.id, e.target.checked)
                      }
                      onClick={(e) => e.stopPropagation()}
                      className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span className={getStatusBadge(result.status)}>
                      {getStatusIcon(result.status)}
                      {result.status}
                    </span>
                  </div>
                </div>

                <div className="space-y-2">
                  <div>
                    <a
                      href={result.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      onClick={(e) => e.stopPropagation()}
                      className="text-blue-600 hover:text-blue-800 text-sm font-medium break-all"
                    >
                      {result.url}
                    </a>
                  </div>

                  {result.title && (
                    <div className="text-sm text-gray-900 font-medium">
                      {result.title}
                    </div>
                  )}

                  <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-gray-600">
                    <div>Internal: {result.internalLinksCount}</div>
                    <div>External: {result.externalLinksCount}</div>
                    <div>
                      Broken:{" "}
                      <span
                        className={
                          result.inaccessibleLinksCount > 0
                            ? "text-red-600"
                            : ""
                        }
                      >
                        {result.inaccessibleLinksCount}
                      </span>
                    </div>
                    <div>Login: {result.hasLoginForm ? "Yes" : "No"}</div>
                  </div>

                  <div className="text-xs text-gray-500">
                    Updated: {new Date(result.updatedAt).toLocaleDateString()}
                  </div>
                </div>

                <div className="mt-3 pt-3 border-t border-gray-200">
                  <div className="w-full flex items-center justify-center gap-1 px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded">
                    <Eye className="w-3 h-3" />
                    View Details
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="hidden md:block overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left">
                <input
                  type="checkbox"
                  checked={
                    paginatedResults.length > 0 &&
                    selectedIds.size === paginatedResults.length
                  }
                  onChange={(e) => handleSelectAll(e.target.checked)}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="status">Status</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="url">URL</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="title">Title</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="htmlVersion">HTML Version</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="internalLinksCount">
                  Internal Links
                </SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="externalLinksCount">
                  External Links
                </SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="inaccessibleLinksCount">
                  Broken Links
                </SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="hasLoginForm">Login Form</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <SortButton column="updatedAt">Last Updated</SortButton>
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {paginatedResults.map((result) => (
              <tr
                key={result.id}
                className="hover:bg-gray-50 cursor-pointer"
                onClick={() => onViewDetails(result.id)}
              >
                <td className="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    checked={selectedIds.has(result.id)}
                    onChange={(e) =>
                      handleSelectItem(result.id, e.target.checked)
                    }
                    onClick={(e) => e.stopPropagation()}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={getStatusBadge(result.status)}>
                    {getStatusIcon(result.status)}
                    {result.status}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="max-w-xs truncate">
                    <a
                      href={result.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      onClick={(e) => e.stopPropagation()}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      {result.url}
                    </a>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="max-w-xs truncate" title={result.title}>
                    {result.title || "-"}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {result.htmlVersion || "-"}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {result.internalLinksCount}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {result.externalLinksCount}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <span
                    className={
                      result.inaccessibleLinksCount > 0 ? "text-red-600" : ""
                    }
                  >
                    {result.inaccessibleLinksCount}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {result.hasLoginForm ? "Yes" : "No"}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(result.updatedAt).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div className="text-blue-600 flex items-center gap-1">
                    <Eye className="w-4 h-4" />
                    View
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {paginatedResults.length === 0 && (
          <div className="text-center py-8">
            <p className="text-gray-500">No results found</p>
          </div>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="px-4 sm:px-6 py-3 bg-gray-50 border-t border-gray-200">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-3">
            <div className="text-sm text-gray-700">
              Page {currentPage} of {totalPages}
            </div>
            <div className="flex gap-1">
              <button
                onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1}
                className="px-3 py-1 text-sm border border-gray-300 rounded disabled:opacity-50 hover:bg-gray-50"
              >
                Previous
              </button>
              {[...Array(totalPages)].map((_, i) => {
                const page = i + 1;
                const isCurrentPage = page === currentPage;
                const shouldShow =
                  page === 1 ||
                  page === totalPages ||
                  Math.abs(page - currentPage) <= 2;

                if (!shouldShow) {
                  if (page === currentPage - 3 || page === currentPage + 3) {
                    return (
                      <span key={page} className="px-2 text-sm text-gray-400">
                        ...
                      </span>
                    );
                  }
                  return null;
                }

                return (
                  <button
                    key={page}
                    onClick={() => setCurrentPage(page)}
                    className={`px-3 py-1 text-sm border border-gray-300 rounded ${
                      isCurrentPage
                        ? "bg-blue-600 text-white border-blue-600"
                        : "hover:bg-gray-50"
                    }`}
                  >
                    {page}
                  </button>
                );
              })}
              <button
                onClick={() =>
                  setCurrentPage(Math.min(totalPages, currentPage + 1))
                }
                disabled={currentPage === totalPages}
                className="px-3 py-1 text-sm border border-gray-300 rounded disabled:opacity-50 hover:bg-gray-50"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
