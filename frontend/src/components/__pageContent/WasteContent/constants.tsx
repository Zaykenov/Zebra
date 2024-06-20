import { Column, Row } from "react-table";
import { confirmRemoveWaste, rejectRemoveWaste } from "@api/wastes";
import { WasteItem } from "./types";
import { formatNumber } from "@utils/formatNumber";
import React from "react";

export const columns: (total: number) => Column[] = (total) => [
  {
    Header: "Дата",
    accessor: "time",
    Footer: () => <span className="font-semibold">Итого</span>,
  },

  {
    Header: "Склад",
    accessor: "sklad",
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
    accessor: "cost",
  },
  {
    Header: "Причина",
    accessor: "reason",
  },
  {
    Header: "Комментарий",
    accessor: "comment",
  },
  {
    Header: "Статус",
    accessor: "status",
    Footer: () => (
      <span className="whitespace-nowrap font-medium">
        {formatNumber(total, true, true)}
      </span>
    ),
  },
];

export const statusToColor = (row: Row<any>): string => {
  if (row.original.status === "Открыто") {
    return "bg-white-200/80";
  } else if (row.original.status === "Закрыто") {
    return "bg-orange-100/80";
  } else {
    return "bg-red-100/80";
  }
};

export const handleConfirmWasteRemoval = async (row: Row<any>) => {
  try {
    await confirmRemoveWaste(row.original.id);
    window.location.reload();
  } catch (e) {
    console.log(e);
  }
};

export const handleRejectWasteRemoval = async (row: Row<any>) => {
  try {
    await rejectRemoveWaste(row.original.id);
    window.location.reload();
  } catch (e) {
    console.log(e);
  }
};

export const headerRow = [
  "Дата",
  "Склад",
  "Товары",
  "Сумма",
  "Причина",
  "Комментарий",
  "Статус",
];

export const propertyNames: (keyof WasteItem)[] = [
  "date",
  "sklad",
  "items",
  "cost",
  "reason",
  "comment",
  "status",
];

export const statusText = {
  opened: "Открыто",
  closed: "Закрыто",
  rejected: "Отклонено",
};

export const filterProperties = (
  arr: any[],
  properties?: (keyof WasteItem)[],
): Partial<WasteItem>[] => {
  return arr.map((obj) => {
    const filteredObj: Partial<WasteItem> = {};
    if (properties === undefined) {
      properties = Object.keys(obj) as (keyof WasteItem)[];
    }
    properties.forEach((prop) => {
      filteredObj[prop] = obj[prop];
    });
    return filteredObj;
  });
};
