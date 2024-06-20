import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllSklads } from "@api/sklad";

const columns: Column[] = [
  {
    Header: "#",
    accessor: "num",
  },
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Адрес",
    accessor: "address",
  },
];

const StoragesPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllSklads().then((res) => {
      const data = res.data.map((sklad: any, idx: number) => ({
        num: idx + 1,
        name: sklad.name,
        address: sklad.address,
        id: sklad.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Склады" addHref="/storages/storage_form">
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default StoragesPage;
