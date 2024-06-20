import React, { FC, useCallback, useEffect, useState } from "react";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import FetchingError from "@common/FetchingError";
import TableLoader from "@common/TableLoader";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeadersAndAccessor from "@utils/extractHeaderAndAccessor";
import useFetching from "@hooks/useFetching";
import { getAllOstatki } from "@api/ostatki";
import filterProperties from "@utils/filterProperties";
import { formatNumber } from "@utils/formatNumber";
import { useFilter } from "@context/filter.context";
import { columns } from "./constants";
import { CalculationItem } from "./types";
import { FilterOption } from "@layouts/MainLayout/types";

const CalculationsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const { headers, accessors } = extractHeadersAndAccessor(columns);

  const [excelData, setExcelData] = useState<Partial<CalculationItem>[]>([]);

  const { fetchData, isLoading, error } = useFetching(
    async (newQueryOptions) => {
      if (!newQueryOptions || !Object.keys(newQueryOptions).length) return;
      const res = await getAllOstatki(newQueryOptions);
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalPages(res.data.totalPages);
      changeTotalResults(data.length);

      const currentQueryOptions = { ...newQueryOptions };
      delete currentQueryOptions.page;
      const resFiltered = await getAllOstatki(currentQueryOptions);
      const filteredData = handleGetAll(resFiltered);
      const filteredProperties = filterProperties<CalculationItem>(
        filteredData,
        accessors,
      );
      setExcelData(filteredProperties);
    },
  );

  const handleGetAll = useCallback((res: any) => {
    const transformItem = (item: any) => ({
      name: item.name,
      skladName: item.skladName,
      type: item.type === "ingredient" ? "Ингредиент" : "Товар",
      measurement: item.measurement,
      category: item.category,
      cost: formatNumber(item.cost, true, true),
      intCost: item.cost,
      quantity: item.quantity.toFixed(2),
      intQuantity: parseInt(item.quantity.toFixed(2)),
      sum: formatNumber(item.sum, true, true),
      intSum: item.sum,
      id: item.id,
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
    if (!queryOptions) return;
    fetchData(queryOptions).then();
  }, [queryOptions]);

  return (
    <MainLayout
      title="Остатки"
      excelDownloadButton={() =>
        getExcelFile(`Остатки ${getDateString()}`, headers, excelData)
      }
      searchFilter
      pagination
      filterOptions={[
        FilterOption.ITEMS_CATEGORY,
        FilterOption.SKLAD,
        FilterOption.ITEM_TYPE,
      ]}
    >
      {!!error && <FetchingError errorMessage={error} />}
      {isLoading && <TableLoader headerRowNames={headers} rowCount={20} />}
      {!isLoading && !error && (
        <>
          <Table columns={columns} data={tableData} editable={false} />
        </>
      )}
    </MainLayout>
  );
};

export default CalculationsContent;
