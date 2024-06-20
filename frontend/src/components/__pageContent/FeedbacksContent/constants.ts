import { Column } from "react-table";

export const columns: Column[] = [
  {
    Header: "Пользователь",
    accessor: "username",
  },
  {
    Header: "Склад",
    accessor: "shopName",
  },
  {
    Header: "Кассир",
    accessor: "workerName",
  },
  {
    Header: "Оценка качества",
    accessor: "scoreQuality",
  },
  {
    Header: "Оценка сервиса",
    accessor: "scoreService",
  },
  {
    Header: "Отзыв",
    accessor: "feedbackText",
  },
  {
    Header: "Дата отзыва",
    accessor: "feedbackDate",
  },
];
