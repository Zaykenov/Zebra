import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import IngredientsContent from "../../components/__pageContent/IngredientsContent";

const IngredientsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <IngredientsContent />
    </PageLayout>
  );
};

export default IngredientsPage;
