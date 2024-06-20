import React, { FC, useCallback, useEffect, useState } from "react";
import { Select } from "antd";
import { getAllIngredientCategories } from "@api/ingredient-category";
import { getAllProductCategories } from "@api/product-categories";
import { useFilter } from "@context/index";
import { QueryOptions } from "@api/index";

export interface ItemsCategorySelectProps {
  type?: "product" | "ingredient" | "all";
  className?: string;
}

export type GroupedSelectOption = {
  label: string;
  value?: number | string;
  options?: {
    label: string;
    value: number | string;
  }[];
};

const ItemsCategorySelect: FC<ItemsCategorySelectProps> = ({
  type = "all",
  className,
}) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<GroupedSelectOption[]>([
    {
      label: "Все категории",
      value: 0,
    },
  ]);

  const handleChange = useCallback(
    (value: number) => {
      handleFilterChange({ [QueryOptions.CATEGORY]: value });
    },
    [handleFilterChange],
  );

  useEffect(() => {
    type === "all" &&
      getAllIngredientCategories().then((resIngr) => {
        getAllProductCategories().then((resProd) => {
          setOptions((prevState) => [
            ...prevState,
            {
              label: "Товары и тех. карты",
              options: resProd.data.map((category: any) => ({
                label: category.name,
                value: category.id,
              })),
            },
            {
              label: "Ингредиенты",
              options: resIngr.data.map((category: any) => ({
                label: category.name,
                value: category.id,
              })),
            },
          ]);
        });
      });

    type === "product" &&
      getAllProductCategories().then((resProd) => {
        setOptions((prevState) => [
          ...prevState,
          {
            label: "Товары и тех. карты",
            options: resProd.data.map((category: any) => ({
              label: category.name,
              value: category.id,
            })),
          },
        ]);
      });

    type === "ingredient" &&
      getAllIngredientCategories().then((resIngr) => {
        setOptions((prevState) => [
          ...prevState,
          {
            label: "Ингредиенты",
            options: resIngr.data.map((category: any) => ({
              label: category.name,
              value: category.id,
            })),
          },
        ]);
      });
  }, []);

  return (
    <Select
      showSearch
      value={(getFilterValue(QueryOptions.CATEGORY) as number) || 0}
      style={{ width: 200 }}
      onChange={handleChange}
      options={options}
      className={className}
      filterOption={(input, option) =>
        (option?.label ?? "").toLowerCase().includes(input.toLowerCase())
      }
    />
  );
};

export default ItemsCategorySelect;
