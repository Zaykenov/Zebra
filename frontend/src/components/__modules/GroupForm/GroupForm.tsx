import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledInput } from "@shared/ui/Input";
import LabeledSelect from "@shared/ui/Select/LabeledSelect";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import {
  createInventoryGroup,
  InventoryGroupData,
  updateInventoryGroup,
} from "@api/groups";
import { getAllSklads } from "@api/sklad";
import { getAllIngredients } from "@api/ingredient";
import { getAllMenuItems } from "@api/menu-items";
import { Select } from "antd";
import { QueryOptions } from "@api/index";

const measureOptions = [
  {
    name: "шт.",
    value: "шт.",
  },
  {
    name: "кг",
    value: "кг",
  },
  {
    name: "л",
    value: "л",
  },
];

const typeOptions = [
  {
    name: "Ингредиент",
    value: "ingredient",
  },
  {
    name: "Товар",
    value: "tovar",
  },
];

export interface GroupFormProps {
  data?: InventoryGroupData | null;
}

const GroupForm: FC<GroupFormProps> = ({ data }) => {
  const router = useRouter();

  const { handleSubmit, register, reset, setValue, watch } =
    useForm<InventoryGroupData>();

  const [skladOptions, setSkladOptions] = useState<
    {
      name: string;
      value: string;
    }[]
  >([]);

  const [itemOptions, setItemOptions] = useState<
    {
      label: string;
      value: string;
      data: {
        type: string;
        measure: string;
      };
    }[]
  >([]);

  const [selectedItems, setSelectedItems] = useState<number[]>([]);

  const getAllItems = useCallback(async (id: number) => {
    await getAllIngredients({ [QueryOptions.SKLAD]: id }).then(
      (ingredientsRes) => {
        getAllMenuItems({ [QueryOptions.SKLAD]: id }).then((tovarsRes) => {
          setItemOptions([
            ...ingredientsRes.data.map((ingredient: any) => ({
              value: ingredient.ingredient_id,
              label: ingredient.name,
              data: {
                type: "ingredient",
                measure: ingredient.measure,
              },
            })),
            ...tovarsRes.data.map((tovar: any) => ({
              value: tovar.tovar_id,
              label: tovar.name,
              data: {
                type: "tovar",
                measure: tovar.measure,
              },
            })),
          ]);
        });
      }
    );
  }, []);

  useEffect(() => {
    getAllSklads().then(async (res) => {
      setSkladOptions(
        res.data.map((sklad: any) => ({
          name: sklad.name,
          value: sklad.id,
        }))
      );
      if (res.data.length === 0) return;
      setValue("sklad_id", res.data[0].id);
      await getAllItems(res.data[0].id);
    });
  }, [setValue]);

  useEffect(() => {
    if (!data) return;
    reset({ ...data });
    setSelectedItems(data.items.map((item) => item.item_id));
  }, [data, reset]);

  const onSubmit = useCallback(
    (submitData: InventoryGroupData) => {
      const postData = {
        ...submitData,
        items: selectedItems.map((item_id) => ({ item_id })),
      };
      if (!data) {
        createInventoryGroup(postData).then(() => router.replace("/groups"));
      } else {
        updateInventoryGroup({
          id: data.id,
          ...postData,
        }).then(() => router.replace("/groups"));
      }
    },
    [data, selectedItems, router]
  );

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-1/2 space-y-5"
    >
      <LabeledInput {...register("name")} label="Название" />
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
      <LabeledSelect
        {...register("measure")}
        options={measureOptions}
        label="Ед. измерения"
      />
      <LabeledSelect
        {...register("type")}
        options={typeOptions}
        label="Группировка по"
      />

      <div className="w-full flex pt-2">
        <label className="w-40 mr-4">Группировка</label>
        <Select
          mode="multiple"
          allowClear
          style={{ width: "100%" }}
          placeholder="Продукты для группировки"
          value={selectedItems}
          onChange={(value) => {
            setSelectedItems(value);
          }}
          options={itemOptions.filter(
            (option) =>
              option.data.type === watch("type") &&
              option.data.measure === watch("measure")
          )}
          filterOption={(input, option) =>
            (option?.label ?? "")
              .trim()
              .toLowerCase()
              .includes(input.trim().toLowerCase())
          }
        />
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

export default GroupForm;
