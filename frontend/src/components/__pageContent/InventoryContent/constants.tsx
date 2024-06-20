import React from "react";
import { Column } from "react-table";
import clsx from "clsx";
import { formatNumber } from "@utils/formatNumber";

export const columns: Column[] = [
  {
    Header: "Склад",
    accessor: "sklad",
  },
  {
    Header: "Дата и время проведения",
    accessor: "time",
  },
  {
    Header: "Тип",
    accessor: "type",
  },
  {
    Header: "Результат",
    accessor: "result",
    Cell: ({ value }) => (
      <span className={clsx([value < 0 && "text-red-500"])}>
        {formatNumber(value, true, true)}
      </span>
    ),
  },
  {
    Header: "Статус",
    accessor: "status",
  },
];
