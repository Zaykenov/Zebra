import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MenuContent from "../../components/__pageContent/MenuContent";

const MenuPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <MenuContent />
    </PageLayout>
  );
};
export default MenuPage;
