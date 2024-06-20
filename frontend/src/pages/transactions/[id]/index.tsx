import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import { useRouter } from "next/router";
import TransactionForm from "@modules/TransactionForm";
import { getTransaction } from "@api/transactions";

const EditTransactionFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getTransaction(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование транзакции" backBtn={true}>
        <div className="p-5">
          <TransactionForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditTransactionFormPage;
