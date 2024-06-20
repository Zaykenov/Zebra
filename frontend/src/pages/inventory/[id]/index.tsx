import React, { FC, useCallback, useEffect, useMemo, useState } from "react";
import { NextPage } from "next";
import MainLayout from "@layouts/MainLayout";
import Table from "@common/Table";
import PageLayout from "@layouts/PageLayout";
import { Column, Row } from "react-table";
import {
  getInventoryById,
  getInventoryExpenseDetails,
  getInventoryIncomeDetails,
  getInventoryWasteDetails,
  InventoryData,
  InventoryItem,
  updateInventory,
} from "@api/inventory";
import { useRouter } from "next/router";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import clsx from "clsx";
import { Input } from "@shared/ui/Input";
import Search from "@common/Search";
import { Spin } from "antd";
import { Popover } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import { getExcelFile, getParitalInventExcelFile } from "@api/excel";
interface InventoryListData {
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

const EditableCell: FC<{
  value: any;
  row: Row;
  column: Column;
  updateMyData: (row: Row, value: string) => void;
  disabled?: boolean;
}> = ({
  value: initialValue,
  row,
  updateMyData, // This is a custom function that we supplied to our table instance
  disabled,
}) => {
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
    <Input
      className=""
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      disabled={disabled}
    />
  );
};

export const sortTable = (
  data: InventoryListData[],
  sortField: string,
  sortType: "string" | "number" | "date",
  isParseFloat: boolean = false,
  isAsc = true
) => {
  const sortedData = [...data];
  sortedData.sort((itemA, itemB) => {
    // @ts-ignore
    const fieldA = isAsc ? itemA[sortField] : itemB[sortField];
    // @ts-ignore
    const fieldB = isAsc ? itemB[sortField] : itemA[sortField];

    switch (sortType) {
      case "string":
        return fieldA.localeCompare(fieldB);
      case "date":
        const dateA = new Date(fieldA);
        const dateB = new Date(fieldB);
        return dateA > dateB ? 1 : -1;
      case "number":
        const valueA = isParseFloat
          ? parseFloat(fieldA.replace(",", ".").replace(" ", ""))
          : fieldA;
        const valueB = isParseFloat
          ? parseFloat(fieldB.replace(",", ".").replace(" ", ""))
          : fieldB;
        return valueA - valueB;
      default:
        return 1;
    }
  });
  return sortedData;
};

