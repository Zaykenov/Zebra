import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import ProductCategoryForm from "@modules/ProductCategoryForm";

const ProductCategoryFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление категории товара" backBtn={true}>
        <div className="p-5">
          <ProductCategoryForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductCategoryFormPage;
