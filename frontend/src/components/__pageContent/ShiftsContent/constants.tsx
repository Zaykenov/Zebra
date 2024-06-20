import { Column } from "react-table";
import { formatNumber } from "@utils/formatNumber";
import React from "react";
import { dateToString } from "@api/check";

export const columns: Column[] = [
  {
    Header: "#",
    accessor: "id",
  },
  {
    Header: "Заведение",
    accessor: "shop",
  },
  {
    Header: "Смена открыта",
    accessor: "created_at",
  },
  {
    Header: "Смена закрыта",
    accessor: "closed_at",
  },
  {
    Header: "Начало смены",
    accessor: "begin_sum",
    Cell: (cellProps) => {
      // @ts-ignore
      const diff = cellProps.row.original.difference_with_previous;
      return (
        <div className="w-2/3 flex items-center justify-between space-x-2">
          <div>{cellProps.value}</div>
          {/* @ts-ignore */}
          {!cellProps.row.original.is_equal && (
            <div className="relative group w-1.5 h-1.5 ml-5 rounded-full bg-red-500">
              <div className="absolute z-20 -left-28 top-3 bg-white border p-4 border-gray-100 hidden group-hover:flex flex-col items-center space-y-2">
                <span className="text-xs font-light">
                  Расхождение с предыдущей сменой
                </span>
                <span className="text-lg">
                  {diff > 0 && "+"}
                  {formatNumber(diff, true, true)}
                </span>
              </div>
            </div>
          )}
        </div>
      );
    },
  },
  {
    Header: "Инкассация",
    accessor: "collection",
  },
  {
    Header: "В кассе",
    accessor: "end_sum_fact",
  },
  {
    Header: "Разница",
    accessor: "difference",
  },
];

export const firstRow = [
  {
    text: "Книжный баланс",
    field: "end_sum_plan",
    color: "text-black",
  },
  {
    text: "Фактический баланс",
    field: "end_sum_fact",
    color: "text-black",
  },
  {
    text: "Разница",
    field: "difference",
    color: "text-black",
  },
];

export const secondRow = [
  {
    text: "Наличная выручка",
    field: "cash",
    color: "text-green-800",
  },
  {
    text: "Безналичная выручка",
    field: "card",
    color: "text-green-800",
  },
  {
    text: "Приходы",
    field: null,
    color: "text-green-800",
  },
  {
    text: "Расходы",
    field: "expense",
    color: "text-red-800",
  },
  {
    text: "Инкассация",
    field: "collection",
    color: "text-red-800",
  },
];

export const transactionColumns: Column[] = [
  {
    Header: "Категория",
    accessor: "category",
  },
  {
    Header: "Время",
    accessor: "time",
  },
  {
    Header: "Сумма",
    accessor: "sum",
  },
  {
    Header: "Сотрудник",
    accessor: "worker",
  },
  {
    Header: "Комментарий",
    accessor: "comment",
  },
  {
    Header: "Редактировал",
    accessor: "updated_worker",
    Cell: ({ value, row }) => {
      // @ts-ignore
      const updatedTime = row.original.updated_time;
      return !value ? (
        <div>---</div>
      ) : (
        <div className="">
          {value}, {dateToString(updatedTime, false, false)}
        </div>
      );
    },
  },
];
