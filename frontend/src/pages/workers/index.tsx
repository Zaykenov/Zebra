import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllWorkers } from "@api/workers";

const columns: Column[] = [
  {
    Header: "Имя",
    accessor: "name",
  },
  {
    Header: "Логин",
    accessor: "username",
  },
  {
    Header: "Доступ",
    accessor: "role",
  },
];

const WorkersPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllWorkers().then((res) => {
      const data = res.data.map((worker: any) => ({
        name: worker.name,
        username: worker.username,
        phone: worker.phone,
        role: worker.role,
        id: worker.id,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Сотрудники" addHref="/workers/worker_form">
        {tableData && <Table columns={columns} data={tableData} />}
      </MainLayout>
    </PageLayout>
  );
};

export default WorkersPage;
