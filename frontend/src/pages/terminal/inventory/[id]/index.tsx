import React, { FC, useEffect, useMemo, useState } from "react";
import { NextPage } from "next";
import { Column, Row } from "react-table";

import { useRouter } from "next/router";
import { Input } from "@shared/ui/Input";
import {
  getInventoryById,
  InventoryData,
  InventoryItem,
  updateInventory,
} from "@api/inventory";
import Search from "@common/Search";
import Table from "@common/Table";
import TerminalLayout from "@layouts/TerminalLayout";
import { sortTable } from "../../../inventory/[id]";
import NumberKeyboard from "@modules/NumberKeyboard";
import clsx from "clsx";

export interface InventoryListData {
  cost: string;
  difference: string;
  difference_sum: string;
  expenses: string;
  fact_quantity: number;
  fact_quantity_sum: string;
  id: number;
  income: string;
  inventarization_id: number;
  item_id: number;
  item_name: string;
  measure: string;
  plan_quantity: string;
  removed: string;
  removed_sum: string;
  sklad_id: number;
  sklad_name: string;
  start_quantity: string;
  status: string;
  time: string;
  type: string;
}

export const inventoryCacheKey = "zebra.cache.inventory.terminal";

const EditableCell: FC<{
  value: any;
  row: Row;
  column: Column;
  updateMyData: (row: Row, value: string) => void;
  activeInput: number;
  onClick: (row: Row) => void;
}> = ({ value: initialValue, row, updateMyData, activeInput, onClick }) => {
  // We need to keep and update the state of the cell normally
  const [value, setValue] = useState(initialValue);

  const onChange = (e: any) => {
    setValue(e.target.value);
  };

  // We'll only update the external data when the input is blurred
  const onBlur = () => {
    if (value !== initialValue) updateMyData(row, value);
  };

  // If the initialValue is changed external, sync it up with our state
  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  return (
    <div className="w-full flex items-center space-x-3">
      <Input
        className={clsx([
          "w-3/5 cursor-pointer",
          // @ts-ignore
          activeInput === row.original.id && "ring ring-primary",
        ])}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        onClick={(e) => {
          onClick(row);
          (e.target as HTMLInputElement).select();
        }}
        readOnly
      />
      {/* @ts-ignore */}
      <span>{row.original.measure}</span>
    </div>
  );
};

