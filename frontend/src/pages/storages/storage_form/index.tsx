import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import StorageForm from "@modules/StorageForm";

const StorageFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление склада" backBtn={true}>
        <div className="p-5">
          <StorageForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default StorageFormPage;
