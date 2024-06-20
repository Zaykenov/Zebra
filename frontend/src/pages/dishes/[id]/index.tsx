import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import DishForm from "@modules/DishForm";
import { useRouter } from "next/router";
import { getDish } from "@api/dishes";

const DishFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getDish(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);
  return (
    <PageLayout>
      <MainLayout title="Редактирование тех. карты" backBtn={true}>
        <div className="p-5">
          <DishForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default DishFormPage;
