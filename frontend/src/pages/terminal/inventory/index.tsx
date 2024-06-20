import React, { useEffect, useMemo, useState } from "react";
import { NextPage } from "next";
import TerminalLayout from "@layouts/TerminalLayout";
import InventoryForm from "@modules/InventoryForm";
import TerminalCabinetLayout from "@layouts/TerminalCabinetLayout";
import Table from "@common/Table/Table";
import { Column } from "react-table";
import { getAllInventory } from "@api/inventory";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import { QueryOptions } from "@api/index";

const columns: Column[] = [
  {
    Header: "Склад",
    accessor: "sklad",
  },
  {
    Header: "Дата и время проведения",
    accessor: "time",
  },
  {
    Header: "Тип",
    accessor: "type",
  },
  {
    Header: "Результат",
    accessor: "result",
  },
  {
    Header: "Статус",
    accessor: "status",
  },
];

const InventoryTerminalPage: NextPage = () => {
  const [tableData, setTableData] = useState([]);

  useEffect(() => {
    getAllInventory({
      [QueryOptions.STATUS]: "opened",
    }).then((res) => {
      const data = res.data.map((item: any) => ({
        sklad: item.sklad,
        time:
          new Date(item.time).getFullYear() < 2022
            ? "-"
            : dateToString(item.time, false, false),
        type: item.type === "partial" ? "Частичная" : "Полная",
        result: formatNumber(item.result, true, true),
        status: item.status === "opened" ? "На редактировании" : "Проведенная",
        id: item.id,
      }));
      setTableData(data);
    });
  }, []);

  const navigation = useMemo(
    () => [
      {
        id: 0,
        name: "История",
        component: tableData ? (
          <Table columns={columns} data={tableData} editable onlyEditable />
        ) : (
          <></>
        ),
      },
      // {
      //   id: 1,
      //   name: "Новая инвентаризация",
      //   component: (
      //     <div className="w-full flex flex-col p-8">
      //       <InventoryContent fromTerminal />
      //     </div>
      //   ),
      // },
    ],
    [tableData],
  );

  return (
    <TerminalLayout>
      <TerminalCabinetLayout navigation={navigation} />
    </TerminalLayout>
  );
};

export default InventoryTerminalPage;
