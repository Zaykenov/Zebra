import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import SupplierForm from "@modules/SupplierForm";

const ProductFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление поставщика" backBtn={true}>
        <div className="p-5">
          <SupplierForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
