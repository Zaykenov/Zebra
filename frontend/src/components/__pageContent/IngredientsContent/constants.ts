import { Column } from "react-table";

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
    Header: "Категория",
    accessor: "category",
  },
  {
    Header: "Ед. измерения",
    accessor: "measure",
  },
  {
    Header: "Себестоимость",
    accessor: "cost",
  },
];
