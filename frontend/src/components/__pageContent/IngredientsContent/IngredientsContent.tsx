import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { formatNumber } from "@utils/formatNumber";
import { getAllIngredients, getMasterIngredients } from "@api/ingredient";
import { useFilter } from "@context/index";
import { columns } from "./constants";
import useMasterRole from "@hooks/useMasterRole";

const IngredientsContent: FC = () => {
  const isMaster = useMasterRole();

  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const handleGetAll = useCallback((res: any) => {
    changeTotalPages(res.data.totalPages);
    const data = res.data.data.map((ingredient: any) => ({
      name: ingredient.name,
      shop_name: ingredient.shop_name,
      category: ingredient.category,
      measure: ingredient.measure,
      cost: formatNumber(ingredient.cost, true, true),
      id: ingredient.id,
    }));
    setTableData(data);
    changeTotalResults(data.length);
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length || isMaster === null)
      return;
    if (isMaster) {
      getMasterIngredients(queryOptions).then(handleGetAll);
    } else {
      getAllIngredients(queryOptions).then(handleGetAll);
    }
  }, [queryOptions, isMaster]);

  return (
    <MainLayout
      title="Ингредиенты"
      addHref="/ingredients/ingredient_form"
      searchFilter
      pagination
      filterOptions={[FilterOption.INGREDIENTS_CATEGORY, FilterOption.SHOP]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          isDetailedConfirmation={true}
        />
      )}
    </MainLayout>
  );
};

export default IngredientsContent;
