import React, { FC, useEffect, useMemo, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllAccounts } from "@api/accounts";
import { getAllSuppliers } from "@api/suppliers";
import { getAllProductCategories } from "@api/product-categories";

export interface PaymentSelectProps {
  onChange: (value: string) => void;
  className?: string;
}

const PaymentSelect: FC<PaymentSelectProps> = ({
  onChange,
  className = "",
}) => {
  const options = useMemo(
    () => [
      {
        name: "Способ оплаты",
        value: "",
      },
      {
        name: "Картой",
        value: "Картой",
      },
      {
        name: "Наличкой",
        value: "Наличными",
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

export default PaymentSelect;
