import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import SupplyForm from "@modules/SupplyForm/SupplyForm";

const SupplyFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление поставки" backBtn={true}>
        <div className="p-5">
          <SupplyForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default SupplyFormPage;
