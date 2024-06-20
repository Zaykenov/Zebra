import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import { getSupply } from "@api/supplies";

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

export default SubRow;
