import React, { FC, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import {
  getAllTransactions,
  mapTransactionCategoryToString,
} from "@api/transactions";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import { useFilter } from "@context/index";
import { columns } from "./constants";

const TransactionsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllTransactions(queryOptions).then((res) => {
      changeTotalPages(res.data.totalPages);
      const data = res.data.data.map((item: any) => ({
        time: dateToString(item.time, false, true),
        category: mapTransactionCategoryToString(item.category),
        comment: item.comment,
        sum: formatNumber(item.sum, true, false),
        schet: item.schet,
        id: item.id,
      }));
      setTableData(data);
      changeTotalResults(data.length);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Транзакции"
      addHref="/transactions/transaction_form"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.SCHET]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          isRowDeletable={(row) =>
            !(
              row.original.category === "Закрытие смены" ||
              row.original.category === "Открытие смены"
            )
          }
          isRowEditable={(row) =>
            !(
              row.original.category === "Закрытие смены" ||
              row.original.category === "Открытие смены"
            )
          }
        />
      )}
    </MainLayout>
  );
};

export default TransactionsContent;
