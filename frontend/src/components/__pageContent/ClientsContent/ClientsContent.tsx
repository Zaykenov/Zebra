import React, { FC, useCallback, useEffect, useState } from "react";
import { getExcelFile } from "@api/excel";
import { getDateAndTime, getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeaderAndAccessor from "@utils/extractHeaderAndAccessor";
import { formatNumber } from "@utils/formatNumber";
import clientsService from "@services/clientsService";
import { columns, filterProperties } from "./constants";
import SubRow from "./SubRow";
import { useFilter } from "@context/index";
import { ClientData } from "./types";

const ClientsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState<ClientData[]>([]);

  const { headers, accessors } = extractHeaderAndAccessor(columns);
  const [excelData, setExcelData] = useState<Partial<ClientData>[]>([]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  const handleGetAll = useCallback((res: any) => {
    const transformItem = (user: any) => ({
      id: user.id,
      email: user.email,
      firstName: user.name,
      birthDate: getDateAndTime(user.birth_date),
      registrationDate: getDateAndTime(user.reg_date),
      discount: user.discount + "%",
      zebraCoinBalance: formatNumber(user.zebra_coin_balance, true, true),
      removeDate: getDateAndTime(user.remove_date),
      status: user.status,
      feedbacks: user.feedbacks,
      checks: user.checks,
    });

    let data;
    if (res.data.data) {
      data = res.data.data.map(transformItem);
    } else {
      data = res.data.map(transformItem);
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    clientsService.getAllClients(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalPages(res.data.totalPages);
      changeTotalResults(data.length);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    clientsService.getAllClients(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties(data, accessors);
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Пользователи"
      excelDownloadButton={() =>
        getExcelFile(`Клиенты ${getDateString()}`, headers, excelData)
      }
      pagination
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          editable={false}
          details
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default ClientsContent;
