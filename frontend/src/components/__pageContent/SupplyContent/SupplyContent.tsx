import React, { FC, useCallback, useEffect, useState } from "react";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { useFilter } from "@context/filter.context";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import { getAllSupplies } from "@api/supplies";
import SubRow from "./SubRow";
import {
  columns,
  filterProperties,
  filterTopItems,
  headerRow,
  propertyNames,
} from "./constants";
import { SupplyItem } from "./types";
import { FilterOption } from "@layouts/MainLayout/types";

const SupplyContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);
  const [totalSum, setTotalSum] = useState<number>(0);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  const [excelData, setExcelData] = useState<Partial<SupplyItem>[]>([]);

  const handleGetAll = useCallback((res: any) => {
    const transformSupply = (supply: any) => ({
      ...supply,
      dealer: supply.dealer,
      sklad: supply.sklad,
      schet: supply.schet,
      isDeleted: supply.deleted,
      time: dateToString(supply.time, false),
      sum: formatNumber(supply.sum, true, true),
      intSum: supply.sum,
      items: filterTopItems(supply.items.map((item: any) => item.name)).join(
        ", ",
      ),
      id: supply.id,
    });

    let data;
    if (res.data.data) {
      data = res.data.data.postavka.map(transformSupply);
    } else {
      data = res.data.postavka.map(transformSupply);
    }
    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    console.log(queryOptions);
    getAllSupplies(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTotalSum(res.data.data.sum);
      changeTotalPages(res.data.totalPages);
      changeTotalResults(data.length);
      setTableData(data);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllSupplies(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties(data, propertyNames);
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Поставки"
      addHref="/supply/supply_form"
      dateFilter
      searchFilter
      pagination
      filterOptions={[
        FilterOption.ITEMS_CATEGORY,
        FilterOption.DEALER,
        FilterOption.SKLAD,
        FilterOption.SCHET,
      ]}
      excelDownloadButton={() =>
        getExcelFile(`Поставки ${getDateString()}`, headerRow, excelData)
      }
    >
      {tableData && (
        <Table
          columns={columns(totalSum)}
          data={tableData}
          details
          hasFooter
          isRowEditable={(row) => !row.original.isDeleted}
          isRowDeletable={(row) => !row.original.isDeleted}
          customRowStyle={(row) =>
            row.original.isDeleted ? "bg-red-200/80" : "bg-white-100/80"
          }
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default SupplyContent;
