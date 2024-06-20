import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import SupplyContent from "../../components/__pageContent/SupplyContent";

const SuppliesPage: NextPage = () => {
  return (
    <PageLayout
      defaultFilters={{
        hasPagination: true,
        hasDatePicker: true,
      }}
    >
      <SupplyContent />
    </PageLayout>
  );
};

export default SuppliesPage;
