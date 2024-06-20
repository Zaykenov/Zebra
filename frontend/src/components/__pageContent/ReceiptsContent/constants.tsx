import { Column } from "react-table";
import { formatNumber } from "@utils/formatNumber";
import clsx from "clsx";
import React from "react";
import { TotalInventoryData } from "./types";

export const columns: (totalData: TotalInventoryData) => Column<any>[] = (
  totalData,
) => [
  {
    Header: "#",
    accessor: "id",
  },

  {
    Header: "Кассир",
    accessor: "worker",
    Footer: () => (
      <div className="text-center flex flex-col items-center">
        <span className="font-semibold text-sm text-gray-500">Оборот:</span>
        <span className="text-xl font-bold whitespace-nowrap">
          {formatNumber(totalData.total_money, true, true)}
        </span>
      </div>
    ),
  },
  {
    Header: "Открыт",
    accessor: "opened_at",
    Footer: () => (
      <div className="text-center flex flex-col items-center">
        <span className="font-semibold text-sm text-gray-500">Наличные:</span>
        <span className="text-xl font-bold whitespace-nowrap">
          {formatNumber(totalData.total_cash, true, true)}
        </span>
      </div>
    ),
  },
  {
    Header: "Закрыт",
    accessor: "closed_at",
    Footer: () => (
      <div className="text-center flex flex-col items-center">
        <span className="font-semibold text-sm text-gray-500">Карточные:</span>
        <span className="text-xl font-bold whitespace-nowrap">
          {formatNumber(totalData.total_card, true, true)}
        </span>
      </div>
    ),
  },
  {
    Header: "Способ оплаты",
    Cell: ({ row }) => {
      if (row.original.card > 0 && row.original.cash > 0) {
        return <span>Смешанный</span>;
      }
      return <span>{row.original.card > 0 ? "Картой" : "Наличными"}</span>;
    },
  },
  {
    Header: "Сумма чека",
    accessor: "sum",
    Footer: () => (
      <div className="text-center flex flex-col items-center">
        <span className="font-semibold text-sm text-gray-500">Скидки:</span>
        <span className="text-xl font-bold whitespace-nowrap">
          {formatNumber(totalData.total_discount, true, true)}
        </span>
      </div>
    ),
  },
  {
    Header: "Себестоимость",
    accessor: "cost",
    Footer: () => (
      <div className="text-center flex flex-col items-center">
        <span className="font-semibold text-sm text-gray-500 whitespace-nowrap">
          Валовая прибыль:
        </span>
        <span className="text-xl font-bold whitespace-nowrap">
          {formatNumber(totalData.total_profit, true, true)}
        </span>
      </div>
    ),
  },
  {
    Header: "Скидка в чеке",
    accessor: "discount",
    Cell: ({ value }) => (
      <span>{value === 0 ? "-" : formatNumber(value, true, true)}</span>
    ),
  },
  {
    Header: "Прибыль",
    accessor: "profit",
    Cell: ({ value }) => (
      <span className={clsx([value < 0 && "text-red-500"])}>
        {formatNumber(value, true)}
      </span>
    ),
  },
  {
    Header: "Статус",
    accessor: "status",
  },
];
