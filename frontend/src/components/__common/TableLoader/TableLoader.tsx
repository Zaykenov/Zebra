import React, { FC } from 'react';

interface TableLoaderProps {
  headerRowNames: string[];
  rowCount: number;
}

const TableLoader: FC<TableLoaderProps> = ({ headerRowNames, rowCount }) => {
  const rows = Array.from({ length: rowCount });
  return (
    <table className="w-full divide-y divide-gray-300">
      <thead>
        <tr>
          {headerRowNames.map((headerName, idx) => (
            <th
              key={idx}
              className="sticky top-0 z-10 bg-gray-50 group py-3.5 pl-4 pr-3 text-left text-xs font-normal hover:bg-indigo-100 text-gray-500 pl-6 cursor-pointer animate-pulse"
            >
              {headerName}
            </th>
          ))}
        </tr>
      </thead>
      <tbody className="divide-y divide-gray-200 bg-white">
        {rows.map((_, idx) => (
          <tr
            key={idx}
            className="animate-pulse h-14"
            style={{
              animationDelay: `${0.1 * idx}s`,
            }}
          >
            {headerRowNames.map((_, idx) => (
              <td
                key={idx}
                className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
              ></td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default TableLoader;