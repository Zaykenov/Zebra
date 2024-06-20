import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import ClientsContent from "../../components/__pageContent/ClientsContent";

const ClientsPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <ClientsContent />
    </PageLayout>
  );
};

export default ClientsPage;
