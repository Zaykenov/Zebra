import React from "react";
import { NextPage } from "next";
import MainLayout from "@layouts/MainLayout";
import SupplyForm from "@modules/SupplyForm/SupplyForm";
import TerminalLayout from "@layouts/TerminalLayout";

const SupplyFormPage: NextPage = () => {
  return (
    <TerminalLayout>
      <MainLayout title="Добавление поставки" backBtn={true}>
        <div className="p-5">
          <SupplyForm />
        </div>
      </MainLayout>
    </TerminalLayout>
  );
};

export default SupplyFormPage;