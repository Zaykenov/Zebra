import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import { getAllWastes } from "@api/wastes";
import SubRow from "./SubRow";
import { WasteItem } from "./types";
import {
  columns,
  filterProperties,
  handleConfirmWasteRemoval,
  handleRejectWasteRemoval,
  headerRow,
  propertyNames,
  statusText,
  statusToColor,
} from "./constants";
import { useFilter } from "@context/filter.context";

const WasteContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);
  const [totalSum, setTotalSum] = useState<number>(0);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  const [excelData, setExcelData] = useState<Partial<WasteItem>[]>([]);

  const handleGetAll = useCallback((res: any) => {
    const transformWasteItem = (item: any) => {
      const items = item.items
        .slice(0, 3)
        .map((removedItem: any) => removedItem.name)
        .join(", ");
      return {
        id: item.id,
        time: dateToString(item.time, false),
        sklad: item.sklad,
        items,
        cost: formatNumber(item.cost, true, true),
        reason: item.reason,
        comment: item.comment,
        status: statusText[item.status as "opened" | "closed" | "rejected"],
      };
    };

    let data;
    if (res.data.data) {
      data = res.data.data.remove_from_sklad.map(transformWasteItem);
    } else {
      data = res.data.remove_from_sklad.map(transformWasteItem);
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllWastes(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalResults(data.length);
      changeTotalPages(res.data.totalPages);
      setTotalSum(res.data.data.sum);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllWastes(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties(data, propertyNames);
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Списания"
      addHref="/waste/waste_form"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.ITEMS_CATEGORY, FilterOption.SKLAD]}
      excelDownloadButton={() =>
        getExcelFile(`Списания ${getDateString()}`, headerRow, excelData)
      }
    >
      {tableData && (
        <Table
          columns={columns(totalSum)}
          data={tableData}
          editable={true}
          details={true}
          hasFooter
          renderRowSubComponent={renderSubComponent}
          customRowStyle={(row) => statusToColor(row)}
          customEditBtn={(row) =>
            row.original.status === "Открыто" ? (
              <button
                onClick={() => handleConfirmWasteRemoval(row)}
                className="text-indigo-500 hover:text-indigo-600 hover:underline"
              >
                Принять
              </button>
            ) : (
              false
            )
          }
          customDeleteBtn={(row) =>
            row.original.status === "Открыто" ? (
              <button
                onClick={() => handleRejectWasteRemoval(row)}
                className="text-red-500 hover:text-red-600 hover:underline"
              >
                Отклонить
              </button>
            ) : (
              false
            )
          }
        />
      )}
    </MainLayout>
  );
};

export default WasteContent;
