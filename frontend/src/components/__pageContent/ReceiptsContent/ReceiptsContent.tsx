import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { dateToString, getAllChecks } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import SubRow from "./SubRow";
import { TotalInventoryData } from "./types";
import { columns } from "./constants";
import { useFilter } from "@context/index";

const ReceiptsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState<any>([]);

  const [totalData, setTotalData] = useState<TotalInventoryData>({
    total_card: 0,
    total_cash: 0,
    total_discount: 0,
    total_money: 0,
    total_net_cost: 0,
    total_profit: 0,
  });

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllChecks(queryOptions).then((res) => {
      changeTotalPages(res.totalPages);
      const { check, ...totalData } = res.data;
      setTotalData(totalData);

      const data = res.data.check.map((check: any) => ({
        ...check,
        worker_id: check.worker_id,
        opened_at: dateToString(check.opened_at, false),
        closed_at:
          check.status === "closed"
            ? dateToString(check.closed_at, false)
            : "-",
        sum: formatNumber(check.sum - check.discount, true),
        cost: formatNumber(check.cost, true),
        profit: check.sum - check.cost,
        status:
          check.status === "closed"
            ? "Закрыт"
            : check.status === "inactive"
            ? "Возврат"
            : "Открыт",
      }));
      setTableData(data);
      changeTotalResults(data.length);
    });
  }, [changeTotalResults, changeTotalPages, queryOptions]);

  return (
    <MainLayout
      title="Чеки"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.WORKER, FilterOption.SHOP]}
    >
      {tableData && (
        <Table
          columns={columns(totalData)}
          data={tableData}
          onlyDeletable
          hasFooter
          customDeleteText="Возврат"
          deleteConfirmationText="Вы уверены, что хотите сделать возврат этого чека?"
          customDeleteBtn={(row) => {
            return row.original.status !== "Закрыт";
          }}
          details={true}
          renderRowSubComponent={renderSubComponent}
          customRowStyle={(row) =>
            row.original.status === "Открыт"
              ? "bg-amber-100/40"
              : row.original.status === "Возврат"
              ? "bg-red-100"
              : ""
          }
        />
      )}
    </MainLayout>
  );
};

export default ReceiptsContent;
