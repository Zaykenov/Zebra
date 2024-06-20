import React, { FC, useEffect, useState } from "react";
import { Select } from "@shared/ui/Select";
import { SelectOption } from "@shared/ui/Select/Select";
import clsx from "clsx";
import { getAllAccounts } from "@api/accounts";
import { getAllSuppliers } from "@api/suppliers";
import { getAllProductCategories } from "@api/product-categories";

export interface ProductCategorySelectProps {
  onChange: (value: number) => void;
  className?: string;
}

const ProductCategorySelect: FC<ProductCategorySelectProps> = ({
  onChange,
  className = "",
}) => {
  const [options, setOptions] = useState<SelectOption[]>([]);

  useEffect(() => {
    getAllProductCategories().then((res) => {
      setOptions([
        {
          name: "Все категории",
          value: 0,
        },
        ...res.data.map((category: any) => ({
          name: category.name,
          value: category.id,
        })),
      ]);
    });
  }, []);

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

export default ProductCategorySelect;
