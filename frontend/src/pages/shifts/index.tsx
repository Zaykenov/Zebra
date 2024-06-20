import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import ShiftsContent from "../../components/__pageContent/ShiftsContent";

const ShiftsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true, hasDatePicker: true }}>
      <ShiftsContent />
    </PageLayout>
  );
};

export default ShiftsPage;