const EditInventoryTerminalPage: NextPage = () => {
  const router = useRouter();

  const inventoryId = useMemo(() => router.query.id, [router]);

  useEffect(() => {
    if (!inventoryId) return;
    localStorage.getItem(`${inventoryCacheKey}.${inventoryId}`);
  }, [inventoryId]);

  const [tableData, setTableData] = useState<Array<InventoryListData>>([]);
  const [filteredData, setFilteredData] = useState<Array<InventoryListData>>(
    []
  );
  const [filterText, setFilterText] = useState("");
  const [data, setData] = useState<InventoryData | null>(null);
  const [items, setItems] = useState<InventoryItem[]>([]);

  const [sortInfo, setSortInfo] = useState<any>({
    item_name: true,
    fact_quantity: true,
  });

  const [submitLoading, setSubmitLoading] = useState(false);

  const [isModified, setIsModified] = useState<boolean>(false);

  const [activeInput, setActiveInput] = useState<number>(1);
  const [resetInputValue, setResetInputValue] = useState<boolean>(true);

  const columns: Column[] = useMemo(
    () => [
      {
        Header: "Наименование",
        accessor: "item_name",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "item_name", "string", false, isAsc);
        },
      },
      {
        Header: "Факт. остаток",
        accessor: "fact_quantity",
        Cell: (cellProps) => (
          <EditableCell
            value={cellProps.value}
            row={cellProps.row}
            column={cellProps.column}
            updateMyData={(row: Row, value: string) => {
              setIsModified(true);
              setTableData((prevState) =>
                prevState.map((inventory: any) => {
                  // @ts-ignore
                  if (row.original.id !== inventory.id) return inventory;
                  const fact_quantity =
                    parseFloat(value.replace(/,/g, ".")) || 0;
                  return {
                    ...inventory,
                    fact_quantity,
                  };
                })
              );
              setItems((prevState) =>
                prevState.map((item) => {
                  // @ts-ignore
                  if (item.id !== row.original.id) return item;
                  return {
                    ...item,
                    fact_quantity: parseFloat(value.replace(/,/g, ".")) || 0,
                  };
                })
              );
            }}
            activeInput={activeInput}
            onClick={(row) => {
              setResetInputValue(true);
              // @ts-ignore
              setActiveInput(row.original.id);
            }}
          />
        ),
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "fact_quantity", "number", false, isAsc);
        },
      },
    ],
    [sortInfo, activeInput]
  );

  useEffect(() => {
    if (!inventoryId) return;
    const savedDataString = localStorage.getItem(
      `${inventoryCacheKey}.${inventoryId}`
    );
    getInventoryById(inventoryId as string).then((res) => {
      setData(res.data);
      setItems(res.data.items);
      const data = res.data.items.map((item: any) => ({
        ...item,
        fact_quantity: item.fact_quantity,
      }));
      if (!savedDataString) {
        setTableData(data);
        setFilteredData(data);
      } else {
        const savedData = JSON.parse(savedDataString);
        setTableData(savedData);
        setFilteredData(savedData);
      }
    });
  }, [inventoryId]);

  useEffect(() => {
    if (!tableData) return;
    const searchValue = filterText.toLowerCase().trim();
    const searchedData = searchValue
      ? tableData.filter((item) =>
          item.item_name.toLowerCase().includes(searchValue)
        )
      : tableData;
    setFilteredData(searchedData);
  }, [filterText, tableData]);

  useEffect(() => {
    if (!tableData || !inventoryId) return;
    isModified &&
      localStorage.setItem(
        `${inventoryCacheKey}.${inventoryId}`,
        JSON.stringify(tableData)
      );
  }, [tableData, inventoryId, isModified]);

  return (
    <TerminalLayout>
      <div className="flex flex-col flex-1">
        <div className="border-b shadow-sm z-10">
          <Search
            onChange={(e) =>
              setFilterText((e.target as HTMLInputElement).value)
            }
          />
        </div>
        <div className="h-full overflow-hidden flex-1 grid grid-cols-2 grid-rows-1">
          <div className="flex flex-col border-r overflow-y-auto">
            {tableData && (
              <Table
                columns={columns}
                data={filteredData.filter((item: any) => item.is_visible)}
                editable={false}
                onHeaderClick={(header) => {
                  // @ts-ignore
                  const sortedData = header.sort(
                    sortInfo[header.id],
                    tableData
                  );
                  setSortInfo((prevState: any) => ({
                    ...prevState,
                    [header.id]: !prevState[header.id],
                  }));
                  setTableData(sortedData);
                }}
              />
            )}
          </div>
          <div className="bg-gray-100 flex items-center justify-center">
            <NumberKeyboard
              onKeyClick={(value) => {
                setIsModified(true);
                setTableData((prevState) =>
                  prevState.map((inventory: any) => {
                    if (activeInput !== inventory.id) return inventory;
                    const fact_quantity = resetInputValue
                      ? value
                      : `${inventory.fact_quantity}${value}`;
                    return {
                      ...inventory,
                      fact_quantity,
                    };
                  })
                );
                setItems((prevState) =>
                  prevState.map((item) => {
                    if (item.id !== activeInput) return item;
                    const fact_quantity = resetInputValue
                      ? value
                      : `${item.fact_quantity}${value}`;
                    return {
                      ...item,
                      fact_quantity,
                    };
                  })
                );
                setResetInputValue(false);
              }}
              onDelete={() => {
                setTableData((prevState) =>
                  prevState.map((inventory: any) => {
                    if (activeInput !== inventory.id) return inventory;
                    const fact_quantity = `${inventory.fact_quantity}`.slice(
                      0,
                      -1
                    );
                    return {
                      ...inventory,
                      fact_quantity,
                    };
                  })
                );
                setItems((prevState) =>
                  prevState.map((item) => {
                    if (item.id !== activeInput) return item;
                    const fact_quantity = `${item.fact_quantity}`.slice(0, -1);
                    return {
                      ...item,
                      fact_quantity,
                    };
                  })
                );
                setResetInputValue(false);
              }}
            />
          </div>
        </div>
        <div className="h-16 px-5 py-3 flex items-center justify-between border-t border-gray-300">
          {data && (
            <div className="flex items-center space-x-3">
              <button
                onClick={() => {
                  if (!data) return;
                  if (window !== undefined) {
                    localStorage.removeItem(
                      `${inventoryCacheKey}.${inventoryId}`
                    );
                  }
                  setSubmitLoading(true);
                  updateInventory({
                    ...data,
                    status: "closed",
                    items: items.map((item) => ({
                      ...item,
                      fact_quantity: parseFloat(item.fact_quantity as string),
                    })),
                  })
                    .then(() => router.push("/terminal/order"))
                    .catch(() => setSubmitLoading(false));
                }}
                disabled={submitLoading}
                className="px-3 pb-2 pt-1.5 disabled:bg-gray-400/60 bg-primary hover:bg-primary/80 rounded text-white font-medium"
              >
                Сохранить
              </button>
              <button
                onClick={() => {
                  if (window !== undefined) {
                    localStorage.removeItem(
                      `${inventoryCacheKey}.${inventoryId}`
                    );
                  }
                  router.reload();
                }}
                className="px-3 pb-2 pt-1.5 hover:bg-gray-100 border border-gray-300 rounded text-gray-600 font-medium"
              >
                Сбросить
              </button>
            </div>
          )}
        </div>
      </div>
    </TerminalLayout>
  );
};

export default EditInventoryTerminalPage;
