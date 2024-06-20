import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { getAllTransfers } from "@api/transfers";
import { dateToString } from "@api/check";
import SubRow from "./SubRow";
import { columns } from "./constants";
import { useFilter } from "@context/filter.context";

const TransferContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllTransfers(queryOptions).then((res) => {
      const data = res.data.data.map((item: any) => {
        const items = item.item_transfers
          .slice(0, 3)
          .map((removedItem: any) => removedItem.item_name)
          .join(", ");
        return {
          ...item,
          time: dateToString(item.time, false),
          sklads: [item.from_sklad_name, item.to_sklad_name],
          item_transfers: items,
          sum: item.sum,
          reason: item.reason,
          comment: item.comment,
        };
      });
      setTableData(data);
      changeTotalResults(data.length);
      changeTotalPages(res.data.totalPages);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Перемещения"
      addHref="/transfer/transfer_form"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.SKLAD]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          editable={true}
          details={true}
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default TransferContent;
