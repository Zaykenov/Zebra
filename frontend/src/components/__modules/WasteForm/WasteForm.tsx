import React, { FC, useCallback, useEffect, useMemo, useState } from "react";
import { Input, LabeledInput } from "@shared/ui/Input";
import { Controller, useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { LabeledSelect, Select } from "@shared/ui/Select";
import { PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { getItems, QueryOptions } from "@api/index";
import { getAllSklads } from "@api/sklad";
import {
  createWaste,
  updateWaste,
  WasteData,
  WasteItemData,
} from "@api/wastes";
import { getAllDishes } from "@api/dishes";
import DatePicker from "react-datepicker";
import SelectWithSearch from "@shared/ui/SelectWithSearch";
import { getAllMenuItems } from "@api/menu-items";
import { getAllIngredients } from "@api/ingredient";

export interface WasteFormProps {
  data?: any;
  isEdit?: boolean;
}

const WasteForm: FC<WasteFormProps> = ({ data, isEdit = false }) => {
  const router = useRouter();

  const [dataLoaded, setDataLoaded] = useState(false);
  const [products, setProducts] = useState<WasteItemData[]>([]);
  const [productOptions, setProductOptions] = useState<any[]>([]);
  const [skladOptions, setSkladOptions] = useState<any[]>([]);
  const { handleSubmit, register, reset, control, getValues } =
    useForm<WasteData>({
      defaultValues: {
        sklad_id: 1,
        reason: "",
        comment: "",
        time: new Date().toISOString(),
        date: new Date(),
      },
    });

  useEffect(() => {
    if (!data || !productOptions || dataLoaded) return;
    setDataLoaded(true);
    reset({
      sklad_id: data.sklad_id,
      reason: data.reason,
      comment: data.comment,
      time: data.time,
      date: new Date(data.time),
    });
    setProducts(
      data.items.map((dataItem: any) => ({
        item_id: `${dataItem.item_id}_${dataItem.type}_${dataItem.measure}`,
        type: dataItem.type,
        measure: dataItem.measure,
        quantity: dataItem.quantity,
        details: dataItem.details,
      }))
    );
  }, [data, reset, productOptions, dataLoaded]);

  const getAllItems = useCallback(async (id: number) => {
    const [resProducts, resIngredients, resDishes] = await Promise.all([
      getAllMenuItems({ [QueryOptions.SKLAD]: id }),
      getAllIngredients({ [QueryOptions.SKLAD]: id }),
      getAllDishes({ [QueryOptions.SKLAD]: id }),
    ]);
    setProductOptions([
      ...resProducts.data.map((item: any) => ({
        name: item.name,
        value: `${item.tovar_id}_tovar_${item.measure}`,
      })),
      ...resIngredients.data.map((ingredient: any) => ({
        name: ingredient.name,
        value: `${ingredient.ingredient_id}_ingredient_${ingredient.measure}`,
      })),
      ...resDishes.data.map((dish: any) => ({
        name: dish.name,
        value: `${dish.tech_cart_id}_techCart_${dish.measure}`,
      })),
    ]);
  }, []);

  useEffect(() => {
    getAllSklads().then(async (resSklads) => {
      setSkladOptions(
        resSklads.data.map((item: any) => ({
          name: item.name,
          value: parseInt(item.id),
        }))
      );
      reset({
        ...getValues(),
        sklad_id: resSklads.data[0].id,
      });
    });
  }, []);

  useEffect(() => {
    if (!skladOptions.length || (isEdit && !data)) return;
    const selectedSklad = isEdit ? data.sklad_id : skladOptions[0].value;
    getAllItems(selectedSklad).then();
  }, [isEdit, data, getAllItems, skladOptions]);

  const defaultProductOption = useMemo(
    () => ({
      name: productOptions[0]?.name || "",
      value: productOptions[0]?.value || "",
    }),
    [productOptions]
  );

  const onSubmit = useCallback(
    (submitData: WasteData) => {
      const time = submitData.date?.toISOString() || new Date().toISOString();
      const items = products.map((product) => ({
        item_id: parseInt(product.item_id as string),
        type: product.type,
        quantity:
          typeof product.quantity === "string"
            ? parseFloat(product.quantity.replace(",", ".").replace(" ", ""))
            : product.quantity,
        details: product.details,
      }));
      if (!data) {
        createWaste({ ...submitData, time, items }).then(() =>
          router.replace("/waste")
        );
      } else {
        updateWaste({ ...submitData, time, items, id: data.id }).then(() =>
          router.push("/waste")
        );
      }
    },
    [products, data, router]
  );

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col space-y-5">
      <div className="w-1/2 flex flex-col space-y-5">
        <LabeledSelect
          {...register("sklad_id", {
            valueAsNumber: true,
            onChange: async (event) => {
              await getAllItems(event.target.value);
            },
          })}
          label="Склад"
          options={skladOptions}
        />
        <LabeledInput {...register("reason")} label="Причина" />
        <LabeledInput {...register("comment")} label="Комментарий" />

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
            <div className="w-1/2 text-sm text-gray-500 font-medium">
              Что списывается
            </div>
            <div className="w-1/4 text-sm text-gray-500 font-medium">
              Кол-во
            </div>
            <div className="w-1/4 text-sm text-gray-500 font-medium">
              Детали
            </div>
            <div className="w-9 text-right text-sm text-gray-500 font-medium"></div>
          </div>
          {products.map((item, idx) => {
            return (
              <div
                className="w-full flex justify-between items-center"
                key={idx}
              >
                <div className="w-1/2 pr-4">
                  <SelectWithSearch
                    className="w-full"
                    options={productOptions}
                    value={item.item_id}
                    // @ts-ignore
                    onChange={(value: string) => {
                      setProducts((prevState) =>
                        prevState.map((option, i) => {
                          if (idx !== i) return option;
                          return {
                            item_id: value,
                            quantity: option.quantity,
                            type: value.split("_")[1],
                            measure: value.split("_")[2],
                            details: option.details,
                          };
                        })
                      );
                    }}
                  />
                </div>
                <div className="flex items-center w-1/4 pr-4 space-x-2">
                  <Input
                    type="text"
                    name="quantity"
                    className="w-full"
                    value={products[idx].quantity}
                    onInput={(e) => {
                      setProducts((prevState) =>
                        prevState.map((option, i) => {
                          if (idx !== i) return option;
                          return {
                            ...option,
                            quantity: (e.target as HTMLInputElement).value,
                          };
                        })
                      );
                    }}
                  />
                  <span className="w-8">{item.measure}</span>
                </div>
                <div className="w-1/4 pr-4">
                  <Input
                    type="text"
                    name="details"
                    className="w-full"
                    value={products[idx].details}
                    onChange={(e) => {
                      setProducts((prevState) =>
                        prevState.map((option, i) => {
                          if (idx !== i) return option;
                          return {
                            ...option,
                            details: (e.target as HTMLInputElement).value,
                          };
                        })
                      );
                    }}
                  />
                </div>
                <button
                  onClick={() => {
                    setProducts((prevState) =>
                      prevState.filter((_, i) => idx !== i)
                    );
                  }}
                  type="button"
                  className="p-2 rounded-md hover:bg-gray-200 transition duration-100"
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>
            );
          })}
        </div>
        <button
          onClick={() => {
            setProducts((prevState) => [
              ...prevState,
              {
                item_id: "",
                quantity: 0,
                type: defaultProductOption.value.split("_")[1],
                measure: defaultProductOption.value.split("_")[2],
                details: "",
              },
            ]);
          }}
          type="button"
          className="flex items-center space-x-1 text-indigo-500 hover:text-indigo-700"
        >
          <PlusIcon className="w-4 h-4" /> <span>Добавить еще</span>
        </button>
      </div>

      <div className="pt-5 border-t border-gray-200">
        <button
          type="submit"
          className="py-2 px-3 bg-primary hover:bg-teal-600 transition duration-300 text-white rounded-md"
        >
          Сохранить
        </button>
      </div>
    </form>
  );
};

export default WasteForm;
