import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import { getProductCategory } from "@api/product-categories";
import { useRouter } from "next/router";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import ProductCategoryForm from "@modules/ProductCategoryForm";

const EditProductCategoriesFormPage: NextPage = () => {
  const router = useRouter();
  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getProductCategory(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование категории товаров" backBtn={true}>
        <div className="p-5">
          <ProductCategoryForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditProductCategoriesFormPage;
