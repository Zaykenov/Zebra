import React, { FC } from "react";
import { Select } from "antd";

interface SelectWithSearchProps {
  options: any[];
  value: any;
  onChange: (value: string | number, option?: any) => void;
  className?: string;
  disabled?: boolean;
}

const SelectWithSearch: FC<SelectWithSearchProps> = ({
  options,
  value,
  onChange,
  className,
  disabled = false,
}) => {
  return (
    <Select
      showSearch
      disabled={disabled}
      className={className}
      options={options}
      value={value}
      optionFilterProp="children"
      filterOption={(input, option) =>
        (option?.name ?? "").toLowerCase().includes(input.toLowerCase())
      }
      filterSort={(optionA, optionB) =>
        (optionA?.name ?? "")
          .toLowerCase()
          .localeCompare((optionB?.name ?? "").toLowerCase())
      }
      fieldNames={{ label: "name" }}
      onChange={onChange}
    />
  );
};

export default SelectWithSearch;
