import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledInput } from "@shared/ui/Input";
import LabeledSelect from "@shared/ui/Select/LabeledSelect";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { AccountData } from "./types";
import { createAccount, updateAccount } from "@api/accounts";
import clsx from "clsx";
import { Select } from "antd";
import { getAllShops } from "@api/shops";

export interface AccountFormProps {
  data?: any;
  width?: string;
  type?: "cash" | "card";
  onCreate?: (id: number) => void;
}

const typeOptions = [
  {
    name: "Безналичный счет",
    value: "Безналичный счет",
  },
  {
    name: "Банковская карта",
    value: "Банковская карта",
  },
  {
    name: "Наличные",
    value: "Наличные",
  },
];

const AccountForm: FC<AccountFormProps> = ({
  data,
  width = "w-1/2",
  type,
  onCreate,
}) => {
  const router = useRouter();

  const [shopOptions, setShopOptions] = useState<
    {
      value: number;
      label: string;
    }[]
  >([]);
  const [selectedShops, setSelectedShops] = useState<number[]>([]);

  const { handleSubmit, register, reset } = useForm<AccountData>({
    defaultValues: {
      name: data?.name || "",
      currency: data?.currency || "tenge",
      type: data?.type || "Безналичный счет",
      shops: [],
    },
  });

  useEffect(() => {
    getAllShops().then((res) => {
      const shops = res.data.map(
        ({ id, name }: { id: number; name: string }) => ({
          label: name,
          value: id,
        }),
      );
      setShopOptions(shops);
    });
  }, []);

  useEffect(() => {
    reset({
      name: "",
      currency: "tenge",
      type: type === "cash" ? "Наличный счет" : "Безналичный счет",
    });
  }, [type]);

  useEffect(() => {
    reset(data);
  }, [data, reset]);

  const onSubmit = useCallback(
    (submitData: AccountData) => {
      const processedData = { ...submitData, shops: selectedShops };
      if (!data)
        createAccount(processedData).then((res) => {
          !!onCreate ? onCreate(res.data.id) : router.replace("/accounts");
        });
      else {
        updateAccount({
          id: data.id,
          ...processedData,
        }).then(() => router.replace("/accounts"));
      }
    },
    [data, router, selectedShops],
  );

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className={clsx(["flex flex-col space-y-5", width])}
    >
      <LabeledInput {...register("name")} label="Название" />
      {!type && (
        <LabeledSelect
          {...register("type")}
          options={typeOptions}
          label="Категория"
        />
      )}
      <LabeledInput
        {...register("start_balance", { valueAsNumber: true })}
        label="Начальный баланс"
      />
      {!data && (
        <div className="w-full flex items-center pt-2">
          <label className="w-40 mr-4">Заведения</label>
          <Select
            mode="multiple"
            allowClear
            style={{ width: "100%", flex: 1 }}
            placeholder="Все заведения"
            value={selectedShops}
            onChange={(value) => {
              setSelectedShops(value);
            }}
            options={shopOptions}
            filterOption={(input, option) =>
              (option?.label ?? "")
                .trim()
                .toLowerCase()
                .includes(input.trim().toLowerCase())
            }
          />
        </div>
      )}
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

export default AccountForm;
