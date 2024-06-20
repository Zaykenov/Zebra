import { Column } from "react-table";
import React from "react";

export const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Заведение",
    accessor: "shop_name",
  },
  {
    Header: "Мин.",
    accessor: "min",
  },
  {
    Header: "Макс.",
    accessor: "max",
  },
  {
    Header: "Ингредиенты",
    accessor: "nabor_ingredient",
    Cell: ({ value }) => (
      <span className="whitespace-normal">
        {value?.map((item: any) => item.name)?.join(", ")}
      </span>
    ),
  },
];
