import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import DishForm from "@modules/DishForm";

const DishFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление тех. карты" backBtn={true}>
        <div className="p-5">
          <DishForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default DishFormPage;
