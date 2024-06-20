import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import TransactionsContent from "../../components/__pageContent/TransactionsContent";

const TransactionsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true, hasDatePicker: true }}>
      <TransactionsContent />
    </PageLayout>
  );
};

export default TransactionsPage;
