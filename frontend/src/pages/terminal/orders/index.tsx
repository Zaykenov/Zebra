import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import Table from "@common/Table";
import {
  dateToString,
  getAllChecks,
  getAllWorkerChecks,
} from "@api/check";
import TerminalLayout from "@layouts/TerminalLayout";
import { useRouter } from "next/router";
import { Column, Row } from "react-table";

const columns: Column<any>[] = [
  { Header: "Открыто", accessor: "opened_at" },
  { Header: "Состав", accessor: "items" },
  { Header: "Статус", accessor: "status" },
  { Header: "Сумма", accessor: "sum" },
];

const techAndTovarCartCheckArrayToString = (
  techCartCheck: [],
  tovarCartCheck: []
) => {
  const techCartCheckString = techCartCheck
    ? techCartCheck.map((techCart: any) => techCart.name).join(", ")
    : "";
  const tovarCartCheckString = tovarCartCheck
    ? tovarCartCheck.map((tovarCart: any) => tovarCart.tovar_name).join(", ")
    : "";
  if (techCartCheckString.length == 0 && tovarCartCheckString.length > 0) {
    return tovarCartCheckString;
  } else if (
    tovarCartCheckString.length == 0 &&
    techCartCheckString.length > 0
  ) {
    return techCartCheckString;
  } else {
    return techCartCheckString + ", " + tovarCartCheckString;
  }
};

const OrdersPage: NextPage = () => {
  const router = useRouter();

  const [tableData, setTableData] = useState<any>([]);

  useEffect(() => {
    getAllWorkerChecks().then((res) => {
      const data = res.data.map((item: any) => ({
        opened_at: dateToString(item.opened_at, false),
        items: techAndTovarCartCheckArrayToString(
          item.techCartCheck,
          item.tovarCheck
        ),
        status: item.status === "opened" ? "Открытый" : "Закрытый",
        sum: item.sum - item.discount,
        check: item,
        id: item.id,
      }));
      console.log(data);
      setTableData(data);
    });
  }, []);

  const onRowClick = (row: Row<any>) => {
    if (row.original.status === "Закрытый") return;
    if (typeof window !== "undefined") {
      localStorage.setItem(
        "zebra.activeCheck",
        JSON.stringify(row.original.check)
      );
    }
    router.push({
      pathname: "/terminal/order",
      query: {
        id: row.original.id,
      },
    });
  };

  const isRowDeletable = (row: Row<any>) => {
    return row.original.status === "Открытый";
  };

  const customRowStyle = (row: Row<any>) => {
    return row.original.status === "Открытый"
      ? "bg-red-50 hover:bg-red-100"
      : "hover:bg-indigo-100";
  };

  return (
    <TerminalLayout>
      <div className="flex w-full overflow-auto">
        <Table
          columns={columns}
          data={tableData}
          editable={true}
          onRowClick={onRowClick}
          onlyDeletable
          isRowDeletable={isRowDeletable}
          customRowStyle={customRowStyle}
        />
      </div>
    </TerminalLayout>
  );
};

export default OrdersPage;
