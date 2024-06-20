import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import WorkerForm from "@modules/WorkerForm";

const WorkerFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление сотрудника" backBtn={true}>
        <div className="p-5">
          <WorkerForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default WorkerFormPage;
