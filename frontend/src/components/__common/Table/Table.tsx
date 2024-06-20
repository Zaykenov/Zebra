import React, { FC, ReactNode, useState } from "react";
import Link from "next/link";
import clsx from "clsx";
import DeleteConfirmationModal from "@modules/DeleteConfirmationModal/DeleteConfirmationModal";

import { ChevronDownIcon } from "@heroicons/react/24/outline";
import {
  Cell,
  Column,
  HeaderGroup,
  Row,
  useSortBy,
  useTable,
} from "react-table";
import { useRouter } from "next/router";
import { deleteResource } from "@api/index";
import { useFilter } from "@context/index";

export interface TableProps {
  columns: Column[];
  data: any;
  editable?: boolean;
  onlyEditable?: boolean;
  onlyDeletable?: boolean;
  isRowDeletable?: (row: Row<any>) => boolean;
  isRowEditable?: (row: Row<any>) => boolean;
  details?: boolean;
  renderRowSubComponent?: (row: any) => ReactNode;
  customEditBtn?: (row: Row<any>) => ReactNode | false;
  customDeleteBtn?: (row: Row<any>) => ReactNode | false;
  onRowClick?: (row: Row<any>) => void;
  customRowStyle?: (row: Row<any>) => string;
  expandOnRowClick?: boolean;
  customEditPath?: string;
  onHeaderClick?: (header: HeaderGroup) => void;
  hasFooter?: boolean;
  isDetailedConfirmation?: boolean;
  enableSorting?: boolean;
  customDeleteText?: string;
  customOnDelete?: (row: Row) => void;
  deleteConfirmationText?: string;
  renderCellStyle?: (cell: Cell<{}>) => string;
}

