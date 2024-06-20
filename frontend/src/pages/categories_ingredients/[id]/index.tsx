import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import { useRouter } from "next/router";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import { getIngredientCategory } from "@api/ingredient-category";
import IngredientCategoryForm from "@modules/IngredientCategoryForm";

const EditIngredientCategoryFormPage: NextPage = () => {
  const router = useRouter();
  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getIngredientCategory(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование категории ингредиентов" backBtn={true}>
        <div className="p-5">
          <IngredientCategoryForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditIngredientCategoryFormPage;
