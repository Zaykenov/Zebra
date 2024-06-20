import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllProductCategories } from "@api/product-categories";

const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
];

const ProductCategoriesPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllProductCategories().then((res) => {
      const data = res.data.map((category: any) => ({
        name: category.name,
        id: category.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Категории товаров" addHref="/categories_products/form">
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default ProductCategoriesPage;
