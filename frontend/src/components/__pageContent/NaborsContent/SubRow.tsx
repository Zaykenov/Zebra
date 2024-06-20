import React, { FC, useEffect, useMemo, useState } from "react";
import { Row } from "react-table";
import { getMasterModifierNabor, getModifier } from "@api/modifiers";
import useMasterRole from "@hooks/useMasterRole";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const isMaster = useMasterRole();

  const [data, setData] = useState<any>(null);

  useEffect(() => {
    if (isMaster === null) return;
    const id = rowData.original.id;
    (async function () {
      const res = isMaster
        ? await getMasterModifierNabor(id)
        : await getModifier(id);
      setData(res.data);
    })();
  }, [rowData, isMaster]);

  const ingredients = useMemo(
    () => (isMaster ? data?.ingredient_nabor : data?.nabor_ingredient),
    [isMaster, data]
  );

  return (
    <div className="p-4 px-10 bg-white">
      <div className="w-full flex flex-col pr-4">
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          <div className="grow">Модификатор</div>
          <div className="w-1/6 text-right">Брутто</div>
          <div className="w-1/6 text-right">Цена</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {ingredients &&
            ingredients.map((ingredient: any) => (
              <div key={`tech-${ingredient.id}`} className="flex text-sm">
                <div className="grow flex flex-col space-y-3">
                  <span>{ingredient.name}</span>
                </div>
                <div className="w-1/6 text-right">{ingredient.brutto}</div>
                <div className="w-1/6 text-right">
                  {ingredient.price.toFixed(2)} ₸
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SubRow;
