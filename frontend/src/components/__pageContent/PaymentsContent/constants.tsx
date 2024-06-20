import { Column } from "react-table";
import React from "react";
import { formatNumber } from "@utils/formatNumber";
import { TotalPaymentData } from "./types";

export const columns: (totalData: TotalPaymentData) => Column[] = (
  totalData,
) => [
  {
    Header: "Дата",
    accessor: "time",
    Footer: () => <span className="font-semibold">Итого</span>,
  },
  {
    Header: "Количество",
    accessor: "check_count",
    Footer: totalData.total_check_count,
  },
  {
    Header: "Наличными",
    accessor: "cash",
    Footer: formatNumber(totalData.total_cash, true, true),
  },
  {
    Header: "Карточкой",
    accessor: "card",
    Footer: formatNumber(totalData.total_card, true, true),
  },
  {
    Header: "Всего",
    accessor: "total",
    Footer: formatNumber(totalData.total_total, true, true),
  },
];
