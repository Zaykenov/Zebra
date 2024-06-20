import React, {useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import { Column } from "react-table";
import Table from "@common/Table";
import { formatNumber } from "@utils/formatNumber";
import { getWorkersStats } from "@api/stats";
import TableLoader from "@common/TableLoader/TableLoader";
import useFetching from "@hooks/useFetching";
import FetchingError from "@common/FetchingError/FetchingError";
import extractHeadersAndAccessor from "@utils/extractHeaderAndAccessor";

interface Cashier {
  name: string;
  revenue: number;
  profit: number;
  checkNum: number;
  avgCheck: number;
}

const columns: Column<any>[] = [
  {
    Header: "Официант",
    accessor: "name",
  },
  {
    Header: "Выручка",
    accessor: "revenue",
  },
  {
    Header: "Прибыль",
    accessor: "profit",
  },
  {
    Header: "Чеки",
    accessor: "checkNum",
  },
  {
    Header: "Средний чек",
    accessor: "avgCheck",
  },
];

const CashiersPage: NextPage = () => {
  const [tableData, setTableData] = useState<Cashier[]>([]);
  const {headers, accessors} = extractHeadersAndAccessor(columns)
  const {fetchData, isLoading, error} = useFetching(async()=>{
    const workerStats = (await getWorkersStats()).data
    const tableData = workerStats.map((worker: any) => ({
      name: worker.name,
      checkNum: worker.check_num,
      revenue: formatNumber(worker.revenue, true, true),
      profit: formatNumber(worker.profit, true, true),
      avgCheck: formatNumber(worker.avg_check, true, true),
    }))
    setTableData(tableData)
  })

  useEffect(() => {
    fetchData()
  }, []);

  const showError = error !== "";
  const showLoader = isLoading && !showError;
  const showTable = !isLoading && !showError;

  return (
    <PageLayout>
      <MainLayout title="Сотрудники">
        {showError && <FetchingError errorMessage={error} />}
        {showLoader && <TableLoader headerRowNames={headers} rowCount={20} />}
        {showTable && (
          <Table columns={columns} data={tableData} editable={false} />
        )}
      </MainLayout>
    </PageLayout>
  );
};

export default CashiersPage;