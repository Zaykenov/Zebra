import React, { FC, useCallback, useMemo } from "react";
import { ChevronLeftIcon, ChevronRightIcon } from "@heroicons/react/24/outline";
import clsx from "clsx";
import { useFilter } from "@context/index";

interface PaginationProps {
  detailed?: boolean;
}

const PER_PAGE = 20;

const Pagination: FC<PaginationProps> = ({ detailed = false }) => {
  const {
    totalPages,
    totalResults,
    curPage,
    toNextPage,
    toPreviousPage,
    toPage,
  } = useFilter();

  const current: number = useMemo(() => curPage || 1, []);

  const beginNum = useMemo(() => (current - 1) * PER_PAGE + 1, [curPage]);
  const endNum = useMemo(
    () => (current - 1) * PER_PAGE + totalResults,
    [current, totalResults],
  );

  const handlePageSwitch = useCallback(
    (pageNum: number) => () => {
      toPage(pageNum);
    },
    [toPage],
  );

  return (
    <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
      <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
        {detailed && (
          <div>
            <p className="text-sm text-gray-700">
              Показано с <span className="font-medium">{beginNum}</span> по{" "}
              <span className="font-medium">{endNum}</span> результатов
            </p>
          </div>
        )}
        <div>
          <nav
            className="isolate max-w-3xl overflow-auto inline-flex -space-x-px rounded-md shadow-sm"
            aria-label="Pagination"
          >
            <button
              type="button"
              onClick={toPreviousPage}
              disabled={curPage === 1}
              className="relative inline-flex items-center rounded-l-md border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20"
            >
              <span className="sr-only">Назад</span>
              <ChevronLeftIcon className="h-5 w-5" aria-hidden="true" />
            </button>
            {Array.from({ length: totalPages || 1 }, (_, i) => i + 1).map(
              (pageNum) => (
                <button
                  type="button"
                  onClick={handlePageSwitch(pageNum)}
                  aria-current="page"
                  className={clsx([
                    "relative inline-flex items-center border px-4 py-2 text-sm font-medium focus:z-20",
                    curPage === pageNum
                      ? "z-10 bg-indigo-50 border-indigo-500 text-indigo-600"
                      : "bg-white border-gray-300 text-gray-500 hover:bg-gray-50",
                  ])}
                >
                  {pageNum}
                </button>
              ),
            )}
            <button
              type="button"
              onClick={toNextPage}
              disabled={curPage === totalPages}
              className="relative inline-flex items-center rounded-r-md border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20"
            >
              <span className="sr-only">Далее</span>
              <ChevronRightIcon className="h-5 w-5" aria-hidden="true" />
            </button>
          </nav>
        </div>
      </div>
    </div>
  );
};

export default Pagination;
