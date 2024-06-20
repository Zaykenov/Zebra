import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import WasteForm from "@modules/WasteForm";

const WasteFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление списания" backBtn={true}>
        <div className="p-5">
          <WasteForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default WasteFormPage;
