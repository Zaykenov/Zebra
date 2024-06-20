import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledSelect } from "@shared/ui/Select";
import { Input } from "@shared/ui/Input";
import { ClockIcon, PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { useRouter } from "next/router";
import { getItems } from "@api/index";
import { createSupplyAsWorker, SupplyData, updateSupply } from "@api/supplies";
import { useForm } from "react-hook-form";
import { getAllAccounts } from "@api/accounts";
import axios, { AxiosError } from "axios";
import { Dropdown } from "semantic-ui-react";
import "semantic-ui-css/semantic.min.css";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage";
import { formatInputValue } from "@utils/formatInputValue";
import clsx from "clsx";
import { formatNumber } from "@utils/formatNumber";
import { Spin } from "antd";

const draftKey = "zebra.cache.supplyDraft";

interface TerminalSupplyFormProps {
  data?: any;
}

const TerminalSupplyForm: FC<TerminalSupplyFormProps> = ({ data }) => {
  const router = useRouter();

  const [loading, setLoading] = useState(false);

  const [schetOptions, setSchetOptions] = useState<any[]>([]);
  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const [totalPrice, setTotalPrice] = useState<number>(0);

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
  const [productOptions, setProductOptions] = useState<
    {
      label: string;
      value: string;
      data: {
        type: string;
        measure: string;
        last_cost: number;
      };
    }[]
  >([]);

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

  const { handleSubmit, register, reset, getValues, setValue } =
    useForm<SupplyData>({
      defaultValues: {
        schet_id: 1,
        items: [],
      },
    });

  useEffect(() => {
    if (!data || !productOptions || !schetOptions) return;
    reset({
      sklad_id: data.sklad_id,
      schet_id: data.schet_id,
      items: data.items,
      time: data.time,
      dealer_id: data.dealer_id,
    });
    setProducts(
      data.items.map((item: any) => ({
        id: `${item.item_id}_${item.type}`,
        quantity: item.quantity,
        type: item.type,
        measure: item.measurement,
        price: item.cost,
        total: item.cost * item.quantity,
      }))
    );
  }, [data, reset, productOptions, schetOptions]);

  useEffect(() => {
    getAllAccounts().then((accountsRes) => {
      setSchetOptions(
        accountsRes.data.map((item: any) => ({
          name: item.name,
          value: parseInt(item.id),
        }))
      );
      reset({
        schet_id: accountsRes.data[0]?.id,
      });
    });
    getItems().then((res) => {
      setProductOptions(
        res.data.map((item: any) => ({
          label: item.name,
          value: `${item.id}_${item.type}`,
          data: {
            type: item.type,
            measure: item.measure,
            last_cost: item.last_postavka_cost,
          },
        }))
      );
      if (data) return;

      res.data.length > 0 &&
        setProducts([
          {
            id: `${res.data[0].id}_${res.data[0].type}`,
            quantity: 0,
            price: 0,
            total: 0,
            type: res.data[0].type,
            measure: res.data[0].measure,
            last_cost: res.data[0].last_postavka_cost,
          },
        ]);
      const draftSupplyString = localStorage
        ? localStorage.getItem(draftKey)
        : null;
      if (draftSupplyString) {
        const draftSupply = JSON.parse(draftSupplyString);
        setProducts(draftSupply.items);
        setValue("schet_id", draftSupply.schet_id);
      }
    });
  }, [data]);

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

  const handleAddProduct = () => {
    if (products.length > 0) {
      setProducts((prevState) => [
        ...prevState,
        {
          id: "",
          quantity: 0,
          type: "",
          price: 0,
          measure: products[0].measure,
          total: 0,
        },
      ]);
    } else {
      setProducts([
        {
          id: "",
          quantity: 0,
          type: "",
          price: 0,
          measure: "кг.",
          total: 0,
        },
      ]);
    }
  };

  const onSubmit = useCallback(
    (submitData: any) => {
      setLoading(true);
      data
        ? updateSupply({
            id: data.id,
            ...submitData,
            items: products.map((product) => ({
              item_id: parseInt(product.id),
              type: product.type,
              quantity: parseFloat(product.quantity as string),
              cost: parseFloat(product.price as string),
            })),
          })
            .then(() => {
              localStorage && localStorage.removeItem(draftKey);
              router.reload();
            })
            .catch((e: Error | AxiosError) => {
              setLoading(false);
              if (axios.isAxiosError(e)) {
                if (e.response?.status === 666) {
                  showAlertMessage("НАЧНИТЕ СМЕНУ!", AlertMessageType.ERROR);
                  return router.push("/terminal/shift");
                } else showAlertMessage(e.message, AlertMessageType.ERROR);
              } else showAlertMessage(e.message, AlertMessageType.ERROR);
            })
        : createSupplyAsWorker({
            ...submitData,
            items: products.map((product) => ({
              item_id: parseInt(product.id),
              type: product.type,
              quantity: parseFloat(product.quantity as string),
              cost: parseFloat(product.price as string),
            })),
          })
            .then(() => {
              localStorage && localStorage.removeItem(draftKey);
              router.reload();
            })
            .catch((e: Error | AxiosError) => {
              setLoading(false);
              if (axios.isAxiosError(e)) {
                if (e.response?.status === 666) {
                  showAlertMessage("НАЧНИТЕ СМЕНУ!", AlertMessageType.ERROR);
                  return router.push("/terminal/shift");
                } else showAlertMessage(e.message, AlertMessageType.ERROR);
              } else showAlertMessage(e.message, AlertMessageType.ERROR);
            });
    },
    [products, data]
  );

  const onDraftSubmit = useCallback(() => {
    const dataString = JSON.stringify({
      schet_id: getValues("schet_id"),
      items: products,
    });
    localStorage && localStorage.setItem(draftKey, dataString);
    return router.push("/terminal/order");
  }, [products, router, getValues]);

  return (
    <div className="w-full h-full max-w-3xl p-4 rounded bg-gray-100 shadow-2xl border border-gray-300 flex flex-col justify-between">
      <div className="w-full h-full">
        <div className="text-lg font-medium mb-4">
          {data ? "Редактировать" : "Добавить"} поставку
        </div>
        <LabeledSelect
          {...register("schet_id", { valueAsNumber: true })}
          label="Счет"
          options={schetOptions}
          fieldClass="mb-5"
        />
        {alertMessage && (
          <AlertMessage
            message={alertMessage.message}
            type={alertMessage.type}
            onClose={hideAlertMessage}
          />
        )}
        {!productOptions.length ? (
          <div className="w-full flex items-center justify-center">
            <Spin />
          </div>
        ) : (
          <div className="flex flex-col space-y-3 mb-4">
            {products.length > 0 && (
              <div className="w-full flex flex-col space-y-3">
                <div className="w-full flex items-center border-b border-gray-200 pb-2">
                  <div className="w-1/2 text-sm text-gray-500 font-medium">
                    Ингредиент
                  </div>
                  <div className="w-1/6 text-sm text-gray-500 font-medium">
                    Кол-во
                  </div>
                  <div className="w-1/6 pr-3 text-sm text-gray-500 font-medium">
                    Цена за единицу
                  </div>
                  <div className="w-1/6 text-sm text-gray-500 font-medium">
                    Итого
                  </div>
                  <div className="w-9 text-sm text-gray-500 font-medium" />
                </div>
                {products.map((ingredient, idx) => (
                  <div className="w-full flex items-center" key={idx}>
                    <div className="w-1/2 pr-3 mb-5">
                      <Dropdown
                        placeholder="Выберите позицию"
                        fluid
                        search
                        selection
                        value={ingredient.id}
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
                    <div className="w-1/6 flex items-center space-x-1 pr-3 mb-5">
                      <Input
                        className="flex-1"
                        type="text"
                        name="quantity"
                        value={products[idx].quantity}
                        onInput={(e) => handleChangeQuantity(idx, e)}
                      />{" "}
                      <span className="block w-6">{products[idx].measure}</span>
                    </div>
                    <div
                      className={clsx([
                        "w-1/5 pr-3 flex flex-col space-y-1",
                        !products[idx].last_cost && "mb-5",
                      ])}
                    >
                      <Input
                        type="text"
                        name="price"
                        value={products[idx].price || 0}
                        onInput={(e) => handleChangePrice(idx, e)}
                      />
                      {!!products[idx].last_cost && (
                        <span className="pl-3 flex items-center space-x-1 text-xs font-medium">
                          <ClockIcon className="w-3 h-3 text-gray-500" />
                          <span>
                            {formatNumber(
                              products[idx].last_cost as number,
                              true,
                              true
                            )}
                          </span>
                        </span>
                      )}
                    </div>
                    <div className="w-1/6 flex items-center pr-1 mb-5">
                      <Input
                        type="text"
                        name="total"
                        className=""
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
                      className="p-2 rounded-md hover:bg-gray-200 transition duration-100"
                    >
                      <XMarkIcon className="w-5 h-5" />
                    </button>
                  </div>
                ))}
              </div>
            )}
            <button
              onClick={handleAddProduct}
              type="button"
              className="flex items-center space-x-1 text-indigo-500 hover:text-indigo-700"
            >
              <PlusIcon className="w-4 h-4" /> <span>Добавить ингредиент</span>
            </button>
            <div className="font-medium text-lg border-t border-gray-300 pt-2 mt-2">
              Итого:{" "}
              <span className="font-bold ml-2">
                {formatNumber(totalPrice, true, true)}
              </span>
            </div>
          </div>
        )}
      </div>
      <div className="flex items-center justify-between">
        <button
          disabled={loading}
          onClick={() => {
            router.back();
          }}
          className="px-4 py-2 bg-transparent hover:bg-gray-300 rounded text-gray-500 hover:text-gray-900"
        >
          Отмена
        </button>
        <div className="flex items-center space-x-3">
          {!data && (
            <button
              onClick={onDraftSubmit}
              className="px-8 py-2 bg-gray-300 hover:bg-gray-400 rounded text-gray-600 hover:text-white font-medium"
            >
              Черновик
            </button>
          )}
          <button
            disabled={loading}
            onClick={handleSubmit(onSubmit)}
            className="disabled:bg-primary/50 px-8 py-2 bg-primary hover:opacity-80 rounded text-white font-medium"
          >
            Добавить
          </button>
        </div>
      </div>
    </div>
  );
};

export default TerminalSupplyForm;
