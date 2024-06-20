import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import PaymentsContent from "../../components/__pageContent/PaymentsContent";

const SuppliesPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true, hasDatePicker: true }}>
      <PaymentsContent />
    </PageLayout>
  );
};

export default SuppliesPage;
