import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllInventoryGroups, InventoryGroupData } from "@api/groups";

const columns: Column[] = [
  {
    Header: "Название",
    accessor: "name",
  },
  {
    Header: "Склад",
    accessor: "sklad_name",
  },
  {
    Header: "Ед. измерения",
    accessor: "measure",
  },
  {
    Header: "Тип",
    accessor: "type",
    Cell: ({ value }) => (
      <span>{value === "ingredient" ? "Ингредиенты" : "Товары"}</span>
    ),
  },
];

const GroupsPage: NextPage = () => {
  const [tableData, setTableData] = useState<InventoryGroupData[]>([]);

  useEffect(() => {
    getAllInventoryGroups().then((res) => {
      setTableData(res.data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout
        title="Группы товаров и ингредиентов для инвентаризации"
        addHref="/groups/groups_form"
      >
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};
export default GroupsPage;
