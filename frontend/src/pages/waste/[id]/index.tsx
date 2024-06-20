import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import WasteForm from "@modules/WasteForm";
import { useRouter } from "next/router";
import { getWasteById } from "@api/wastes";

const WasteFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getWasteById(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование списания" backBtn={true}>
        <div className="p-5">
          <WasteForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default WasteFormPage;
