import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import { getFeedbackById } from "@api/feedbacks";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    getFeedbackById(rowData.original.id).then((res) => {
      setData(res.data);
    });
  }, []);

  return (
    <div className="p-4 px-10 bg-white">
      <div className="w-full flex flex-col pr-4">
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          <div className="w-1/2 text-left">Товары</div>
          <div className="w-1/2 text-left">Тех. карты</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {data && (
            <div key={`tech-${data.id}`} className="flex text-sm">
              <div className="w-1/2 text-left">
                {JSON.parse(data.check_json).TechCarts.map(
                  (item: any) => item.Name
                )}
              </div>
              <div className="w-1/2 text-left">
                {JSON.parse(data.check_json).Tovars.map(
                  (item: any) => item.Name
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default SubRow;
