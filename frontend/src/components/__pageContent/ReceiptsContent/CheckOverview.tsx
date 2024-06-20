import React, { FC, useEffect, useState } from "react";
import { Row } from "@tanstack/react-table";
import { getCheckById } from "@api/check";
import { formatNumber } from "@utils/formatNumber";

const CheckOverview: FC<{ rowData: Row<any>; withIngredients?: boolean }> = ({
  rowData,
  withIngredients = false,
}) => {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    getCheckById({ id: rowData.original.id }).then((res) => {
      setData(res.data);
    });
  }, [rowData]);

  return data ? (
    <div className="w-full flex flex-col pr-4">
      <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
        <div className="grow">Товар</div>
        {!withIngredients && <div className="w-1/6 text-right">Цена</div>}
        <div className="w-1/6 text-right">Кол-во</div>
        <div className="w-1/6 text-right">Итого</div>
      </div>
      <div className="flex flex-col space-y-3 mb-4 pt-2">
        {data.techCartCheck.map((techCart: any) => (
          <>
            <div key={`tech-${techCart.id}`} className="flex text-sm">
              <div className="grow flex flex-col space-y-3">
                <span>{techCart.name}</span>
              </div>
              {!withIngredients && (
                <div className="w-1/6 text-right">
                  {techCart.price.toFixed(2)} ₸
                </div>
              )}
              <div className="w-1/6 text-right">{techCart.quantity}</div>
              <div className="w-1/6 text-right">
                {withIngredients
                  ? techCart.cost.toFixed(2)
                  : (techCart.price * techCart.quantity).toFixed(2)}{" "}
                ₸
              </div>
            </div>
            {withIngredients && (
              <ul className="w-full flex flex-col space-y-2 text-xs text-gray-500">
                {JSON.parse(techCart.ingredients).map((ingredient: any) => (
                  <li className="w-full flex items-center">
                    <div className="grow flex flex-col space-y-3 pl-3">
                      <span>--- {ingredient.name}</span>
                    </div>
                    <div className="w-1/6 text-right">
                      {formatNumber(
                        ingredient.brutto,
                        false,
                        ingredient.brutto % 1 !== 0,
                      )}{" "}
                      {ingredient.measure}
                    </div>
                    <div className="w-1/6 text-right">
                      {formatNumber(ingredient.cost, true, true)}
                    </div>
                  </li>
                ))}
                {JSON.parse(techCart.modificators).map((modificator: any) => (
                  <li className="w-full flex items-center">
                    <div className="grow flex flex-col space-y-3 pl-3">
                      <span>--- {modificator.name}</span>
                    </div>
                    <div className="w-1/6 text-right">
                      {modificator.brutto} {modificator.measure}
                    </div>
                    <div className="w-1/6 text-right">
                      {modificator.cost.toFixed(2)} ₸
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </>
        ))}
        {data.tovarCheck.map((tovar: any) => (
          <div key={`tovar-${tovar.id}`} className="flex text-sm">
            <div className="grow flex flex-col space-y-3">
              <span>{tovar.tovar_name}</span>
            </div>
            {!withIngredients && (
              <div className="w-1/6 text-right">{tovar.price.toFixed(2)} ₸</div>
            )}
            <div className="w-1/6 text-right">{tovar.quantity.toFixed(2)}</div>
            <div className="w-1/6 text-right">
              {withIngredients ? tovar.cost.toFixed(2) : tovar.price.toFixed(2)}{" "}
              ₸
            </div>
          </div>
        ))}
      </div>
      <div className="flex flex-col py-2 space-y-1.5">
        <div className="flex items-center space-x-3 font-semibold">
          <div className="text-sm">
            {withIngredients ? "Итоговая себестоимость" : "Итого"}
          </div>
          <hr className="grow" />
          <div className="text-sm tracking-wide">
            {withIngredients ? data.cost.toFixed(2) : data.sum.toFixed(2)} ₸
          </div>
        </div>
        {!withIngredients && (
          <div className="flex items-center space-x-3 font-semibold">
            <div className="text-sm">{"К оплате (с учетом скидки)"}</div>
            <hr className="grow" />
            <div className="text-base tracking-wide">
              {formatNumber(data.sum - data.discount, true, true)}
              {"   "}
              {data.discount_percent > 0 && (
                <span className="font-normal text-green-700">
                  ({data.discount_percent * 100}%)
                </span>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  ) : null;
};

export default CheckOverview;
