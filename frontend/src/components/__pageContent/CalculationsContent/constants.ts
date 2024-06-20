import { Column } from "react-table";

export const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Склад",
    accessor: "skladName",
  },
  {
    Header: "Тип",
    accessor: "type",
  },
  {
    Header: "Категория",
    accessor: "category",
  },
  {
    Header: "Себестоимость",
    accessor: "cost",
  },
  {
    Header: "Остатки",
    accessor: "quantity",
  },
  {
    Header: "Ед. измерения",
    accessor: "measurement",
  },
  {
    Header: "Сумма",
    accessor: "sum",
  },
];
