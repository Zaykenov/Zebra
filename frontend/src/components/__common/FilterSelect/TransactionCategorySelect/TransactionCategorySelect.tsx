import React, { FC, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";

export interface TransactionCategorySelectProps {
  onChange: (value: number) => void;
  className?: string;
}

const TransactionCategorySelect: FC<TransactionCategorySelectProps> = ({
  onChange,
  className = "",
}) => {
  const [options, setOptions] = useState<SelectOption[]>([
    {
      name: "Все транзакции",
      value: 0,
    },
  ]);

  return (
    <Select
      options={options}
      onChange={(e) =>
        onChange(parseInt((e.target as HTMLSelectElement).value))
      }
      className={clsx(["w-40", className])}
    />
  );
};

export default TransactionCategorySelect;
