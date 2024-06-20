import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import { Column } from "react-table";
import { getAllWorkers } from "@api/workers";
import { getAllShops } from "@api/shops";

const columns: Column[] = [
  {
    Header: "#",
    accessor: "id",
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

const ShopsPage: NextPage = () => {
  const [tableData, setTableData] = useState(null);

  useEffect(() => {
    getAllShops().then((res) => {
      const data = res.data.map((shop: any) => ({
        id: shop.id,
        name: shop.name,
        address: shop.address,
      }));
      setTableData(data);
    });
  }, []);

  return (
    <PageLayout>
      <MainLayout title="Заведения" addHref="/shops/shop_form">
        {tableData && <Table columns={columns} data={tableData} editable={false}/>}
      </MainLayout>
    </PageLayout>
  );
};

export default ShopsPage;
