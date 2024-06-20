import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import ReceiptsContent from "../../components/__pageContent/ReceiptsContent";

const ReceiptsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true, hasDatePicker: true }}>
      <ReceiptsContent />
    </PageLayout>
  );
};

export default ReceiptsPage;
