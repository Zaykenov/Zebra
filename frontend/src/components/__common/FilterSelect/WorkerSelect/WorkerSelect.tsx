import React, { FC, FormEvent, useCallback, useEffect, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import { getAllWorkers } from "@api/workers";
import clsx from "clsx";
import { useFilter } from "@context/filter.context";
import { QueryOptions } from "@api/index";

export interface WorkerSelectProps {
  className?: string;
}

const WorkerSelect: FC<WorkerSelectProps> = ({ className = "" }) => {
  const { handleFilterChange, getFilterValue } = useFilter();

  const [options, setOptions] = useState<SelectOption[]>([]);

  useEffect(() => {
    getAllWorkers().then((res) => {
      setOptions([
        {
          name: "Все кассиры",
          value: 0,
        },
        ...res.data.map((worker: any) => ({
          name: worker.name,
          value: worker.id,
        })),
      ]);
    });
  }, []);

  const handleChange = useCallback(
    (e: FormEvent<HTMLSelectElement>) => {
      handleFilterChange({
        [QueryOptions.WORKER]: parseInt((e.target as HTMLSelectElement).value),
      });
    },
    [handleFilterChange],
  );

  return (
    <Select
      options={options}
      onChange={handleChange}
      value={(getFilterValue(QueryOptions.WORKER) as number) || 0}
      className={clsx(["w-40", className])}
    />
  );
};

export default WorkerSelect;
