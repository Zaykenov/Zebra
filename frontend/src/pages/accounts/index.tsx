import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllAccounts } from "@api/accounts";
import { formatNumber } from "@utils/formatNumber";

const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Тип",
    accessor: "type",
  },
  {
    Header: "Баланс",
    accessor: "start_balance",
  },
];

const AccountsPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllAccounts().then((res) => {
      const data = res.data.map((item: any) => ({
        name: item.name,
        type: item.type,
        start_balance: formatNumber(item.start_balance, true, true),
        id: item.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Счета" addHref="/accounts/account_form">
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default AccountsPage;
