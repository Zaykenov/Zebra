import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import MenuItemForm from "@modules/MenuItemForm";
import { useRouter } from "next/router";
import { getMenuItem } from "@api/menu-items";
import AccountForm from "@modules/AccountForm";
import { getAccount } from "@api/accounts";

const EditProductFormPage: NextPage = () => {
  const router = useRouter();

  const [data, setData] = useState(null);

  useEffect(() => {
    const id = router.query.id;
    if (!id) return;
    getAccount(id as string).then((res) => {
      setData(res.data);
    });
  }, [router]);

  return (
    <PageLayout>
      <MainLayout title="Редактирование товара" backBtn={true}>
        <div className="p-5">
          <AccountForm data={data} />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default EditProductFormPage;
