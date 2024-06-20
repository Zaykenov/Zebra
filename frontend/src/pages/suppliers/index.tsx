import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllSuppliers } from "@api/suppliers";

const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Адрес",
    accessor: "address",
  },
  {
    Header: "Телефон",
    accessor: "phone",
  },
  {
    Header: "Комментарий",
    accessor: "comment",
  },
];

const SuppliersPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllSuppliers().then((res) => {
      const data = res.data.map((item: any) => ({
        name: item.name,
        address: item.address,
        phone: item.phone,
        comment: item.comment,
        id: item.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Поставщики" addHref="/suppliers/supplier_form">
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default SuppliersPage;
