import React, { FC, FormEvent, useCallback } from "react";
import { Select } from "@shared/ui/Select";
import clsx from "clsx";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";

export interface ItemTypeSelectProps {
  className?: string;
}

const options = [
  {
    name: "Все типы",
    value: "",
  },
  {
    name: "Товары",
    value: "tovar",
  },
  {
    name: "Ингредиенты",
    value: "ingredient",
  },
];

const ItemTypeSelect: FC<ItemTypeSelectProps> = ({ className = "" }) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const handleChange = useCallback(
    (e: FormEvent<HTMLSelectElement>) => {
      handleFilterChange({
        [QueryOptions.TYPE]: (e.target as HTMLSelectElement).value,
      });
    },
    [handleFilterChange],
  );

  return (
    <Select
      options={options}
      value={(getFilterValue(QueryOptions.TYPE) as string) || ""}
      onChange={handleChange}
      className={clsx(["w-40", className])}
    />
  );
};

export default ItemTypeSelect;
