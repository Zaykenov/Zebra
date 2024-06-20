import React, { FC, useCallback, useEffect, useMemo, useState } from "react";
import { Input } from "@shared/ui/Input";
import { Controller, useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { getAllSuppliers } from "@api/suppliers";
import { createSupply, SupplyData, updateSupply } from "@api/supplies";
import { LabeledSelect } from "@shared/ui/Select";
import { ClockIcon, PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { getItems } from "@api/index";
import { getAllSklads } from "@api/sklad";
import { getAllAccounts } from "@api/accounts";
import DatePicker from "react-datepicker";
import { formatInputValue } from "@utils/formatInputValue";
import Loader from "@common/Loader/Loader";
import SelectWithSearch from "@shared/ui/SelectWithSearch";
import clsx from "clsx";
import { formatNumber } from "@utils/formatNumber";
import { ChevronDownIcon, ChevronUpIcon } from "@heroicons/react/24/solid";
export interface SupplierFormProps {
  data?: any;
  isEdit?: boolean;
}

interface Ingredient {
  id: string;
  sklad_id: number;
  idx: number;
  type: string;
  measure: string;
  last_postavka_cost: number;
  name: string;
  value: string;
  data?: any;
}

const SupplyForm: FC<SupplierFormProps> = ({ data, isEdit = false }) => {
  const router = useRouter();

  const [products, setProducts] = useState<
    {
      id: string;
      quantity: number | string;
      type: string;
      measure: string;
      price: number | string;
      total: number | string;
      last_cost?: number;
    }[]
  >([]);
  const [productOptions, setProductOptions] = useState<Ingredient[]>([]);
  const [dealerOptions, setDealerOptions] = useState<any[]>([]);
  const [schetOptions, setSchetOptions] = useState<any[]>([]);
  const [skladOptions, setSkladOptions] = useState<any[]>([]);
  const [totalPrice, setTotalPrice] = useState<number>(0);

  const [optionsLoaded, setOptionsLoaded] = useState<boolean>(false);
  const [dataLoaded, setDataLoaded] = useState<boolean>(false);

  const [isLoading, setIsLoading] = useState<boolean>(false);

  const defaultValues = useMemo(
    () => ({
      dealer_id: 1,
      sklad_id: 1,
      schet_id: 1,
      time: new Date().toISOString(),
      date: new Date(),
      items: [],
    }),
    []
  );

  const { handleSubmit, register, reset, control } = useForm<SupplyData>({
    defaultValues,
  });

  useEffect(() => {
    if (!data) setIsLoading(true);
    else setIsLoading(false);
  }, [data]);

  useEffect(() => {
    if (!data || !productOptions.length || !optionsLoaded || dataLoaded) return;
    setDataLoaded(true);
    const dataValues = {
      dealer_id: data.dealer_id,
      sklad_id: data.sklad_id,
      schet_id: data.schet_id,
      time: data.time,
      date: new Date(data.time),
      items: data.items,
    };
    reset(dataValues);
    setProducts(
      data.items.map((item: any) => {
        const selectedOption = productOptions.find(
          (option) => option.value === item.item_id
        );
        return {
          id: `${item.item_id}_${item.type}`,
          quantity: item.quantity,
          type: item.type,
          measure: item.measurement,
          price: item.cost,
          total: item.quantity * item.cost,
          last_cost: selectedOption?.data?.last_cost || 0,
        };
      })
    );
  }, [data, reset, productOptions, optionsLoaded]);

  const fetchData = useCallback(async () => {
    if (isEdit === undefined || (isEdit && !data)) return;
    setIsLoading(true);
    const [sklads, accounts, suppliers] = await Promise.all([
      getAllSklads(),
      getAllAccounts(),
      getAllSuppliers(),
    ]);
    setIsLoading(false);
    const selectedSklad = isEdit ? data.sklad_id : sklads.data[0].id;
    getItems().then((res) => {
      const productData = res.data
        .filter((item: Ingredient) => item.sklad_id === selectedSklad)
        .map((item: Ingredient) => ({
          name: item.name,
          value: `${item.id}_${item.type}`,
          data: {
            type: item.type,
            measure: item.measure,
            last_cost: item.last_postavka_cost,
          },
        }));
      setProductOptions(productData);
    });
    setSkladOptions(
      sklads.data.map((item: any) => ({
        name: item.name,
        value: parseInt(item.id),
      }))
    );
    setSchetOptions(
      accounts.data.map((item: any) => ({
        name: item.name,
        value: parseInt(item.id),
      }))
    );
    setDealerOptions(
      suppliers.data.map((item: any) => ({
        name: item.name,
        value: parseInt(item.id),
      }))
    );
    setOptionsLoaded(true);
  }, [isEdit, data]);

  useEffect(() => {
    fetchData().then();
  }, [fetchData]);

  useEffect(() => {
    if (isEdit) return;
    if (
      skladOptions.length > 0 &&
      dealerOptions.length > 0 &&
      schetOptions.length > 0
    ) {
      const updatedValues = {
        ...defaultValues,
        sklad_id: skladOptions[0].value,
        dealer_id: dealerOptions[0].value,
        schet_id: schetOptions[0].value,
      };
      reset(updatedValues);
    }
  }, [isEdit, skladOptions, dealerOptions, schetOptions, defaultValues]);

  const defaultProductOption = useMemo(
    () => ({
      id: productOptions[0]?.value,
      last_cost: productOptions[0]?.data.last_cost,
      type: productOptions[0]?.data.type,
      measure: productOptions[0]?.data.measure,
      quantity: 0,
      price: 0,
      total: 0,
    }),
    [productOptions]
  );

  const onSubmit = useCallback(
    (submitData: SupplyData) => {
      setIsLoading(true);
      if (!data) {
        createSupply({
          ...submitData,
          time: submitData.date?.toISOString(),
          items: products.map((product) => ({
            item_id: parseInt(product.id),
            type: product.type,
            quantity: parseFloat(product.quantity as string),
            cost: parseFloat(product.price as string),
          })),
        })
          .then(() => router.replace("/supply"))
          .finally(() => setIsLoading(false));
      } else {
        updateSupply({
          ...submitData,
          time: submitData.date?.toISOString(),
          items: products.map((product) => ({
            item_id: parseInt(product.id),
            type: product.type,
            quantity: parseFloat(product.quantity as string),
            cost: parseFloat(product.price as string),
          })),
          id: data.id,
        })
          .then(() => router.replace("/supply"))
          .finally(() => setIsLoading(false));
      }
    },
    [products, data, router]
  );

  const handleSelectIngredient = (idx: number, value: string, data?: any) => {
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: value,
          quantity: option.quantity,
          price: option.price,
          total: option.total,
          type: data.type,
          measure: data.measure,
          last_cost: data.last_cost,
        };
      })
    );
  };

  const handleChangeQuantity = (
    idx: number,
    e: React.FormEvent<HTMLInputElement>
  ) => {
    const { inputValue, numberValue } = formatInputValue(
      (e.target as HTMLSelectElement).value
    );
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: option.id,
          price: option.price,
          type: option.type,
          measure: option.measure,
          total: numberValue * parseFloat(option.price as string),
          quantity: inputValue,
          last_cost: option.last_cost,
        };
      })
    );
  };

  const handleChangePrice = (
    idx: number,
    e: React.FormEvent<HTMLInputElement> | null,
    value?: string
  ) => {
    const { inputValue, numberValue } = formatInputValue(
      value || (e ? (e.target as HTMLSelectElement).value : "")
    );
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: option.id,
          quantity: option.quantity,
          type: option.type,
          measure: option.measure,
          total: numberValue * parseFloat(option.quantity as string),
          price: inputValue,
          last_cost: option.last_cost,
        };
      })
    );
  };

  const handleChangeTotalPrice = (
    idx: number,
    e: React.FormEvent<HTMLInputElement>
  ) => {
    const { inputValue, numberValue } = formatInputValue(
      (e.target as HTMLSelectElement).value
    );
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: option.id,
          quantity: option.quantity,
          type: option.type,
          measure: option.measure,
          total: inputValue,
          price: !isFinite(numberValue / parseFloat(option.quantity as string))
            ? 0
            : numberValue / parseFloat(option.quantity as string),
          last_cost: option.last_cost,
        };
      })
    );
  };

  const getTotalPrice = () => {
    let total = 0;
    products.map((product) => {
      total += product.total
        ? Math.round(
            (parseFloat(product.total as string) + Number.EPSILON) * 100
          ) / 100
        : 0;
    });
    return total;
  };

  useEffect(() => {
    setTotalPrice(getTotalPrice());
  }, [products]);
  if (isLoading) return <Loader />;
  else
    return (
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="flex flex-col space-y-5"
      >
        <div className="w-1/2 flex flex-col space-y-5">
          <LabeledSelect
            {...register("dealer_id", { valueAsNumber: true })}
            label="Поставщик"
            options={dealerOptions}
          />
          <LabeledSelect
            {...register("sklad_id", {
              valueAsNumber: true,
              onChange: (e) => {
                getItems().then((res) => {
                  const productData = res.data
                    .filter((item: Ingredient) => {
                      return item.sklad_id === parseInt(e.target.value);
                    })
                    .map((item: Ingredient) => ({
                      name: item.name,
                      value: `${item.id}_${item.type}`,
                      data: {
                        type: item.type,
                        measure: item.measure,
                        last_cost: item.last_postavka_cost,
                      },
                    }));
                  setProductOptions(productData);
                  setProducts((prevProductsState) => {
                    return prevProductsState.map((product) => {
                      const selectedOption = productData.find(
                        (option: any) => option.value === product.id
                      );
                      return {
                        ...product,
                        last_cost: selectedOption?.data?.last_cost || 0,
                      };
                    });
                  });
                });
              },
            })}
            label="Склад"
            options={skladOptions}
          />
          <LabeledSelect
            {...register("schet_id", { valueAsNumber: true })}
            label="Счет"
            options={schetOptions}
          />
          <div className="w-full flex items-center">
            <label htmlFor="date" className="w-40 mr-4">
              Дата
            </label>
            <div>
              <Controller
                name="date"
                control={control}
                render={({ field }) => (
                  <DatePicker
                    locale="ru"
                    renderCustomHeader={({ monthDate }) => (
                      <span className="font-medium font-inter text-sm capitalize">
                        {monthDate.toLocaleDateString("default", {
                          month: "long",
                        })}{" "}
                        {monthDate.getFullYear()}
                      </span>
                    )}
                    selected={field.value}
                    onChange={field.onChange}
                    timeInputLabel="Время:"
                    dateFormat="dd.MM.yyyy HH:mm"
                    showTimeInput
                    className="rounded text-gray-800 py-2 px-3 border border-gray-300 focus:outline-none focus:border-indigo-500"
                  />
                )}
              />
            </div>
          </div>
        </div>

        <div className="w-2/3 flex flex-col space-y-3 mb-4">
          <div className="w-full flex flex-col space-y-3">
            <div className="w-full flex items-center border-b border-gray-200 pb-2">
              <div className="w-2/5 text-sm text-gray-500 font-medium">
                Ингредиент
              </div>
              <div className="w-1/5 text-sm text-gray-500 font-medium">
                Кол-во
              </div>
              <div className="w-1/5 text-sm text-gray-500 font-medium">
                Цена за единицу
              </div>
              <div className="w-1/5 text-sm text-gray-500 font-medium">
                Итого
              </div>
              <div className="w-9 text-right text-sm text-gray-500 font-medium" />
            </div>
            {products.map((ingredient, idx) => (
              <div
                className="w-full flex justify-between items-center"
                key={idx}
              >
                <div className="w-2/5 pr-3 mb-5">
                  <SelectWithSearch
                    className="w-full"
                    options={productOptions}
                    value={ingredient.id}
                    onChange={(value) => {
                      const selectedOption = productOptions.find(
                        (elem) => elem.value === value
                      );
                      handleSelectIngredient(
                        idx,
                        value as string,
                        selectedOption && selectedOption.data
                      );
                    }}
                  />
                </div>
                <div className="w-1/5 pr-3 flex items-center space-x-2 mb-5">
                  <Input
                    type="text"
                    name="quantity"
                    className="grow"
                    value={products[idx].quantity || 0}
                    onInput={(e) => handleChangeQuantity(idx, e)}
                  />
                  <span className="block w-8">{products[idx]?.measure}</span>
                </div>
                <div
                  className={clsx([
                    "w-1/5 pr-3 flex flex-col items-start space-y-1",
                    !products[idx].last_cost && "mb-5",
                  ])}
                >
                  <Input
                    type="text"
                    name="price"
                    value={products[idx].price || 0}
                    className="w-full"
                    onInput={(e) => handleChangePrice(idx, e)}
                  />
                  {!!products[idx].last_cost && (
                    <div className="relative group">
                      <button
                        type="button"
                        onClick={() =>
                          handleChangePrice(
                            idx,
                            null,
                            `${products[idx].last_cost}`
                          )
                        }
                        className="pl-3 flex items-center space-x-1 text-xs font-medium hover:bg-gray-100"
                      >
                        <ClockIcon className="w-3 h-3 text-gray-500" />
                        <span>
                          {formatNumber(
                            products[idx].last_cost as number,
                            true,
                            true
                          )}
                        </span>
                        {(products[idx].last_cost as number) >=
                        products[idx].price ? (
                          <ChevronUpIcon className="text-green-500 w-3 h-3" />
                        ) : (
                          <ChevronDownIcon className="text-red-500 h-3 w-3" />
                        )}
                        {products[idx].last_cost !== undefined &&
                          products[idx].price > 0 && (
                            <div className="absolute left-0 top-6 border border-gray-500 rounded p-2 hidden group-hover:flex flex-col items-start space-y-1.5 w-max bg-white">
                              <span className="text-xs font-normal">
                                Предыдущая цена{" "}
                                {(products[idx].last_cost as number) >=
                                products[idx].price
                                  ? "больше"
                                  : "меньше"}{" "}
                                текущей на{" "}
                                <strong>
                                  {Math.round(
                                    (100 *
                                      Math.abs(
                                        (products[idx].price as number) -
                                          (products[idx].last_cost as number)
                                      )) /
                                      (products[idx].price as number)
                                  )}
                                  %
                                </strong>
                              </span>
                              <span className="text-[11px] text-gray-400 font-normal">
                                Нажмите на цену, чтобы добавить её в поле
                              </span>
                            </div>
                          )}
                      </button>
                    </div>
                  )}
                </div>
                <div className="w-1/5 pr-2 mb-5">
                  <Input
                    type="text"
                    name="total"
                    value={products[idx].total || 0}
                    onInput={(e) => handleChangeTotalPrice(idx, e)}
                  />
                </div>
                <button
                  onClick={() => {
                    setProducts((prevState) =>
                      prevState.filter((_, i) => idx !== i)
                    );
                  }}
                  type="button"
                  className="p-2 rounded-md hover:bg-gray-200 transition duration-100 mb-5"
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>
            ))}
          </div>
          <button
            onClick={() => {
              setProducts((prevState) => [...prevState, defaultProductOption]);
            }}
            type="button"
            className="flex items-center space-x-1 text-indigo-500 hover:text-indigo-700"
          >
            <PlusIcon className="w-4 h-4" />
            <span>Добавить ингредиент</span>
          </button>
          <div className="font-medium text-lg border-t border-gray-300 pt-2 mt-2">
            Итого: {totalPrice}
          </div>
        </div>

        <div className="pt-5 border-t border-gray-200">
          <button
            type="submit"
            className="py-2 px-3 bg-primary disabled:bg-gray-400/60 hover:bg-teal-600 transition duration-300 text-white rounded-md"
            disabled={isLoading}
          >
            Сохранить
          </button>
        </div>
      </form>
    );
};

export default SupplyForm;
