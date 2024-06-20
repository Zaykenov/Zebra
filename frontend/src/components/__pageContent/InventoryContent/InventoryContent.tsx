import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeadersAndAccessor from "@utils/extractHeaderAndAccessor";
import { dateToString } from "@api/check";
import { getAllInventory } from "@api/inventory";
import filterProperties from "@utils/filterProperties";
import { columns } from "./constants";
import { InventoryItem } from "./types";
import { useFilter } from "@context/index";

const InventoryContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const { headers, accessors } = extractHeadersAndAccessor(columns);

  const [excelData, setExcelData] = useState<Partial<InventoryItem>[]>([]);

  const handleGetAll = useCallback((res: any) => {
    const transformInventoryItem = (item: any) => ({
      sklad: item.sklad,
      time:
        new Date(item.time).getFullYear() < 2022
          ? "-"
          : dateToString(item.time, false, false),
      type: item.type === "partial" ? "Частичная" : "Полная",
      result: item.result,
      status: item.status === "opened" ? "На редактировании" : "Проведенная",
      id: item.id,
    });

    let data;
    if (res.data.data) {
      data = res.data.data.map(transformInventoryItem);
    } else {
      data = res.data.map(transformInventoryItem);
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllInventory(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalResults(data.length);
      changeTotalPages(res.data.totalPages);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllInventory(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties<InventoryItem>(
        data,
        accessors,
      );
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Инвентаризации"
      addHref="/inventory/inventory_form"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.SKLAD]}
      excelDownloadButton={() =>
        getExcelFile(`Инвентаризации ${getDateString()}`, headers, excelData)
      }
    >
      {tableData && <Table columns={columns} data={tableData} />}
    </MainLayout>
  );
};

export default InventoryContent;
