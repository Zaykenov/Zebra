import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import WorkerForm from "@modules/WorkerForm";
import { useRouter } from "next/router";
import { getWorker } from "@api/workers";

const EditWorkerFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getWorker(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование сотрудника" backBtn={true}>
        <div className="p-5">
          <WorkerForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditWorkerFormPage;
