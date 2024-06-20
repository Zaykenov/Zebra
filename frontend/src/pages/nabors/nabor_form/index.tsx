import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import NaborForm from "@modules/NaborForm";

const NaborFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление набора" backBtn={true}>
        <div className="p-5">
          <NaborForm />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default NaborFormPage;
