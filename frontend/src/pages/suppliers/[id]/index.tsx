import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import { useRouter } from "next/router";
import { getSupplier } from "@api/suppliers";
import SupplierForm from "@modules/SupplierForm";

const EditSupplierFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getSupplier(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование товара" backBtn={true}>
        <div className="p-5">
          <SupplierForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditSupplierFormPage;
