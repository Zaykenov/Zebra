import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import FetchingError from "@common/FetchingError";
import TableLoader from "@common/TableLoader";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeaderAndAccessor from "@utils/extractHeaderAndAccessor";
import useFetching from "@hooks/useFetching";
import { getAbcAnalysis } from "@api/abc";
import filterProperties from "@utils/filterProperties";
import { formatNumber } from "@utils/formatNumber";
import { columns, renderCellStyle } from "./constants";
import { useFilter } from "@context/index";
import { AbcItem } from "./types";

const AbcContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const { headers, accessors } = extractHeaderAndAccessor(columns);
  const [excelData, setExcelData] = useState<Partial<AbcItem>[]>([]);

  const { fetchData, isLoading, error } = useFetching(
    async (newQueryOptions) => {
      const res = await getAbcAnalysis(newQueryOptions);
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalResults(data.length);
      changeTotalPages(res.data.totalPages);

      const currentQueryOptions = { ...newQueryOptions };
      delete currentQueryOptions.page;
      const resFiltered = await getAbcAnalysis(currentQueryOptions);
      const filteredData = handleGetAll(resFiltered);
      const filteredProperties = filterProperties<AbcItem>(
        filteredData,
        accessors,
      );
      setExcelData(filteredProperties);
    },
  );

  const handleGetAll = useCallback((res: any) => {
    let count = 1;
    const transformItem = (item: any) => ({
      id: count++,
      product: item.item_name,
      productType: item.item_type,
      sales: item.sales + " шт.",
      sales_percent: item.sales_percent.toFixed(2) + "%",
      revenue: formatNumber(item.revenue, true, true),
      revenue_percent: item.revenue_percent.toFixed(2) + "%",
      profit: formatNumber(item.profit, true, true),
      profit_percent: item.profit_percent.toFixed(2) + "%",
      sales_letter: item.sales_letter,
      revenue_letter: item.revenue_letter,
      profit_letter: item.profit_letter,
    });

    let data;
    if (res.data.data) {
      data = res.data.data.map(transformItem);
    } else {
      data = res.data.map(transformItem);
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    fetchData(queryOptions).then();
  }, [queryOptions]);

  return (
    <MainLayout
      title="ABC анализ"
      dateFilter
      searchFilter
      filterOptions={[FilterOption.ITEMS_CATEGORY, FilterOption.SHOP]}
      excelDownloadButton={() =>
        getExcelFile(`ABC analys ${getDateString()}`, headers, excelData)
      }
    >
      {!!error && <FetchingError errorMessage={error} />}
      {isLoading && <TableLoader headerRowNames={headers} rowCount={20} />}
      {!isLoading && !error && (
        <>
          <Table
            columns={columns}
            data={tableData}
            editable={false}
            renderCellStyle={renderCellStyle}
          />
        </>
      )}
    </MainLayout>
  );
};

export default AbcContent;
