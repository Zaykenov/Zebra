import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import MainLayout from "@layouts/MainLayout";
import { useRouter } from "next/router";
import SupplyForm from "@modules/SupplyForm/SupplyForm";
import { getSupply } from "@api/supplies";
import TerminalLayout from "@layouts/TerminalLayout/TerminalLayout";
import TerminalSupplyForm from "@modules/TerminalSupplyForm";

const EditSupplyFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getSupply(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <TerminalLayout>
      <MainLayout title="Редактирование поставки" backBtn={true}>
        <div className="p-5 flex flex-col items-center justify-center mb-16">
          <TerminalSupplyForm data={data} />
        </div>
      </MainLayout>
    </TerminalLayout>
  );
};

export default EditSupplyFormPage;
