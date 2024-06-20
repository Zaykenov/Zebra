import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import AccountForm from "@modules/AccountForm";

const ProductFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление счета" backBtn={true}>
        <div className="p-5">
          <AccountForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
