import React from "react";
import PageLayout from "@layouts/PageLayout";
import { NextPage } from "next";
import DishesContent from "../../components/__pageContent/DishesContent/DishesContent";

const DishesPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <DishesContent />
    </PageLayout>
  );
};

export default DishesPage;
