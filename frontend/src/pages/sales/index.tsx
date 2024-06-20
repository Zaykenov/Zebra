import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import SalesContent from "../../components/__pageContent/SalesContent";

const SalesPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasDatePicker: true }}>
      <SalesContent />
    </PageLayout>
  );
};

export default SalesPage;
