import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import TerminalLayout from "@layouts/TerminalLayout";
import { Input } from "@shared/ui/Input";
import { PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { useRouter } from "next/router";
import { getItems } from "@api/index";
import { createWaste, WasteData } from "@api/wastes";
import { useForm } from "react-hook-form";
import { getAllDishes } from "@api/dishes";
import { Dropdown } from "semantic-ui-react";
import { formatInputValue } from "@utils/formatInputValue";
import "semantic-ui-css/semantic.min.css";

const TerminalWastePage: NextPage = () => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState<
    {
      id: string;
      quantity: number | string;
      details: string;
      type: string;
      measure: string;
    }[]
  >([]);
  const [productOptions, setProductOptions] = useState<
    {
      label: string;
      value: string;
      data: {
        type: string;
        measure: string;
      };
    }[]
  >([]);

  const { handleSubmit, register } = useForm<WasteData>({
    defaultValues: {
      reason: "",
      comment: "",
      items: [],
    },
  });

  useEffect(() => {
    Promise.all([getItems(), getAllDishes()]).then(([resItems, resDishes]) => {
      setProductOptions([
        ...resItems.data.map((item: any) => ({
          label: item.name,
          value: `${item.id}_${item.type}`,
          data: {
            type: item.type,
            measure: item.measure,
          },
        })),
        ...resDishes.data.map((dish: any) => ({
          label: dish.name,
          value: `${dish.tech_cart_id}_techCart`,
          data: {
            type: "techCart",
            measure: dish.measure,
          },
        })),
      ]);
    });
  }, []);

  const handleSelectIngredient = (idx: number, value: string, data?: any) => {
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: value,
          quantity: option.quantity,
          details: option.details,
          type: data.type,
          measure: data.measure,
        };
      })
    );
  };

  const handleChangeQuantity = (
    idx: number,
    e: React.FormEvent<HTMLInputElement>
  ) => {
    const { inputValue } = formatInputValue(
      (e.target as HTMLInputElement).value
    );
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: option.id,
          type: option.type,
          measure: option.measure,
          details: option.details,
          quantity: inputValue,
        };
      })
    );
  };

  const handleChangeDetails = (
    idx: number,
    e: React.FormEvent<HTMLInputElement>
  ) => {
    const inputValue = (e.target as HTMLInputElement).value;
    setProducts((prevState) =>
      prevState.map((option, i) => {
        if (idx !== i) return option;
        return {
          id: option.id,
          type: option.type,
          measure: option.measure,
          quantity: option.quantity,
          details: inputValue,
        };
      })
    );
  };

  const onSubmit = async (submitData: WasteData) => {
    setLoading(true);
    try {
      await createWaste({
        ...submitData,
        items: products.map((product) => ({
          item_id: parseInt(product.id),
          type: product.id.split("_")[1],
          quantity: parseFloat(product.quantity as string),
          details: product.details,
        })),
      });
      router.reload();
    } catch (error) {
      setLoading(false);
    }
  };

  return (
    <TerminalLayout>
      <div className="w-full min-h-full py-10 flex flex-col items-center overflow-auto">
        <div className="w-full max-w-3xl p-4 rounded bg-gray-100 shadow-2xl border border-gray-300 flex flex-col">
          <div className="text-lg font-medium mb-4">Добавить списание</div>
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="h-full flex flex-col justify-between"
          >
            <div className="flex flex-col space-y-5 mb-5">
              <Input
                {...register("reason")}
                placeholder="Причина"
                className="w-2/3"
              />
              <Input
                {...register("comment")}
                placeholder="Комментарий"
                className="w-2/3"
              />

              <div className="flex flex-col space-y-3 mb-4">
                {products.length > 0 && (
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
                      <div className="w-9 text-sm text-gray-500 font-medium"></div>
                    </div>
                    {products.map((item, idx) => {
                      return (
                        <div
                          className="w-full flex justify-between items-center"
                          key={idx}
                        >
                          <div className="w-1/2 pr-4">
                            <Dropdown
                              placeholder="Выберите позицию"
                              fluid
                              search
                              selection
                              value={item.id}
                              onChange={(_, data) => {
                                const selectedOption = productOptions.find(
                                  (elem) => elem.value === data.value
                                );
                                handleSelectIngredient(
                                  idx,
                                  data.value as string,
                                  selectedOption && selectedOption.data
                                );
                              }}
                              options={productOptions.map((option) => ({
                                key: option.value,
                                value: option.value,
                                text: option.label,
                              }))}
                            />
                          </div>
                          <div className="flex items-center w-1/4 pr-4 space-x-2">
                            <Input
                              type="text"
                              name="quantity"
                              className="w-full"
                              value={products[idx].quantity}
                              onInput={(e) => handleChangeQuantity(idx, e)}
                            />
                            <span className="w-8">{item.measure}</span>
                          </div>
                          <div className="w-1/4 pr-4">
                            <Input
                              type="text"
                              name="details"
                              className="w-full"
                              value={products[idx].details}
                              onChange={(e) => handleChangeDetails(idx, e)}
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
                )}
                <button
                  onClick={() => {
                    setProducts((prevState) => [
                      ...prevState,
                      {
                        id: "",
                        quantity: 0,
                        type: "",
                        details: "",
                        measure: "",
                      },
                    ]);
                  }}
                  type="button"
                  className="flex items-center space-x-1 text-indigo-500 hover:text-indigo-700"
                >
                  <PlusIcon className="w-4 h-4" /> <span>Добавить еще</span>
                </button>
              </div>
            </div>

            <div className="pt-5 border-t border-gray-200">
              <div className="flex items-center justify-between">
                <button
                  disabled={loading}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    router.back();
                  }}
                  className="px-4 py-2 bg-transparent hover:bg-gray-300 rounded text-gray-500 hover:text-gray-900"
                >
                  Отмена
                </button>
                <button
                  disabled={loading}
                  type="submit"
                  className="disabled:bg-primary/50 px-8 py-2 bg-primary hover:opacity-80 rounded text-white font-medium"
                >
                  Сохранить
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </TerminalLayout>
  );
};

export default TerminalWastePage;
