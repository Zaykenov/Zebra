import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import InventoryForm from "@modules/InventoryForm";

const InventoryFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление инвентаризации" backBtn={true}>
        <div className="p-5">
          <InventoryForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default InventoryFormPage;
