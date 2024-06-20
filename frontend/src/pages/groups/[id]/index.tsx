import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import GroupForm from "@modules/GroupForm";
import { useRouter } from "next/router";
import { getInventoryGroup } from "@api/groups";

const ProductFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    if (!router) return;
    const id = router.query.id;
    if (!id) return;
    getInventoryGroup(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование группы продуктов" backBtn={true}>
        <div className="p-5">
          <GroupForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ProductFormPage;
