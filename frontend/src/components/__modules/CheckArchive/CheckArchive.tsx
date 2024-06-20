import React, { FC, useEffect, useState } from "react";
import {
  CheckData,
  dateToString,
  getCheckForPrint,
  printCheck,
} from "@api/check";
import clsx from "clsx";
import { processPrintObject } from "@utils/processPrintObject";

export interface CheckArchiveProps {
  data: CheckData[];
}

const CheckArchive: FC<CheckArchiveProps> = ({ data }) => {
  const [selectedCheck, setSelectedCheck] = useState<CheckData | null>(null);
  const [checkToPrint, setCheckToPrint] = useState<string | null>(null);

  useEffect(() => {
    selectedCheck &&
      getCheckForPrint({ id: selectedCheck.id as number })
        .then((res) => {
          const printObj = processPrintObject(res.data);
          setCheckToPrint(JSON.stringify(printObj));
        })
        .catch(() => {
          setCheckToPrint(null);
        });
  }, [selectedCheck]);

  return (
    <div className="w-full flex">
      <div className="w-[320px] flex flex-col bg-stone-100 overflow-y-scroll">
        {data &&
          data.map((check) => (
            <button
              key={check.id}
              onClick={() => {
                setSelectedCheck(check);
              }}
              className={clsx([
                "flex flex-col py-6 px-4 hover:bg-stone-200 space-y-1",
                selectedCheck?.id === check.id && "bg-stone-200",
              ])}
            >
              <div className="w-full flex items-center justify-between">
                <div className="font-medium">№{check.id}</div>
                <div className="text-sm">{check.sum.toFixed(2)} ₸</div>
              </div>
              <div className="text-sm text-gray-400">
                {/*{check.techCartCheck.map((techCart) => techCart.name)}*/}
              </div>
            </button>
          ))}
      </div>
      <div className="grow flex">
        {selectedCheck && (
          <div className="py-6 px-8 w-full flex flex-col">
            <div className="text-2xl mb-4">Чек №{selectedCheck.id}</div>
            <div className="overflow-y-scroll grow flex flex-col">
              <div className="w-1/2 mb-5 flex flex-col items-stretch space-y-3">
                <div className="flex items-center justify-between">
                  <span className="grow text-sm text-gray-500">Сотрудник</span>
                  <span className="w-1/3 text-lg">{selectedCheck.worker}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="grow text-sm text-gray-500">Открыт</span>
                  <span className="w-1/3 text-lg">
                    {dateToString(selectedCheck.opened_at as string, false)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="grow text-sm text-gray-500">
                    Счет закрыт
                  </span>
                  <span className="w-1/3 text-lg">
                    {dateToString(selectedCheck.closed_at as string, false)}
                  </span>
                </div>
              </div>
              <div className="w-full flex flex-col pr-4">
                <div className="flex text-sm font-bold py-3">
                  <div className="grow">Наименование</div>
                  <div className="w-1/6 text-right">Кол-во</div>
                  <div className="w-1/6 text-right">Цена</div>
                  <div className="w-1/6 text-right">Итого</div>
                </div>
                <div className="flex flex-col space-y-3 mb-6">
                  {selectedCheck.techCartCheck.map((techCart) => {
                    let modificators = techCart.modificators;
                    let productName = techCart.name;
                    if (modificators) {
                      //@ts-ignore
                      const modificatorsArray = JSON.parse(modificators).map(
                        (modificator: any) => modificator.name
                      );
                      const joinedStringOfModificators =
                        modificatorsArray.join(" + ");
                      productName =
                        joinedStringOfModificators === ""
                          ? techCart.name
                          : `${techCart.name} + ${joinedStringOfModificators}`;
                    }
                    return (
                      <div key={`tech-${techCart.id}`} className="flex text-sm">
                        <div className="grow flex flex-col space-y-3">
                          <span>{productName}</span>
                        </div>
                        <div className="w-1/6 text-right">
                          {techCart.quantity}
                        </div>
                        <div className="w-1/6 text-right">
                          {techCart.price.toFixed(2)}
                        </div>
                        <div className="w-1/6 text-right">
                          {(techCart.quantity * techCart.price).toFixed(2)}
                        </div>
                      </div>
                    );
                  })}
                  {selectedCheck.tovarCheck.map((tovar) => (
                    <div key={`tovar-${tovar.id}`} className="flex text-sm">
                      <div className="grow flex flex-col space-y-3">
                        <span>{tovar.tovar_name}</span>
                      </div>
                      <div className="w-1/6 text-right">{tovar.quantity}</div>
                      <div className="w-1/6 text-right">
                        {tovar.price.toFixed(2)}
                      </div>
                      <div className="w-1/6 text-right">
                        {(tovar.quantity * tovar.price).toFixed(2)}
                      </div>
                    </div>
                  ))}
                </div>
                <div className="flex flex-col border-t-2 border-b-2 border-gray-200 py-4 space-y-1.5">
                  <div className="flex items-center space-x-3">
                    <div className="text-lg">Итого</div>
                    <hr className="grow" />
                    <div className="text-lg tracking-wide">
                      {selectedCheck.sum.toFixed(2)} ₸
                    </div>
                  </div>
                  <div className="flex items-center space-x-3 font-bold">
                    <div className="text-lg">К оплате (со скидкой)</div>
                    <hr className="grow" />
                    <div className="text-lg tracking-wide">
                      {(selectedCheck.sum - selectedCheck.discount).toFixed(2)}{" "}
                      ₸ ({selectedCheck.discount_percent * 100}%)
                    </div>
                  </div>
                </div>
                <div className="flex flex-col space-y-1.5 pt-4">
                  <div className="text-lg font-bold">Оплата</div>
                  <div className="flex items-center space-x-3">
                    <div className="text-lg">{selectedCheck.payment}</div>
                    <hr className="grow" />
                    <div className="text-lg tracking-wide">
                      {(selectedCheck.sum - selectedCheck.discount).toFixed(2)}{" "}
                      ₸
                    </div>
                  </div>
                </div>
              </div>
              <div className="flex flex-col items-start py-5">
                <button
                  onClick={() => {
                    checkToPrint &&
                      printCheck(checkToPrint).then((res) => {
                        console.log(res);
                      });
                  }}
                  className="px-3 py-1 rounded bg-primary hover:bg-primary/80 text-white border border-primary"
                >
                  Печать
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default CheckArchive;
