import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import TransferForm from "@modules/TransferForm";
import { useRouter } from "next/router";
import { getTransferById } from "@api/transfers";

const TransferFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getTransferById(id as string).then((res) => {
      setData(res.data);
      setLoading(false);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Перемещение товаров" backBtn={true}>
        <div className="p-5">{!loading && <TransferForm data={data} />}</div>
      </MainLayout>
    </PageLayout>
  );
};

export default TransferFormPage;
