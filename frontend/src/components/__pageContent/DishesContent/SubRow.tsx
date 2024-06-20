import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import { getDish } from "@api/dishes";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    getDish(rowData.original.id).then((res) => {
      setData(res.data);
    });
  }, []);

  return (
    <div className="p-4 px-10 bg-white">
      <div className="w-full flex flex-col pr-4">
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          <div className="grow">Ингредиент</div>
          <div className="w-1/6 text-right">Брутто</div>
          <div className="w-1/6 text-right">Себестоимость без НДС</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {data &&
            data.ingredient_tech_cart.map((ingredient: any) => (
              <div key={`tech-${ingredient.id}`} className="flex text-sm">
                <div className="grow flex flex-col space-y-3">
                  <span>{ingredient.name}</span>
                </div>
                <div className="w-1/6 text-right">{ingredient.brutto}</div>
                <div className="w-1/6 text-right">
                  {ingredient.cost.toFixed(2)} ₸
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SubRow;
