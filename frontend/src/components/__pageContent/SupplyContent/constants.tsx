import { Column } from "react-table";
import { formatNumber } from "@utils/formatNumber";
import React from "react";
import { SupplyItem } from "./types";

export const filterTopItems = (items: any[]) => {
  return items.sort((a, b) => b.price - a.price).slice(0, 3);
};

export const filterProperties = (
  arr: any[],
  properties?: (keyof SupplyItem)[],
): Partial<SupplyItem>[] => {
  return arr.map((obj) => {
    const filteredObj: Partial<SupplyItem> = {};
    if (properties === undefined) {
      properties = Object.keys(obj) as (keyof SupplyItem)[];
    }
    properties.forEach((prop) => {
      filteredObj[prop] = obj[prop];
    });
    return filteredObj;
  });
};

export const columns: (total: number) => Column[] = (total) => [
  {
    Header: "Дата",
    accessor: "time",
    Footer: () => <span className="font-semibold">Итого</span>,
  },
  {
    Header: "Поставщик",
    accessor: "dealer",
  },
  {
    Header: "Склад",
    accessor: "sklad",
  },
  {
    Header: "Счет",
    accessor: "schet",
  },
  {
    Header: "Товары",
    accessor: "items",
    Cell: ({ value }) => (
      <div className="w-full flex items-center">
        <div className="whitespace-normal">{value}</div>
      </div>
    ),
  },
  {
    Header: "Сумма",
    accessor: "sum",
    Cell: (cellProps) => {
      return (
        <div className="w-full flex items-center justify-between space-x-2">
          <div>{cellProps.value}</div>
        </div>
      );
    },
    Footer: () => (
      <span className="whitespace-nowrap font-medium">
        {formatNumber(total, true, true)}
      </span>
    ),
  },
];

export const headerRow = [
  "Дата",
  "Поставщик",
  "Склад",
  "Счет",
  "Товары",
  "Сумма",
];

export const propertyNames: (keyof SupplyItem)[] = [
  "time",
  "dealer",
  "sklad",
  "schet",
  "items",
  "intSum",
];
