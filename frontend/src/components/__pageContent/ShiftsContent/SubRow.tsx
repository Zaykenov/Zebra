import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import { getShiftById, ShiftCategory, shiftMapping } from "@api/shifts";
import clsx from "clsx";
import { formatNumber } from "@utils/formatNumber";
import Link from "next/link";
import { PlusIcon } from "@heroicons/react/24/outline";
import Table from "@common/Table";
import { dateToString } from "@api/check";
import { firstRow, secondRow, transactionColumns } from "./constants";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [data, setData] = useState<any>();
  useEffect(() => {
    rowData &&
      getShiftById(rowData.original.id).then((res) => {
        setData(res.data);
      });
  }, [rowData]);

  return data ? (
    <div className="w-full py-3 px-24 flex flex-col divide-y divide-y-gray-200">
      <div className="flex items-center space-x-4">
        {firstRow.map((item) => (
          <div className="flex flex-col py-2 w-48">
            <span className="text-sm">{item.text}:</span>
            <div className={clsx(["text-lg font-medium", item.color])}>
              {formatNumber(item.field ? data[item.field] : 0, true, true)}
            </div>
          </div>
        ))}
      </div>
      <div className="flex items-center space-x-4">
        {secondRow.map((item) => (
          <div className="flex flex-col py-2 w-48">
            <span className="">{item.text}:</span>
            <div className={clsx(["text-lg font-medium", item.color])}>
              {formatNumber(item.field ? data[item.field] : 0, true, true)}
            </div>
          </div>
        ))}
      </div>
      <div className="w-full flex flex-col items-start py-3">
        <Link
          href={{
            pathname: "/transactions/transaction_form",
            query: {
              shift: rowData?.original?.id || 1,
              routerName: "shifts",
            },
          }}
        >
          <button className="flex items-center space-x-2 text-indigo-500 mb-3 text-sm py-1 border-b border-transparent hover:border-indigo-500">
            <PlusIcon className="w-3 h-3" />
            <span>Добавить транзакцию</span>
          </button>
        </Link>
        <div className="w-full flex shadow-md">
          <Table
            columns={transactionColumns}
            data={data.transactions.map((transaction: any) => ({
              ...transaction,
              time: dateToString(transaction.time, false),
              category: shiftMapping[transaction.category as ShiftCategory],
              sum: formatNumber(transaction.sum, true, true),
            }))}
            editable
            isRowDeletable={(row) =>
              !(
                row.original.category === "Закрытие смены" ||
                row.original.category === "Открытие смены"
              )
            }
            customEditPath="/transactions"
            customEditBtn={(row) => (
              <Link
                href={{
                  pathname: `/transactions/${row.original.id}`,
                  query: {
                    shift: rowData.original.id || 1,
                  },
                }}
              >
                <a className="text-indigo-500 hover:text-indigo-600 hover:underline">
                  Ред.
                </a>
              </Link>
            )}
          />
        </div>
      </div>
    </div>
  ) : (
    <></>
  );
};

export default SubRow;
