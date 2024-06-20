import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import MenuItemForm from "@modules/MenuItemForm";
import { useRouter } from "next/router";
import { getIngredient } from "@api/ingredient";
import IngredientForm from "@modules/IngredientForm/IngredientForm";

const EditIngredientFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getIngredient(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование ингредиента" backBtn={true}>
        <div className="p-5">
          <IngredientForm data={data} isEdit />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditIngredientFormPage;