const EditInventoryPage: NextPage = () => {
  const router = useRouter();

  const inventoryId = useMemo(() => router.query.id, [router]);

  const [tableData, setTableData] = useState<Array<InventoryListData>>([]);
  const [filteredData, setFilteredData] = useState<Array<InventoryListData>>(
    []
  );
  const [filterText, setFilterText] = useState("");
  const [data, setData] = useState<InventoryData | null>(null);
  const [items, setItems] = useState<InventoryItem[]>([]);

  const [isModified, setIsModified] = useState<boolean>(false);

  useEffect(() => {
    if (!inventoryId) return;
    localStorage.getItem(`zebra.cache.inventory.${inventoryId}`);
  }, [inventoryId]);

  const [sortInfo, setSortInfo] = useState<any>({
    item_name: true,
    time: true,
    start_quantity: true,
    income: true,
    expenses: true,
    removed: true,
    removed_sum: true,
    plan_quantity: true,
    fact_quantity: true,
    fact_quantity_sum: true,
    difference: true,
    difference_sum: true,
  });

  const [isLoading, setIsLoading] = useState(true);

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
        Header: "Последняя проверка",
        accessor: "before_time",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "before_time", "date", false, isAsc);
        },
        Cell: ({ value }) => (
          <span>
            {/* @ts-ignore */}
            {new Date(value).getFullYear() < 2022
              ? "-"
              : dateToString(value, false, false)}
          </span>
        ),
      },
      {
        Header: "Нач. остаток",
        accessor: "start_quantity",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "start_quantity", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <span>
            {/* @ts-ignore */}
            {formatNumber(value, false, false)} {row.original.measure}
          </span>
        ),
      },
      {
        Header: "Поступление",
        accessor: "income",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "income", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <DetailsPopover value={value} row={row} option="income" />
        ),
      },
      {
        Header: "Расход",
        accessor: "expenses",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "expenses", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <DetailsPopover value={value} row={row} option="expense" />
        ),
      },
      {
        Header: "Списано",
        accessor: "removed",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "removed", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <DetailsPopover value={value} row={row} option="waste" />
        ),
      },
      {
        Header: "Списано, тг.",
        accessor: "removed_sum",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "removed_sum", "number", false, isAsc);
        },
        Cell: ({ value }) => (
          <span>
            {/* @ts-ignore */}
            {formatNumber(value, true, true)}
          </span>
        ),
      },
      {
        Header: "План. остаток",
        accessor: "plan_quantity",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "plan_quantity", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <span>
            {/* @ts-ignore */}
            {formatNumber(value, false, false)} {row.original.measure}
          </span>
        ),
      },
      {
        Header: "Факт. остаток",
        accessor: "fact_quantity",
        Cell: ({ value, row, column }) => {
          if (!data) return <></>;
          // @ts-ignore
          return data.status === "opened" ? (
            <EditableCell
              value={value}
              row={row}
              column={column}
              updateMyData={(row: Row, value: string) => {
                setIsModified(true);
                // @ts-ignore
                setTableData((prevState) =>
                  prevState.map((inventory: any) => {
                    // @ts-ignore
                    if (row.original.id !== inventory.id) return inventory;
                    const fact_quantity =
                      parseFloat(value.replace(/,/g, ".").replace(/\s/g, "")) ||
                      0;
                    const fact_quantity_sum = fact_quantity * inventory.cost;
                    const difference = fact_quantity - inventory.plan_quantity;
                    const difference_sum = difference * inventory.cost;
                    return {
                      ...inventory,
                      fact_quantity,
                      fact_quantity_sum,
                      difference,
                      difference_sum,
                    };
                  })
                );
                setItems((prevState) =>
                  prevState.map((item) => {
                    // @ts-ignore
                    if (item.id !== row.original.id) return item;
                    return {
                      ...item,
                      fact_quantity: parseFloat(value.replace(/,/g, ".")),
                    };
                  })
                );
              }}
            />
          ) : (
            <span>
              {/* @ts-ignore */}
              {value} {row.original.measure}
            </span>
          );
        },
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "fact_quantity", "number", false, isAsc);
        },
      },
      {
        Header: "Сумма факт. остатка",
        accessor: "fact_quantity_sum",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "fact_quantity_sum", "number", false, isAsc);
        },
        Cell: ({ value }) => (
          <span>
            {formatNumber(typeof value === "undefined" ? 0 : value, true, true)}
          </span>
        ),
      },
      {
        Header: "Разница",
        accessor: "difference",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "difference", "number", false, isAsc);
        },
        Cell: ({ value, row }) => (
          <span>
            {/* @ts-ignore */}
            {formatNumber(value, false, false)} {row.original.measure}
          </span>
        ),
      },
      {
        Header: "Разница, тг.",
        accessor: "difference_sum",
        sort: (isAsc: boolean, data: InventoryListData[]) => {
          return sortTable(data, "difference_sum", "number", false, isAsc);
        },
        Cell: ({ value }) => (
          <div
            className={clsx([
              "w-full text-right",
              value > 0 && "text-teal-600",
              value < 0 && "text-red-500",
            ])}
          >
            {/* @ts-ignore */}
            {value > 0 && "+"}
            {formatNumber(value, true, true)}
          </div>
        ),
      },
    ],
    [data, sortInfo]
  );

  useEffect(() => {
    if (!data) return;
    if (!inventoryId) return;
    if (isLoading) {
      setInterval(() => {
        getInventoryById(inventoryId as string).then((res) => {
          res.data.loading_status === "completed" && router.reload();
        });
      }, 1000);
    }
  }, [data, inventoryId]);

  useEffect(() => {
    if (!inventoryId) return;
    const savedDataString = localStorage.getItem(
      `zebra.cache.inventory.${inventoryId}`
    );
    getInventoryById(inventoryId as string).then((res) => {
      setData(res.data);
      if (res.data.loading_status !== "completed") {
        setIsLoading(true);
      } else {
        setIsLoading(false);
        setItems(res.data.items);
        const data = res.data.items.map((item: any) => ({
          ...item,
          fact_quantity_sum: item.fact_quantity * item.cost,
        }));
        if (!savedDataString) {
          setTableData(data);
          setFilteredData(data);
        } else {
          const savedData = JSON.parse(savedDataString);
          const newInventoryItems =
            savedData.length === res.data.items.length
              ? null
              : res.data.items.filter(
                  (loadedInventoryItem: any) =>
                    !savedData.some(
                      (savedInventoryItem: any) =>
                        loadedInventoryItem.id === savedInventoryItem.id
                    )
                );
          const mixData = newInventoryItems
            ? [...newInventoryItems, ...savedData]
            : savedData;
          localStorage.setItem(
            `zebra.cache.inventory.${inventoryId}`,
            JSON.stringify(mixData)
          );
          setTableData(mixData);
          setFilteredData(mixData);
        }
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
    if (!data || !tableData || !inventoryId) return;
    data.status === "opened" &&
      isModified &&
      localStorage.setItem(
        `zebra.cache.inventory.${inventoryId}`,
        JSON.stringify(tableData)
      );
  }, [data, tableData, inventoryId, isModified]);

  const downloadWorkbook = async () => {
    const generateTable = (products: InventoryListData[]) => {
      return [
        ["Название товара", "Разница(кол-во)", "Разница(сумма)"],
        ...products.map((product) => [
          product.item_name,
          product.difference,
          parseInt(product.difference_sum),
        ]),
        [
          "",
          "Итого",
          products.reduce(
            (sum, product) => sum + parseInt(product.difference_sum),
            0
          ),
        ],
      ];
    };
    const tables: [InventoryListData[], string, number][] = [
      [
        tableData.filter(
          (product) =>
            Number(product.difference_sum) === 0 &&
            Number(product.difference) === 0
        ),
        "Ровно",
        1,
      ],
      [
        tableData.filter(
          (product) =>
            Number(product.difference_sum) < 0 || Number(product.difference) < 0
        ),
        "Недостачи",
        6,
      ],
      [
        tableData.filter(
          (product) =>
            Number(product.difference_sum) > 0 || Number(product.difference) > 0
        ),
        "Излишки",
        11,
      ],
    ];

    const excelTableData: any[] = [];
    tables.forEach(([products]) => {
      //@ts-ignore
      const filteredProducts = products.filter((product) => product.is_visible);
      const mappedProducts = filteredProducts.map((product) => {
        //@ts-ignore
        product.difference = Number(product.difference);
        //@ts-ignore
        product.difference_sum = Number(product.difference_sum);
        return product;
      });
      const sortedProducts = mappedProducts.sort(
        (a, b) => parseInt(a.difference_sum) - parseInt(b.difference_sum)
      );
      excelTableData.push(generateTable(sortedProducts));
    });
    await getParitalInventExcelFile(
      `${tableData[0].sklad_name} инвентаризация`,
      excelTableData
    );
  };

  const EditButton = (
    <button
      className="text-blue-500 hover:underline text-sm"
      onClick={() => {
        data && router.push(`/inventory/inventory_form/${data.id}`);
      }}
    >
      Изменить параметры
    </button>
  );

  if (isLoading) {
    return (
      <PageLayout>
        <MainLayout title="Инвентаризация">
          <div className="w-full h-full flex items-center justify-center">
            <Spin size="large" />
          </div>
        </MainLayout>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <MainLayout
        title="Инвентаризация"
        content={
          data ? (
            <HeaderContent
              date={dateToString(
                data.time || new Date().toISOString(),
                false,
                false
              )}
              sklad={data.sklad || ""}
            />
          ) : (
            <></>
          )
        }
        backBtn
        customBtns={[EditButton]}
      >
        <div className="flex flex-col flex-1 h-full overflow-y-hidden">
          <div className="flex justify-between">
            <Search
              onChange={(e) =>
                setFilterText((e.target as HTMLInputElement).value)
              }
            />
            <button
              className="text-white shadow-md h-10 mt-2 pt-1.5 pb-2 mr-5 px-3 bg-primary text-sm font-semibold rounded-md hover:bg-teal-600"
              onClick={downloadWorkbook}
            >
              Экспорт в Excel
            </button>
          </div>

          <div className="flex flex-col flex-1 h-full overflow-y-auto">
            {filteredData && (
              <Table
                columns={columns}
                data={filteredData.filter((item: any) => item.is_visible)}
                editable={data?.status === "opened"}
                onlyDeletable
                customOnDelete={(row) => {
                  const inventoryId = router.query.id;
                  const cachedDataStr = localStorage.getItem(
                    //  @ts-ignore
                    `zebra.cache.inventory.${inventoryId}`
                  );
                  if (!cachedDataStr) return;
                  const cachedData = JSON.parse(cachedDataStr);

                  localStorage.setItem(
                    `zebra.cache.inventory.${inventoryId}`,
                    JSON.stringify(
                      cachedData.filter(
                        (inventoryItem: any) =>
                          //  @ts-ignore
                          inventoryItem.id !== row.original.id
                      )
                    )
                  );
                }}
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
          <div className="px-5 py-3 flex items-center justify-between border-t border-gray-300">
            {data && data.status === "opened" ? (
              <div className="flex items-center space-x-3">
                <button
                  onClick={() => {
                    if (window !== undefined) {
                      localStorage.removeItem(
                        `zebra.cache.inventory.${inventoryId}`
                      );
                    }
                    data &&
                      updateInventory({
                        ...data,
                        status: "closed",
                        items,
                      }).then(() => router.reload());
                  }}
                  className="px-3 pb-2 pt-1.5 bg-primary hover:bg-primary/80 rounded text-white font-medium"
                >
                  Сохранить
                </button>
                <button
                  onClick={() => {
                    if (window !== undefined) {
                      localStorage.removeItem(
                        `zebra.cache.inventory.${inventoryId}`
                      );
                    }
                    router.reload();
                  }}
                  className="px-3 pb-2 pt-1.5 hover:bg-gray-100 border border-gray-300 rounded text-gray-600 font-medium"
                >
                  Сбросить
                </button>
              </div>
            ) : (
              <button
                onClick={() => {
                  data &&
                    updateInventory({
                      ...data,
                      status: "opened",
                    }).then(() => router.reload());
                }}
                className="px-3 pb-2 pt-1.5 bg-primary hover:bg-primary/80 rounded text-white font-medium"
              >
                Редактировать
              </button>
            )}
            <div className="text-lg flex items-center space-x-3">
              <span>Итого:</span>
              <span
                className={clsx([
                  data && data.result !== undefined && data.result >= 0
                    ? "text-primary"
                    : "text-red-500",
                  "font-semibold",
                ])}
              >
                {(data && data.result !== undefined && data.result > 0
                  ? "+"
                  : "") +
                  formatNumber(data ? (data.result as number) : 0, true, true)}
              </span>
            </div>
          </div>
        </div>
      </MainLayout>
    </PageLayout>
  );
};

const HeaderContent: FC<{ date: string; sklad: string }> = ({
  date,
  sklad,
}) => {
  return (
    <div className="flex items-center space-x-2">
      <div className="font-light">
        Инвентаризация за <span className="font-medium">{date}</span> в{" "}
        <span className="font-medium">{sklad}</span>
      </div>
    </div>
  );
};

interface DetailsPopoverProps {
  value: number;
  row: Row;
  option: "income" | "expense" | "waste";
}

const DetailsPopover: FC<DetailsPopoverProps> = ({ value, row, option }) => {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<
    {
      cost?: number;
      measurement: string;
      quantity: number;
      sum?: number;
      time: string;
      name?: string;
    }[]
  >([]);

  const handleResponse = useCallback((res: any) => {
    setLoading(false);
    setData(res.data);
  }, []);

  const handleClick = useCallback(() => {
    // @ts-ignore
    const inventoryItemId = row.original.id;
    switch (option) {
      case "income":
        getInventoryIncomeDetails(inventoryItemId).then(handleResponse);
        break;
      case "expense":
        getInventoryExpenseDetails(inventoryItemId).then(handleResponse);
        break;
      case "waste":
        getInventoryWasteDetails(inventoryItemId).then(handleResponse);
        break;
      default:
        break;
    }
  }, [handleResponse, row]);

  const text = useMemo(() => {
    switch (option) {
      case "income":
        return "Поставки";
      case "expense":
        return "Расходы";
      case "waste":
        return "Списания";
      default:
        return "";
    }
  }, [option]);

  return (
    <Popover className="relative">
      <Popover.Button
        onClick={handleClick}
        className="flex items-center space-x-3"
      >
        <span>
          {/* @ts-ignore */}
          {formatNumber(value, false, false)} {row.original.measure}
        </span>
        <div className="p-0.5 border border-indigo-300 rounded">
          <ChevronDownIcon className="w-2 h-2 text-indigo-300" />
        </div>
      </Popover.Button>

      <Popover.Panel className="absolute z-10 mt-5 min-w-[400px] -ml-40 p-4 border border-gray-400 shadow-lg rounded bg-white">
        {loading ? (
          <div className="w-full h-full flex items-center justify-center">
            <Spin size="default" />
          </div>
        ) : data.length === 0 ? (
          <div className="w-full h-full flex items-center justify-center">
            {text} отсутствуют
          </div>
        ) : (
          <div className="w-full h-full flex flex-col">
            <div className="py-1 border-b border-gray-300 font-bold">
              {text}:
            </div>
            <ul className="flex flex-col space-y-2 py-2">
              {data.map((detail) => (
                <li className="w-full flex items-center justify-between px-3 text-sm">
                  <span>{dateToString(detail.time, false, false)}</span>
                  {detail.cost ? (
                    <span>
                      {formatNumber(detail.quantity, false, true)}{" "}
                      {detail.measurement}
                    </span>
                  ) : (
                    <span>{detail.name}</span>
                  )}
                  {detail.sum ? (
                    <span>{formatNumber(detail.sum, true, true)}</span>
                  ) : (
                    <span>
                      {formatNumber(detail.quantity, false, true)}{" "}
                      {detail.measurement}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          </div>
        )}
      </Popover.Panel>
    </Popover>
  );
};

export default EditInventoryPage;
