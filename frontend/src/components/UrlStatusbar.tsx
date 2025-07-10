import React from "react";

import type { CrawlResult } from "../types";

export const UrlStatusbar: React.FC<{
  crawlResults: CrawlResult[];
}> = ({ crawlResults }) => {
  return (
    <div className="grid grid-cols-3 gap-2 lg:flex lg:items-center lg:gap-4">
      <div className="bg-blue-50 px-3 py-2 rounded-lg text-center lg:bg-transparent lg:p-0">
        <div className="text-lg font-semibold text-blue-600 lg:text-sm lg:text-gray-500">
          {crawlResults.filter((r) => r.status === "running").length}
        </div>
        <div className="text-xs text-blue-600 lg:hidden">crawling</div>
        <div className="hidden lg:inline text-sm text-gray-500">crawling</div>
      </div>
      <div className="bg-yellow-50 px-3 py-2 rounded-lg text-center lg:bg-transparent lg:p-0">
        <div className="text-lg font-semibold text-yellow-600 lg:text-sm lg:text-gray-500">
          {crawlResults.filter((r) => r.status === "queued").length}
        </div>
        <div className="text-xs text-yellow-600 lg:hidden">queued</div>
        <div className="hidden lg:inline text-sm text-gray-500">queued</div>
      </div>
      <div className="bg-green-50 px-3 py-2 rounded-lg text-center lg:bg-transparent lg:p-0">
        <div className="text-lg font-semibold text-green-600 lg:text-sm lg:text-gray-500">
          {crawlResults.filter((r) => r.status === "completed").length}
        </div>
        <div className="text-xs text-green-600 lg:hidden">completed</div>
        <div className="hidden lg:inline text-sm text-gray-500">completed</div>
      </div>
    </div>
  );
};
