import { Cell, Column } from "react-table";

export const columns: Column[] = [
  {
    Header: "Топ",
    accessor: "id",
  },
  {
    Header: "Товар",
    accessor: "product",
  },
  {
    Header: "Продажи",
    accessor: "sales",
  },
  {
    Header: "Продажи %",
    accessor: "sales_percent",
  },
  {
    Header: "Выручка",
    accessor: "revenue",
  },
  {
    Header: "Выручка %",
    accessor: "revenue_percent",
  },
  {
    Header: "Прибыль",
    accessor: "profit",
  },
  {
    Header: "Прибыль %",
    accessor: "profit_percent",
  },
  {
    Header: "Продажи",
    accessor: "sales_letter",
  },
  {
    Header: "Выручка",
    accessor: "revenue_letter",
  },
  {
    Header: "Прибыль",
    accessor: "profit_letter",
  },
];

export const renderCellStyle = (cell: Cell<{}, any>): string => {
  switch (cell.value) {
    case "A":
      return "bg-green-100";
    case "B":
      return "bg-yellow-100";
    case "C":
      return "bg-red-100";
    default:
      return "";
  }
};
