import React, { FC, FormEvent, useCallback, useEffect, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllAccounts } from "@api/accounts";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";

export interface SchetSelectProps {
  className?: string;
}

const SchetSelect: FC<SchetSelectProps> = ({ className = "" }) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<SelectOption[]>([]);

  useEffect(() => {
    getAllAccounts().then((res) => {
      setOptions([
        {
          name: "Все счета",
          value: 0,
        },
        ...res.data.map((schet: any) => ({
          name: schet.name,
          value: schet.id !== 0 ? schet.id : -1,
        })),
      ]);
    });
  }, []);

  const handleChange = useCallback(
    (e: FormEvent<HTMLSelectElement>) => {
      handleFilterChange({
        [QueryOptions.SCHET]: parseInt((e.target as HTMLSelectElement).value),
      });
    },
    [handleFilterChange],
  );

  return (
    <Select
      options={options}
      onChange={handleChange}
      value={(getFilterValue(QueryOptions.SCHET) as number) || 0}
      className={clsx(["w-40", className])}
    />
  );
};

export default SchetSelect;
