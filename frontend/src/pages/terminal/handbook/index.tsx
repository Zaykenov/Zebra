import React from "react";
import { NextPage } from "next";
import TerminalLayout from "@layouts/TerminalLayout";

const HandbookPage: NextPage = () => {
  return (
    <TerminalLayout>
      <iframe
        title="Памятка"
        src="/assets/handbook.pdf#view=fitV"
        height="100%"
        width="100%"
      />
    </TerminalLayout>
  );
};

export default HandbookPage;