const Table: FC<TableProps> = ({
  columns,
  data,
  editable = true,
  details = false,
  renderRowSubComponent,
  customEditBtn,
  customDeleteBtn,
  customDeleteText,
  deleteConfirmationText,
  onlyDeletable = false,
  isRowDeletable = () => true,
  isRowEditable = () => true,
  onlyEditable = false,
  onRowClick,
  customRowStyle,
  expandOnRowClick = false,
  customEditPath,
  onHeaderClick,
  hasFooter = false,
  isDetailedConfirmation,
  customOnDelete,
  renderCellStyle,
}) => {
  const router = useRouter();

  const { handleSort } = useFilter();

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    footerGroups,
    rows,
    prepareRow,
    visibleColumns,
  } = useTable(
    {
      columns,
      data,
    },
    useSortBy,
  );
  const [isDeleteConfirmationModalOpen, setIsDeleteConfirmationModalOpen] =
    useState<boolean>(false);
  const [expandedRows, setExpandedRows] = useState<string[]>([]);
  const [deletableRow, setDeletableRow] = useState<Row<any>>();

  return (
    <table {...getTableProps()} className="w-full divide-y divide-gray-300">
      <thead className="">
        {headerGroups.map((headerGroup) => (
          <tr
            {...headerGroup.getHeaderGroupProps()}
            key={`head_${headerGroup.id}`}
          >
            {headerGroup.headers.map((column) => {
              return (
                <th
                  {...column.getHeaderProps()}
                  key={`${router.pathname}_column_${column.id}`}
                  onClick={() => {
                    onHeaderClick
                      ? onHeaderClick(column)
                      : handleSort(column.id);
                  }}
                  scope="col"
                  className="sticky top-0 z-10 bg-gray-50 group py-3.5 pl-4 pr-3 text-left text-xs font-normal hover:bg-indigo-100 text-gray-500 pl-6 cursor-pointer"
                >
                  <div className="inline-flex justify-between w-full">
                    {column.render("Header")}
                    <span className="invisible ml-2 flex-none rounded text-gray-400 group-hover:visible group-focus:visible">
                      <ChevronDownIcon className="h-4 w-4" aria-hidden="true" />
                    </span>
                  </div>
                </th>
              );
            })}
            {details && (
              <>
                <th
                  key={`details_column_${headerGroup.id}`}
                  scope="col"
                  className="sticky top-0 z-10 bg-gray-50 group w-10 py-3.5 pl-4 pr-3 text-left text-xs font-normal text-gray-500 pl-6"
                >
                  <div className="inline-flex"></div>
                </th>
              </>
            )}
            {editable && (
              <>
                {!onlyDeletable && (
                  <th
                    key={`edit_column_${headerGroup.id}`}
                    scope="col"
                    className="sticky top-0 z-10 bg-gray-50 group w-10 py-3.5 pl-4 pr-3 text-left text-xs font-normal text-gray-500 pl-6"
                  >
                    <div className="inline-flex"></div>
                  </th>
                )}
                {!onlyEditable && (
                  <th
                    key={`del_column_${headerGroup.id}`}
                    scope="col"
                    className="sticky top-0 z-10 bg-gray-50 group w-10 py-3.5 pl-4 pr-3 text-left text-xs font-normal text-gray-500 pl-6"
                  >
                    <div className="inline-flex"></div>
                  </th>
                )}
              </>
            )}
          </tr>
        ))}
      </thead>
      <tbody
        {...getTableBodyProps()}
        className="divide-y divide-gray-200 bg-white"
      >
        {rows.map((row) => {
          prepareRow(row);
          return (
            <>
              <tr
                {...row.getRowProps()}
                onClick={() => {
                  onRowClick && onRowClick(row);
                  expandOnRowClick &&
                    (expandedRows.includes(row.id)
                      ? setExpandedRows((prevState) =>
                          prevState.filter((rowId) => rowId !== row.id),
                        )
                      : setExpandedRows((prevState) => [...prevState, row.id]));
                }}
                key={`row_${row.id}`}
                className={clsx([
                  (expandOnRowClick || onRowClick) && "cursor-pointer",
                  customRowStyle && customRowStyle(row),
                ])}
              >
                {row.cells.map((cell) => {
                  const customCellStyle = renderCellStyle
                    ? renderCellStyle(cell)
                    : "";
                  return (
                    <td
                      {...cell.getCellProps()}
                      key={`cell_${cell.row.id}_${cell.column.id}`}
                      className={`whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6 ${customCellStyle}`}
                    >
                      {cell.render("Cell")}
                    </td>
                  );
                })}
                {details && (
                  <td
                    key={`details_cell_${row.id}`}
                    className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                  >
                    <button
                      onClick={() => {
                        expandedRows.includes(row.id)
                          ? setExpandedRows((prevState) =>
                              prevState.filter((rowId) => rowId !== row.id),
                            )
                          : setExpandedRows((prevState) => [
                              ...prevState,
                              row.id,
                            ]);
                      }}
                      className="text-indigo-500 hover:text-indigo-600 hover:underline"
                    >
                      Детали
                    </button>
                  </td>
                )}
                {editable && (
                  <>
                    <td
                      key={`edit_cell_${row.id}`}
                      className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                    >
                      {!onlyDeletable &&
                        isRowEditable(row) &&
                        ((customEditBtn && customEditBtn(row)) || (
                          <Link
                            href={`${customEditPath || router.pathname}/${
                              // @ts-ignore
                              row.original.id
                            }`}
                          >
                            <a className="text-indigo-500 hover:text-indigo-600 hover:underline">
                              Ред.
                            </a>
                          </Link>
                        ))}
                    </td>
                    <td
                      key={`del_cell_${row.id}`}
                      className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-6"
                    >
                      {!onlyEditable &&
                        isRowDeletable(row) &&
                        ((customDeleteBtn && customDeleteBtn(row)) || (
                          <button
                            className="text-red-500 hover:text-red-600 hover:underline"
                            onClick={() => {
                              setDeletableRow(row);
                              setIsDeleteConfirmationModalOpen(true);
                            }}
                          >
                            {customDeleteText ? customDeleteText : "Удалить"}
                          </button>
                        ))}
                    </td>
                  </>
                )}
              </tr>
              {/*  @ts-ignore */}
              {expandedRows.includes(row.id) ? (
                <tr key={`expanded_${row.id}`} className="bg-gray-100">
                  <td colSpan={visibleColumns.length}>
                    {renderRowSubComponent && renderRowSubComponent({ row })}
                  </td>
                </tr>
              ) : null}
            </>
          );
        })}
      </tbody>
      {hasFooter && (
        <tfoot>
          {footerGroups.map((group) => (
            <tr
              {...group.getFooterGroupProps()}
              key={`footer_${group.id}`}
              className="text-base font-normal bg-gray-100 pl-6"
            >
              {group.headers.map((column) => (
                <td
                  {...column.getFooterProps()}
                  className="py-3.5 pl-6 pr-3 text-left"
                >
                  {column.render("Footer")}
                </td>
              ))}
              {details && (
                <td
                  key={`footer_details_${group.id}`}
                  className="py-3.5 pl-6 pr-3 text-left"
                ></td>
              )}
              {editable && (
                <>
                  <td
                    key={`footer_edit_${group.id}`}
                    className="py-3.5 pl-6 pr-3 text-left"
                  ></td>
                  {!onlyEditable && !onlyDeletable && (
                    <td
                      key={`footer_del_${group.id}`}
                      className="py-3.5 pl-6 pr-3 text-left"
                    ></td>
                  )}
                </>
              )}
            </tr>
          ))}
        </tfoot>
      )}
      <DeleteConfirmationModal
        isOpen={isDeleteConfirmationModalOpen}
        setIsOpen={setIsDeleteConfirmationModalOpen}
        row={deletableRow}
        deleteConfirmationText={deleteConfirmationText}
        path={router.pathname}
        isDetailed={isDetailedConfirmation}
        onDelete={(e) => {
          e.preventDefault();
          e.stopPropagation();
          customOnDelete && deletableRow && customOnDelete(deletableRow);
          //@ts-ignore
          deleteResource(
            // @ts-ignore
            deletableRow.original.id,
            customEditPath ? customEditPath : router.pathname,
          ).then(() => {
            router.reload();
          });
        }}
      />
    </table>
  );
};

export default Table;
