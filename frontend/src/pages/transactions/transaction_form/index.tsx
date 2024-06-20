import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import TransactionForm from "@modules/TransactionForm";

const ProductFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление транзакции" backBtn={true}>
        <div className="p-5">
          <TransactionForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
