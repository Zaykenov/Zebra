import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import AbcContent from "../../components/__pageContent/AbcContent";

const LeftPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true, hasDatePicker: true }}>
      <AbcContent />
    </PageLayout>
  );
};

export default LeftPage;
