import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import { getTransferById } from "@api/transfers";
import { formatNumber } from "@utils/formatNumber";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    getTransferById(rowData.original.id).then((res) => {
      setData(res.data);
    });
  }, []);

  return (
    <div className="p-4 px-10 bg-white">
      <div className="w-full flex flex-col pr-4">
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          <div className="grow">Название</div>
          <div className="w-1/5 text-right">Кол-во</div>
          <div className="w-1/5 text-right">Сумма</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {data &&
            data.item_transfers.map((item: any) => (
              <div key={`tech-${item.id}`} className="flex text-sm">
                <div className="grow flex flex-col space-y-3">
                  <span>{item.item_name}</span>
                </div>
                <div className="w-1/5 text-right">{item.quantity}</div>
                <div className="w-1/5 text-right">
                  {formatNumber(item.sum, true, true)}
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SubRow;
