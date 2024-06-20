import { Column } from "react-table";

export const columns: Column[] = [
  {
    Header: "Дата",
    accessor: "time",
  },
  {
    Header: "Категория",
    accessor: "category",
  },
  {
    Header: "Комментарий",
    accessor: "comment",
  },
  {
    Header: "Сумма",
    accessor: "sum",
  },
  {
    Header: "Счет",
    accessor: "schet",
  },
];
