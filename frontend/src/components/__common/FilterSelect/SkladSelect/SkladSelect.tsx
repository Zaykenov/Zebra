import React, { FC, FormEvent, useCallback, useEffect, useState } from "react";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllSklads } from "@api/sklad";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";
import { Select } from "antd";

export interface SkladSelectProps {
  className?: string;
  defaultOption?: boolean;
}

const SkladSelect: FC<SkladSelectProps> = ({
  className = "",
  defaultOption = true,
}) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<SelectOption[]>([]);

  useEffect(() => {
    if (defaultOption === undefined) return;
    const options: SelectOption[] = [];
    defaultOption &&
      options.push({
        name: "Все склады",
        value: 0,
      });
    getAllSklads().then((res) => {
      setOptions([
        ...options,
        ...res.data.map((sklad: any) => ({
          name: sklad.name,
          value: sklad.id !== 0 ? sklad.id : -1,
        })),
      ]);
    });
  }, [defaultOption]);

  const handleChange = useCallback(
    (value: number[]) => {
      handleFilterChange({
        [QueryOptions.SKLAD]: value,
      });
    },
    [handleFilterChange]
  );

  return (
    <div className="w-40">
      <Select
        mode="multiple"
        allowClear
        value={(getFilterValue(QueryOptions.SKLAD) as number[]) || []}
        style={{ width: "100%", flex: 1 }}
        placeholder="Все склады"
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

export default SkladSelect;
