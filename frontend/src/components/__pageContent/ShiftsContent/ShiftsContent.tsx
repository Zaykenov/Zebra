import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import { getAllShifts } from "@api/shifts";
import SubRow from "./SubRow";
import { columns } from "./constants";
import { useFilter } from "@context/index";

const ShiftsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const handleGetAll = useCallback((res: any) => {
    changeTotalPages(res.data.totalPages);
    const data = res.data.data.map((shift: any) => ({
      ...shift,
      created_at: dateToString(shift.created_at, false),
      closed_at: shift.is_closed
        ? dateToString(shift.closed_at, false)
        : "Не закрыта",
      begin_sum: formatNumber(shift.begin_sum, true, true),
      collection: formatNumber(shift.collection, true, true),
      end_sum_fact: formatNumber(shift.end_sum_fact, true, true),
      difference: formatNumber(shift.difference, true, true),
    }));
    setTableData(data);
    changeTotalResults(data.length);
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllShifts(queryOptions).then(handleGetAll);
  }, [queryOptions]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  return (
    <MainLayout
      title="Кассовые смены"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.SHOP]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          expandOnRowClick
          editable={false}
          customRowStyle={(row) =>
            row.original.is_closed ? "bg-red-200/80" : "bg-orange-100/80"
          }
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default ShiftsContent;
