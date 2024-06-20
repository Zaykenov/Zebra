import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import TransferContent from "../../components/__pageContent/TransferContent";

const TransferPage: NextPage = () => {
  return (
    <PageLayout
      defaultFilters={{
        hasDatePicker: true,
        hasPagination: true,
      }}
    >
      <TransferContent />
    </PageLayout>
  );
};

export default TransferPage;
