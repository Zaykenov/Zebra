import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { formatNumber } from "@utils/formatNumber";
import { getAllMenuItems, getMasterMenuItems } from "@api/menu-items";
import { useFilter } from "@context/index";
import { IMenuItem } from "./types";
import { columns } from "./constants";
import useMasterRole from "@hooks/useMasterRole";

const MenuContent: FC = () => {
  const isMaster = useMasterRole();

  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState<Array<IMenuItem>>([]);

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
      margin: `${item.margin.toFixed(2)}%`,
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
      getMasterMenuItems(queryOptions).then(handleGetAll);
    } else {
      getAllMenuItems(queryOptions).then(handleGetAll);
    }
  }, [queryOptions, isMaster]);

  return (
    <MainLayout
      title="Товары"
      addHref="/menu/product_form"
      searchFilter
      pagination
      filterOptions={[FilterOption.PRODUCTS_CATEGORY, FilterOption.SHOP]}
    >
      {tableData && <Table columns={columns} data={tableData} />}
    </MainLayout>
  );
};

export default MenuContent;
