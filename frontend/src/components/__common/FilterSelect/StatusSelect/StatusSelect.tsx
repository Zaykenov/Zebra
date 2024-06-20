import React, { FC, useEffect, useMemo, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllAccounts } from "@api/accounts";
import { getAllSuppliers } from "@api/suppliers";
import { getAllProductCategories } from "@api/product-categories";

export interface StatusSelectProps {
  onChange: (value: string) => void;
  className?: string;
}

const StatusSelect: FC<StatusSelectProps> = ({ onChange, className = "" }) => {
  const options = useMemo(
    () => [
      {
        name: "Статус",
        value: "",
      },
      {
        name: "Открыт",
        value: "opened",
      },
      {
        name: "Наличкой",
        value: "closed",
      },
    ],
    []
  );

  return (
    <Select
      options={options}
      onChange={(e) => onChange((e.target as HTMLSelectElement).value)}
      className={clsx(["w-40", className])}
    />
  );
};

export default StatusSelect;
