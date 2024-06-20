import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import TransferForm from "@modules/TransferForm";

const TransferFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Перемещение товаров" backBtn={true}>
        <div className="p-5">
          <TransferForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default TransferFormPage;
