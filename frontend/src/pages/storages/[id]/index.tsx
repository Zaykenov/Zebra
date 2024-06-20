import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import StorageForm from "@modules/StorageForm";
import { useRouter } from "next/router";
import { getSklad } from "@api/sklad";

const EditStorageFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getSklad(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование склада" backBtn={true}>
        <div className="p-5">
          <StorageForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditStorageFormPage;
