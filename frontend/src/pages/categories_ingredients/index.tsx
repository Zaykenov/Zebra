import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllIngredientCategories } from "@api/ingredient-category";

const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
];

const IngredientCategoriesPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllIngredientCategories().then((res) => {
      const data = res.data.map((category: any) => ({
        name: category.name,
        id: category.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout
        title="Категории ингредиентов"
        addHref="/categories_ingredients/form"
      >
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default IngredientCategoriesPage;
