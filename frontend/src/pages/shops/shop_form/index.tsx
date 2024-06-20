import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import ShopWizardForm from "@modules/ShopWizardForm";

const ShopFormPage: NextPage = () => {
  return (
    <PageLayout>
      <MainLayout title="Добавление заведения" backBtn={true}>
        <div className="h-full overflow-visible">
          <div className="p-5 bg-gray-100">
            <ShopWizardForm />
          </div>
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default ShopFormPage;
