import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import MenuItemForm from "@modules/MenuItemForm";

const ProductFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление товара" backBtn={true}>
        <div className="p-5">
          <MenuItemForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
