import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import IngredientCategoryForm from "@modules/IngredientCategoryForm";

const IngredientCategoryFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление категории ингредиента" backBtn={true}>
        <div className="p-5">
          <IngredientCategoryForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default IngredientCategoryFormPage;
