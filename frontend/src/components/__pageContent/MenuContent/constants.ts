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
    Header: "Налог",
    accessor: "tax",
  },
  {
    Header: "Себестоимость с НДС",
    accessor: "cost",
  },
  {
    Header: "Ед. измерения",
    accessor: "measure",
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
