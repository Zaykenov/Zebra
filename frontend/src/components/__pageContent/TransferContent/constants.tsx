import { Column } from "react-table";
import { formatNumber } from "@utils/formatNumber";
import React from "react";

export const columns: Column[] = [
  {
    Header: "Дата",
    accessor: "time",
  },
  {
    Header: "Наименование",
    accessor: "item_transfers",
  },
  {
    Header: "Сумма",
    accessor: "sum",
    Cell: ({ value }) => <span>{formatNumber(value, true, true)}</span>,
  },
  {
    Header: "Сотрудник",
    accessor: "worker_name",
  },
  {
    Header: "Склады",
    accessor: "sklads",
    Cell: ({ value }) => (
      <span>
        {value[0]} → {value[1]}
      </span>
    ),
  },
];
