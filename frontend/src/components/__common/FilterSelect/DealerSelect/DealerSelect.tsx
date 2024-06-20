import React, { FC, FormEvent, useCallback, useEffect, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllAccounts } from "@api/accounts";
import { getAllSuppliers } from "@api/suppliers";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";

export interface DealerSelectProps {
  className?: string;
}

const DealerSelect: FC<DealerSelectProps> = ({ className = "" }) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<SelectOption[]>([]);

  const handleChange = useCallback(
    (e: FormEvent<HTMLSelectElement>) => {
      handleFilterChange({
        [QueryOptions.DEALER]: parseInt((e.target as HTMLSelectElement).value),
      });
    },
    [handleFilterChange],
  );

  useEffect(() => {
    getAllSuppliers().then((res) => {
      setOptions([
        {
          name: "Все поставщики",
          value: 0,
        },
        ...res.data.map((dealer: any) => ({
          name: dealer.name,
          value: dealer.id,
        })),
      ]);
    });
  }, []);

  return (
    <Select
      options={options}
      onChange={handleChange}
      value={(getFilterValue(QueryOptions.DEALER) as number) || 0}
      className={clsx(["w-40", className])}
    />
  );
};

export default DealerSelect;
