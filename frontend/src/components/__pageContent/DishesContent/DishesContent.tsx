import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { formatNumber } from "@utils/formatNumber";
import { getAllDishes, getMasterDishes } from "@api/dishes";
import SubRow from "./SubRow";
import { useFilter } from "@context/index";
import { columns } from "./constants";
import useMasterRole from "@hooks/useMasterRole";

const DishesContent: FC = () => {
  const isMaster = useMasterRole();

  const { queryOptions, changeTotalResults, changeTotalPages } = useFilter();

  const [tableData, setTableData] = useState([]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  const handleGetAll = useCallback((res: any) => {
    changeTotalPages(res.data.totalPages);
    const data = res.data.data.map((item: any) => ({
      name: item.name,
      shop_name: item.shop_name,
      category: item.category,
      tax: item.tax,
      cost: formatNumber(item.cost, true, true),
      measure: item.measure,
      price: formatNumber(item.price, true, true),
      margin: `${Math.round(item.margin)}%`,
      profit: formatNumber(item.profit, true, true),
      id: item.id,
    }));
    setTableData(data);
    changeTotalResults(data.length);
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length || isMaster === null)
      return;
    if (isMaster) {
      getMasterDishes(queryOptions).then(handleGetAll);
    } else {
      getAllDishes(queryOptions).then(handleGetAll);
    }
  }, [queryOptions, isMaster]);

  return (
    <MainLayout
      title="Тех. карты"
      addHref={isMaster ? "/dishes/dish_form" : ""}
      searchFilter
      pagination
      filterOptions={[FilterOption.PRODUCTS_CATEGORY, FilterOption.SHOP]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          details={true}
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default DishesContent;
