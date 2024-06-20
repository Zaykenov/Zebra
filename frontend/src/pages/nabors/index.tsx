import React from "react";
import PageLayout from "@layouts/PageLayout";
import { NextPage } from "next";
import NaborsContent from "../../components/__pageContent/NaborsContent";

const NaborsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <NaborsContent />
    </PageLayout>
  );
};

export default NaborsPage;
