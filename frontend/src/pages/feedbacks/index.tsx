import React from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import FeedbacksContent from "../../components/__pageContent/FeedbacksContent";

const FeedbacksPage: NextPage = () => {
  return (
    <PageLayout defaultFilters={{ hasPagination: true }}>
      <FeedbacksContent />
    </PageLayout>
  );
};

export default FeedbacksPage;
