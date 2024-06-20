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
    Header: "Себестоимость",
    accessor: "cost",
  },
  {
    Header: "Ед. измерения",
    accessor: "measure",
  },
  {
    Header: "Налог",
    accessor: "tax",
  },
  {
    Header: "Цена",
    accessor: "price",
  },
  {
    Header: "Наценка",
    accessor: "margin",
  },
  {
    Header: "Прибыль",
    accessor: "profit",
  },
];
