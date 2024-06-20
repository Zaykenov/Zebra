import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import WasteContent from "../../components/__pageContent/WasteContent";

const WastesPage: NextPage = () => {
  return (
    <PageLayout
      defaultFilters={{
        hasPagination: true,
        hasDatePicker: true,
      }}
    >
      <WasteContent />
    </PageLayout>
  );
};

export default WastesPage;
