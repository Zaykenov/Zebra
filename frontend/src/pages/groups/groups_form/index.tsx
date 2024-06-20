import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import GroupForm from "@modules/GroupForm";

const ProductFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление группы продуктов" backBtn={true}>
        <div className="p-5">
          <GroupForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
