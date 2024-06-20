import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import { getExcelFile } from "@api/excel";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeadersAndAccessor from "@utils/extractHeaderAndAccessor";
import { getAllSklads } from "@api/sklad";
import { QueryOptions } from "@api/index";
import { getAllTrafficReports } from "@api/reports";
import filterProperties from "@utils/filterProperties";
import { ReportItem } from "./types";
import { columns } from "./constants";
import { useFilter } from "@context/index";

const ReportsContent: FC = () => {
  const {
    queryOptions,
    changeTotalPages,
    changeTotalResults,
    handleFilterChange,
  } = useFilter();

  const [tableData, setTableData] = useState([]);

  const [totalSum, setTotalSum] = useState({ initialSum: 0, finalSum: 0 });

  const { headers, accessors } = extractHeadersAndAccessor(columns(totalSum));

  const [excelData, setExcelData] = useState<Partial<ReportItem>[]>([]);

  useEffect(() => {
    getAllSklads().then((res) => {
      const firstSklad = res.data[0];
      firstSklad && handleFilterChange({ [QueryOptions.SKLAD]: firstSklad.id });
    });
  }, []);

  const handleGetAll = useCallback((res: any) => {
    const transformItem = (item: any) => ({
      ...item,
      item_name: item.item_name,
      type: item.type,
      initial_ostatki: item.initial_ostatki,
      initial_netCost: item.initial_netCost,
      initial_sum: item.initial_sum,
      income: item.income,
      consumption: item.consumption,
      final_ostatki: item.final_ostatki,
      final_netCost: item.final_netCost,
      final_sum: item.final_sum,
    });

    let data;
    if (res.data.data) {
      data = res.data.data.traffic_reports.map(transformItem);
      setTotalSum({
        initialSum: res.data.data.initial_sum,
        finalSum: res.data.data.final_sum,
      });
    } else {
      data = res.data.traffic_reports.map(transformItem);
      setTotalSum({
        initialSum: res.data.initial_sum,
        finalSum: res.data.final_sum,
      });
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllTrafficReports(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalResults(data.length);
      changeTotalPages(res.data.totalPages);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllTrafficReports(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties<ReportItem>(data, accessors);
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Отчет по движению"
      dateFilter
      searchFilter
      pagination
      filterOptions={[
        FilterOption.ITEMS_CATEGORY,
        FilterOption.SKLAD,
        FilterOption.ITEM_TYPE,
      ]}
      excelDownloadButton={() =>
        getExcelFile(`Отчет по движению`, headers, excelData)
      }
    >
      {tableData && (
        <Table
          hasFooter
          columns={columns(totalSum)}
          data={tableData}
          editable={false}
        />
      )}
      {tableData.length === 0 && (
        <div className="py-52 font-light flex items-center justify-center text-center">
          По этим фильтрам нет движения на складе. Выберите другой отчетный
          период или измените фильтры.
        </div>
      )}
    </MainLayout>
  );
};

export default ReportsContent;
