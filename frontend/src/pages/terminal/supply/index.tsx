import React, { FC, useCallback, useEffect, useMemo, useState } from "react";
import { NextPage } from "next";
import TerminalLayout from "@layouts/TerminalLayout";
import TerminalSupplyForm from "@modules/TerminalSupplyForm";
import TerminalCabinetLayout from "@layouts/TerminalCabinetLayout";
import { getAllSupplies, getSupply } from "@api/supplies";
import { dateToString } from "@api/check";
import Table from "@common/Table";
import { Column, Row } from "react-table";
import OldPagination from "@modules/OldPagination";

const columns: Column[] = [
  {
    Header: "Дата",
    accessor: "date",
  },
  {
    Header: "Поставщик",
    accessor: "dealer",
  },
  {
    Header: "Склад",
    accessor: "sklad",
  },
  {
    Header: "Счет",
    accessor: "schet",
  },
  {
    Header: "Товары",
    accessor: "items",
    Cell: ({ value }) => <span className="whitespace-normal">{value}</span>,
  },
  {
    Header: "Сумма",
    accessor: "sum",
  },
];

const TerminalSupplyPage: NextPage = () => {
  const [curPage, setCurPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [tableData, setTableData] = useState([]);
  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  useEffect(() => {
    getAllSupplies({ page: curPage }, false).then((res) => {
      setTotalPages(res.data.totalPages);
      const data = res.data.data.postavka.map((supply: any) => ({
        dealer: supply.dealer,
        sklad: supply.sklad,
        schet: supply.schet,
        date: dateToString(supply.time, false),
        sum: supply.sum,
        items: supply.items.map((item: any) => item.name).join(", "),
        id: supply.id,
        isDeleted: supply.deleted,
      }));
      setTableData(data);
    });
  }, [curPage]);

  const navigation = useMemo(
    () => [
      {
        id: 0,
        name: "История",
        component: tableData ? (
          <>
            <Table
              columns={columns}
              data={tableData}
              editable
              details={true}
              onlyEditable
              isRowEditable={(row) => !row.original.isDeleted}
              renderRowSubComponent={renderSubComponent}
              customRowStyle={(row) =>
                row.original.isDeleted ? "bg-red-200/80" : "bg-white-100/80"
              }
            />
            <OldPagination
              curPage={curPage}
              totalPages={totalPages}
              setPage={setCurPage}
              resultsNum={tableData.length}
              perPage={20}
              detailed
            />
          </>
        ) : (
          <></>
        ),
      },
      {
        id: 1,
        name: "Новая поставка",
        component: (
          <div className="w-full min-h-full py-10 flex flex-col items-center overflow-auto">
            <TerminalSupplyForm />
          </div>
        ),
      },
    ],
    [tableData]
  );

  const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
    const [data, setData] = useState<any>(null);

    useEffect(() => {
      getSupply(rowData.original.id).then((res) => {
        setData(res.data);
      });
    }, []);

    return (
      <div className="p-4 px-10 bg-white">
        <div className="w-full flex flex-col pr-4">
          <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
            <div className="grow">Товар</div>
            <div className="w-1/6 text-right">Кол-во</div>
            <div className="w-1/6 text-right">Сумма</div>
          </div>
          <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
            {data &&
              data.items.map((item: any) => (
                <div key={`tech-${item.id}`} className="flex text-sm">
                  <div className="grow flex flex-col space-y-3">
                    <span>{item.name}</span>
                  </div>
                  <div className="w-1/6 text-right">
                    {item.quantity} {item.measurement}
                  </div>
                  <div className="w-1/6 text-right">
                    {(item.quantity * item.cost).toFixed(2)} ₸
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
    );
  };

  return (
    <TerminalLayout>
      <TerminalCabinetLayout navigation={navigation} />
    </TerminalLayout>
  );
};

export default TerminalSupplyPage;
