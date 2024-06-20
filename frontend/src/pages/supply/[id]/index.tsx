import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import { useRouter } from "next/router";
import SupplyForm from "@modules/SupplyForm/SupplyForm";
import { getSupply } from "@api/supplies";

const EditSupplyFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getSupply(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование поставки" backBtn={true}>
        <div className="p-5">
          <SupplyForm data={data} isEdit />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditSupplyFormPage;
