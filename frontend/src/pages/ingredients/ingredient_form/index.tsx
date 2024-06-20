import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import IngredientForm from "@modules/IngredientForm/IngredientForm";

const IngredientFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление ингредиента" backBtn={true}>
        <div className="p-5">
          <IngredientForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default IngredientFormPage;
