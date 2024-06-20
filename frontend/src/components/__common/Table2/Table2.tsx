import React, { FC, Fragment, ReactNode } from "react";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import {
  useReactTable,
  ColumnDef,
  getCoreRowModel,
  getExpandedRowModel,
  flexRender,
  Row,
} from "@tanstack/react-table";
import Link from "next/link";
import { useRouter } from "next/router";
import { deleteResource } from "@api/index";

export interface TableProps<TData> {
  columns: ColumnDef<TData>[];
  data: any;
  editable?: boolean;
  renderSubComponent?: (props: { row: Row<TData> }) => ReactNode;
  getRowCanExpand: (row: Row<TData>) => boolean;
  customOnEdit?: (row: Row<TData>) => void;
}

const Table: FC<TableProps<any>> = ({
  columns,
  data,
  editable = true,
  renderSubComponent,
  getRowCanExpand,
  customOnEdit,
}) => {
  const router = useRouter();

  const { getHeaderGroups, getRowModel } = useReactTable({
    columns,
    data,
    getRowCanExpand,
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
  });

  return (
    <table className="min-w-full divide-y divide-gray-300">
      <thead className="bg-gray-50">
        {getHeaderGroups().map((headerGroup, idx) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((column, idx) => {
              return (
                <th
                  colSpan={column?.colSpan}
                  key={column.id}
                  scope="col"
                  className="group py-3.5 pl-4 pr-3 text-left text-xs font-normal hover:bg-indigo-100 text-gray-500 pl-6 cursor-pointer"
                >
                  <div className="inline-flex justify-between w-full">
                    {flexRender(
                      column?.column?.columnDef?.header,
                      column?.getContext()
                    )}
                    <span className="invisible ml-2 flex-none rounded text-gray-400 group-hover:visible group-focus:visible">
                      <ChevronDownIcon className="h-4 w-4" aria-hidden="true" />
                    </span>
                  </div>
                </th>
              );
            })}
            {editable && (
              <>
                <th
                  key={`edit_column_${idx}`}
                  scope="col"
                  className="group w-10 py-3.5 pl-4 pr-3 text-left text-xs font-normal text-gray-500 pl-6"
                >
                  <div className="inline-flex"></div>
                </th>
                <th
                  key={`del_column_${idx}`}
                  scope="col"
                  className="group w-10 py-3.5 pl-4 pr-3 text-left text-xs font-normal text-gray-500 pl-6"
                >
                  <div className="inline-flex"></div>
                </th>
              </>
            )}
          </tr>
        ))}
      </thead>
      <tbody className="divide-y divide-gray-200 bg-white">
        {getRowModel()?.rows?.map((row, idx) => {
          return (
            <Fragment key={row.id}>
              <tr className="hover:bg-gray-100">
                {row?.getVisibleCells()?.map((cell, idx) => {
                  return (
                    <td
                      key={cell.id}
                      className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                    >
                      {flexRender(
                        cell?.column?.columnDef?.cell,
                        cell?.getContext()
                      )}
                    </td>
                  );
                })}
                {editable && (
                  <>
                    <td
                      key={`edit_cell_${idx}`}
                      className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                    >
                      {/*@ts-ignore*/}
                      {customOnEdit ? (
                        <button
                          onClick={() => customOnEdit(row)}
                          className="text-indigo-500 hover:text-indigo-600 hover:underline"
                        >
                          Ред.
                        </button>
                      ) : (
                        <Link href={`${router.pathname}/${row.original.id}`}>
                          <a className="text-indigo-500 hover:text-indigo-600 hover:underline">
                            Ред.
                          </a>
                        </Link>
                      )}
                    </td>
                    <td
                      key={`del_cell_${idx}`}
                      className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                    >
                      <button
                        onClick={() => {
                          // @ts-ignore
                          deleteResource(row.original.id, router.pathname).then(
                            () => {
                              router.reload();
                            }
                          );
                        }}
                        className="text-red-500 hover:text-red-600 hover:underline"
                      >
                        Удалить
                      </button>
                    </td>
                  </>
                )}
              </tr>
              {row?.getIsExpanded() && (
                <tr>
                  <td
                    className="bg-gray-100"
                    colSpan={row?.getVisibleCells()?.length}
                  >
                    {renderSubComponent && renderSubComponent({ row })}
                  </td>
                </tr>
              )}
            </Fragment>
          );
        })}
      </tbody>
    </table>
  );
};

export default Table;
