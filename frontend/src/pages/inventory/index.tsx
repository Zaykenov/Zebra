import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import InventoryContent from "../../components/__pageContent/InventoryContent";

const InventoryPage: NextPage = () => {
  return (
    <PageLayout
      defaultFilters={{
        hasPagination: true,
        hasDatePicker: true,
      }}
    >
      <InventoryContent />
    </PageLayout>
  );
};

export default InventoryPage;
