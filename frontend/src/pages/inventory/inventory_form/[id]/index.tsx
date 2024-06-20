import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import MainLayout from "@layouts/MainLayout";
import PageLayout from "@layouts/PageLayout";
import InventoryForm from "@modules/InventoryForm";
import { getInventoryById } from "@api/inventory";
import { useRouter } from "next/router";

const InventoryFormPage: NextPage = () => {
  const router = useRouter();
  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    id &&
      getInventoryById(id as string).then((res) => {
        setData(res.data);
      });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout
        title="Редактирование параметров инвентаризации"
        backBtn={true}
      >
        <div className="p-5">
          <InventoryForm isEdit data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default InventoryFormPage;
