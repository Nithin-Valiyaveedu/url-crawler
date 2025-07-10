import { ArrowLeft, ExternalLink, AlertTriangle } from "lucide-react";
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { CrawlResult } from "../types";

interface CrawlResultDetailProps {
  result: CrawlResult;
  onBack: () => void;
}

const COLORS = {
  internal: "#3B82F6",
  external: "#10B981",
  broken: "#EF4444",
};

export const CrawlResultDetail: React.FC<CrawlResultDetailProps> = ({
  result,
  onBack,
}) => {
  const linkData = [
    {
      name: "Internal Links",
      value: result.internalLinksCount,
      color: COLORS.internal,
    },
    {
      name: "External Links",
      value: result.externalLinksCount,
      color: COLORS.external,
    },
  ];

  const headingData = [
    { level: "H1", count: result.headingCounts.h1 },
    { level: "H2", count: result.headingCounts.h2 },
    { level: "H3", count: result.headingCounts.h3 },
    { level: "H4", count: result.headingCounts.h4 },
    { level: "H5", count: result.headingCounts.h5 },
    { level: "H6", count: result.headingCounts.h6 },
  ].filter((item) => item.count > 0);

  const totalLinks = result.internalLinksCount + result.externalLinksCount;

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  };

  const getStatusBadge = (status: CrawlResult["status"]) => {
    const baseClasses =
      "inline-flex items-center px-3 py-1 rounded-full text-sm font-medium";
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

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-0">
      <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-6">
        <div className="flex flex-col sm:flex-row sm:items-center gap-3 sm:gap-4 mb-4">
          <button
            onClick={onBack}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors self-start"
          >
            <ArrowLeft className="w-5 h-5" />
            <span className="text-sm sm:text-base cursor-pointer">
              Back to Results
            </span>
          </button>
          <span className={getStatusBadge(result.status)}>
            {result.status.toUpperCase()}
          </span>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div>
            <h1 className="text-xl sm:text-2xl font-bold text-gray-900 mb-3">
              {result.title || "Untitled Page"}
            </h1>
            <div className="mb-4">
              <div className="flex items-start gap-2">
                <a
                  href={result.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:text-blue-800 text-sm sm:text-base break-all flex-1"
                >
                  {result.url}
                </a>
                <ExternalLink className="w-4 h-4 mt-0.5 flex-shrink-0" />
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-gray-500 block">HTML Version:</span>
              <p className="font-medium">{result.htmlVersion || "Unknown"}</p>
            </div>
            <div>
              <span className="text-gray-500 block">Login Form:</span>
              <p className="font-medium">
                {result.hasLoginForm ? "Present" : "Not Found"}
              </p>
            </div>
            <div>
              <span className="text-gray-500 block">Crawled:</span>
              <p className="font-medium">{formatDate(result.createdAt)}</p>
            </div>
            <div>
              <span className="text-gray-500 block">Last Updated:</span>
              <p className="font-medium">{formatDate(result.updatedAt)}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6 mb-6">
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 text-center">
          <div className="text-xl sm:text-2xl font-bold text-blue-600">
            {result.internalLinksCount}
          </div>
          <div className="text-xs sm:text-sm text-gray-600">Internal Links</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 text-center">
          <div className="text-xl sm:text-2xl font-bold text-green-600">
            {result.externalLinksCount}
          </div>
          <div className="text-xs sm:text-sm text-gray-600">External Links</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 text-center">
          <div className="text-xl sm:text-2xl font-bold text-red-600">
            {result.inaccessibleLinksCount}
          </div>
          <div className="text-xs sm:text-sm text-gray-600">Broken Links</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 text-center">
          <div className="text-xl sm:text-2xl font-bold text-purple-600">
            {Object.values(result.headingCounts).reduce((a, b) => a + b, 0)}
          </div>
          <div className="text-xs sm:text-sm text-gray-600">Total Headings</div>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 mb-6">
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6">
          <h3 className="text-base sm:text-lg font-semibold text-gray-800 mb-4">
            Links Distribution
          </h3>
          {totalLinks > 0 ? (
            <div className="h-48 sm:h-64">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={linkData}
                    cx="50%"
                    cy="50%"
                    innerRadius={40}
                    outerRadius={window.innerWidth < 640 ? 70 : 100}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {linkData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value, name) => [`${value} links`, name]}
                    labelFormatter={() => ""}
                  />
                </PieChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="h-48 sm:h-64 flex items-center justify-center text-gray-500">
              <p className="text-sm text-center">No link data available</p>
            </div>
          )}
          <div className="flex flex-col sm:flex-row justify-center gap-3 sm:gap-6 mt-4">
            <div className="flex items-center gap-2 justify-center sm:justify-start">
              <div className="w-3 h-3 rounded-full bg-blue-500"></div>
              <span className="text-xs sm:text-sm text-gray-600">
                Internal ({result.internalLinksCount})
              </span>
            </div>
            <div className="flex items-center gap-2 justify-center sm:justify-start">
              <div className="w-3 h-3 rounded-full bg-green-500"></div>
              <span className="text-xs sm:text-sm text-gray-600">
                External ({result.externalLinksCount})
              </span>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6">
          <h3 className="text-base sm:text-lg font-semibold text-gray-800 mb-4">
            Heading Tags Distribution
          </h3>
          {headingData.length > 0 ? (
            <div className="h-48 sm:h-64">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={headingData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="level" fontSize={12} />
                  <YAxis fontSize={12} />
                  <Tooltip
                    formatter={(value) => [`${value}`, "Count"]}
                    labelFormatter={(label) => `${label} Tags`}
                  />
                  <Bar dataKey="count" fill="#8B5CF6" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="h-48 sm:h-64 flex items-center justify-center text-gray-500">
              <p className="text-sm text-center">No heading data available</p>
            </div>
          )}
        </div>
      </div>

      {result.inaccessibleLinksCount > 0 && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center gap-2 mb-4">
            <AlertTriangle className="w-5 h-5 text-red-500" />
            <h3 className="text-lg font-semibold text-gray-800">
              Broken Links ({result.brokenLinks.length})
            </h3>
          </div>

          {result.brokenLinks.length > 0 ? (
            <div className="space-y-3">
              {result.brokenLinks.map((link, index) => (
                <div
                  key={index}
                  className="border border-red-200 rounded-lg p-4 bg-red-50"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-sm font-mono text-red-600">
                          {link.statusCode}
                        </span>
                        <span className="text-sm text-red-700">
                          {link.statusText}
                        </span>
                      </div>
                      <p className="text-sm text-gray-700 break-all">
                        {link.url}
                      </p>
                    </div>
                    <a
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-red-600 hover:text-red-800 flex-shrink-0"
                    >
                      <ExternalLink className="w-4 h-4" />
                    </a>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">
              No broken link details available, but{" "}
              {result.inaccessibleLinksCount} broken links were detected.
            </p>
          )}
        </div>
      )}

      {result.externalLinksCount > 0 && (
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-6">
          <div className="flex items-center gap-2 mb-4">
            <ExternalLink className="w-5 h-5 text-green-500" />
            <h3 className="text-base sm:text-lg font-semibold text-gray-800">
              External Links (
              {result.externalLinks?.length || result.externalLinksCount})
            </h3>
          </div>

          {result.externalLinks && result.externalLinks.length > 0 ? (
            <div className="space-y-3">
              {result.externalLinks.map((link, index) => (
                <div
                  key={index}
                  className="border border-green-200 rounded-lg p-3 sm:p-4 bg-green-50 hover:bg-green-100 transition-colors"
                >
                  <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4">
                    <div className="flex-1 min-w-0">
                      <p className="text-xs sm:text-sm text-gray-700 break-all">
                        {link}
                      </p>
                    </div>
                    <a
                      href={link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-green-600 hover:text-green-800 flex-shrink-0 self-start sm:self-center"
                      title="Open in new tab"
                    >
                      <ExternalLink className="w-4 h-4" />
                    </a>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500">
              No external link details available, but{" "}
              {result.externalLinksCount} external links were detected.
            </p>
          )}
        </div>
      )}
    </div>
  );
};
