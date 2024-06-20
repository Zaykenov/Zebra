import { Column } from "react-table";
import { formatNumber } from "@utils/formatNumber";
import React from "react";
import DetailsPopover from "./DetailsPopover";

export const columns: (total: {
  initialSum: number;
  finalSum: number;
}) => Column[] = (total) => [
  {
    Header: "Название",
    accessor: "item_name",
    Footer: "Итого",
  },

  {
    Header: "Тип",
    accessor: "type",
    Cell: ({ value }) => (
      <span>{value === "tovar" ? "Товар" : "Ингредиент"}</span>
    ),
  },
  {
    Header: "Нач. остаток",
    accessor: "initial_ostatki",
    Cell: ({ value, row }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, false, false)} {row.original.measure}
      </span>
    ),
  },
  {
    Header: "Средняя себест. на начало",
    accessor: "initial_netCost",
    Cell: ({ value }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, true, true)}
      </span>
    ),
  },
  {
    Header: "Нач. сумма",
    accessor: "initial_sum",
    Cell: ({ value }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, true, true)}
      </span>
    ),
    Footer: () => (
      <span className="whitespace-nowrap">
        {formatNumber(total.initialSum, true, true)}
      </span>
    ),
  },
  {
    Header: "Поступления",
    accessor: "income",
    Cell: ({ value, row }) => (
      <DetailsPopover
        value={value}
        row={row}
        details={{
          // @ts-ignore
          measure: row.original.measure,
          // @ts-ignore
          postavka: row.original.postavka,
          // @ts-ignore
          postavka_cost: row.original.postavka_cost,
          // @ts-ignore
          ...(row.original.inventarization > 0
            ? // @ts-ignore
              { inventarization: row.original.inventarization }
            : {}),
          // @ts-ignore
          ...(row.original.transfer > 0
            ? // @ts-ignore
              { transfer: row.original.transfer }
            : {}),
        }}
      />
    ),
  },
  {
    Header: "Расход",
    accessor: "consumption",
    Cell: ({ value, row }) => (
      <DetailsPopover
        value={value}
        row={row}
        details={{
          // @ts-ignore
          measure: row.original.measure,
          // @ts-ignore
          sales: row.original.sales,
          // @ts-ignore
          ...(row.original.inventarization < 0
            ? // @ts-ignore
              { inventarization: row.original.inventarization }
            : {}),
          // @ts-ignore
          ...(row.original.transfer < 0
            ? // @ts-ignore
              { transfer: row.original.transfer }
            : {}),
        }}
      />
    ),
  },
  {
    Header: "Итог. остаток",
    accessor: "final_ostatki",
    Cell: ({ value, row }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, false, false)} {row.original.measure}
      </span>
    ),
  },
  {
    Header: "Средняя себест. на конец",
    accessor: "final_netCost",
    Cell: ({ value }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, true, true)}
      </span>
    ),
  },
  {
    Header: "Итог. сумма",
    accessor: "final_sum",
    Cell: ({ value }) => (
      <span>
        {/* @ts-ignore */}
        {formatNumber(value, true, true)}
      </span>
    ),
    Footer: () => (
      <span className="whitespace-nowrap">
        {formatNumber(total.finalSum, true, true)}
      </span>
    ),
  },
];
