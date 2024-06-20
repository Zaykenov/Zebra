import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import CalculationsContent from "../../components/__pageContent/CalculationsContent";

const CalculationsPage: NextPage = () => {
  return (
    <PageLayout
      defaultFilters={{
        hasPagination: true,
      }}
    >
      <CalculationsContent />
    </PageLayout>
  );
};

export default CalculationsPage;
