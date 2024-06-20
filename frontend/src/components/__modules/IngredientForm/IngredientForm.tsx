import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledInput } from "@shared/ui/Input";
import LabeledSelect from "@shared/ui/Select/LabeledSelect";
import { Controller, useForm } from "react-hook-form";
import {
  createIngredient,
  updateIngredient,
  updateMasterIngredient,
} from "@api/ingredient";
import { useRouter } from "next/router";
import { getAllIngredientCategories } from "@api/ingredient-category";
import { IngredientData } from "./types";
import { getAllShops } from "@api/shops";
import { Select } from "antd";
import useMasterRole from "@hooks/useMasterRole";

const measureOptions = [
  {
    name: "кг",
    value: "кг",
  },
  {
    name: "шт.",
    value: "шт.",
  },
  {
    name: "л",
    value: "л",
  },
];

export interface IngredientFormProps {
  isEdit?: boolean;
  data?: any;
}

const IngredientForm: FC<IngredientFormProps> = ({ data, isEdit = false }) => {
  const router = useRouter();

  const isMaster = useMasterRole();

  const [categoryOptions, setCategoryOptions] = useState([]);
  const [shopOptions, setShopOptions] = useState<
    {
      value: number;
      label: string;
    }[]
  >([]);

  const { handleSubmit, register, reset, getValues, control } =
    useForm<IngredientData>({
      defaultValues: {
        name: data?.name || "",
        category: data?.category || "",
        measure: data?.measure || "кг",
        count: data?.count || 0,
        cost: data?.cost || 0,
        sklad: data?.sklad || "sklad1",
      },
    });

  const [selectedShops, setSelectedShops] = useState<number[]>([]);

  useEffect(() => {
    getAllIngredientCategories().then((res) => {
      const categories = res.data.map(
        ({ id, name }: { id: string; name: string }) => ({
          label: name,
          value: parseInt(id),
        })
      );
      setCategoryOptions(categories);
      reset({
        ...getValues(),
        category: categories[0]?.value || 1,
      });
    });
  }, [reset, getValues]);

  useEffect(() => {
    getAllShops().then((res) => {
      const shops = res.data.map(
        ({ id, name }: { id: number; name: string }) => ({
          label: name,
          value: id,
        })
      );
      setShopOptions(shops);
    });
  }, []);

  const [role, setRole] = useState<string>("");

  useEffect(() => {
    const role = localStorage.getItem("zebra.role");
    role && setRole(role);
  }, []);

  const onSubmit = useCallback(
    (submitData: IngredientData) => {
      if (role !== "master") {
        submitData.shop_id =
          selectedShops.length > 0
            ? selectedShops
            : shopOptions.map((shop) => shop.value);
      }

      if (!data)
        createIngredient(submitData).then(() => router.replace("/ingredients"));
      else {
        role === "master"
          ? updateMasterIngredient({
              id: data?.id,
              ...submitData,
            }).then(() => router.replace("/ingredients"))
          : updateIngredient({
              id: data?.id,
              ...submitData,
            }).then(() => router.replace("/ingredients"));
      }
    },
    [data, router, selectedShops, role]
  );

  useEffect(() => {
    if (!data || !shopOptions.length) return;
    reset({ ...data, category: data.category_id });
    setSelectedShops(
      typeof data.shop_id === "number" ? [data.shop_id] : data.shop_id || []
    );
  }, [data, reset, shopOptions]);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col w-1/2">
      <LabeledInput fieldClass="mb-5" {...register("name")} label="Название" />
      <div className="w-full flex items-center pt-2 mb-5">
        <label className="w-40 mr-4">Категория</label>
        <Controller
          control={control}
          {...register("category")}
          render={({ field: { onChange, value } }) => (
            <Select
              allowClear
              disabled={!isMaster && isEdit}
              style={{ width: "100%", flex: 1 }}
              placeholder="Все заведения"
              options={categoryOptions}
              value={value}
              onChange={onChange}
              filterOption={(input, option) =>
                (option?.label ?? "")
                  .trim()
                  .toLowerCase()
                  .includes(input.trim().toLowerCase())
              }
            />
          )}
        />
      </div>
      <LabeledSelect
        fieldClass="mb-5 w-[280px]"
        {...register("measure")}
        options={measureOptions}
        label="Ед. измерения"
        disabled={!isMaster && isEdit}
      />
      {role !== "master" && (
        <div className="w-full flex items-center pt-2 mb-5">
          <label className="w-40 mr-4">Заведения</label>
          <Select
            mode="multiple"
            allowClear
            disabled={!isMaster && isEdit}
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

      {!data && (
        <div className="mb-5 mt-9 flex items-center space-x-2">
          <LabeledInput
            fieldClass="w-[250px]"
            {...register("cost", { valueAsNumber: true })}
            label="Цена за единицу"
            disabled={!isMaster && isEdit}
          />
          <span className="font-medium flex items-center justify-center">
            ₸
          </span>
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

export default IngredientForm;
