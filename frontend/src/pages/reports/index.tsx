import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import ReportsContent from "../../components/__pageContent/ReportsContent";

const ReportsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasDatePicker: true, hasPagination: true }}>
      <ReportsContent />
    </PageLayout>
  );
};

export default ReportsPage;
