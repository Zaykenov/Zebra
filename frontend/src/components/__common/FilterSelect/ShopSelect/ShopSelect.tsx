import React, { FC, useCallback, useEffect, useState } from "react";
import { SelectOption } from "@shared/ui/Select/Select";
import { getAllShops } from "@api/shops";
import { Select } from "antd";
import clsx from "clsx";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";

export interface ShopSelectProps {
  onChange?: (value: number[]) => void;
  className?: string;
}

const ShopSelect: FC<ShopSelectProps> = ({ onChange, className = "" }) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<SelectOption[]>([]);

  useEffect(() => {
    getAllShops().then((res) => {
      setOptions([
        {
          name: "Все заведения",
          value: 0,
        },
        ...res.data.map((shop: any) => ({
          name: shop.name,
          value: shop.id !== 0 ? shop.id : -1,
        })),
      ]);
    });
  }, []);

  const handleChange = useCallback(
    (value: number[]) => {
      onChange
        ? onChange(value)
        : handleFilterChange({
            [QueryOptions.SHOP]: value,
          });
    },
    [handleFilterChange, onChange]
  );

  return (
    <div className="w-40">
      <Select
        mode="multiple"
        allowClear
        value={(getFilterValue(QueryOptions.SHOP) as number[]) || []}
        style={{ width: "100%", flex: 1 }}
        placeholder="Все заведения"
        onChange={handleChange}
        className={clsx(["w-full", className])}
        options={options.map((option) => ({
          label: option.name,
          value: option.value,
        }))}
        filterOption={(input, option) =>
          (option?.label ?? "")
            .trim()
            .toLowerCase()
            .includes(input.trim().toLowerCase())
        }
      />
    </div>
  );
};

export default ShopSelect;
