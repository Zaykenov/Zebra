import React, { ForwardedRef, forwardRef, HTMLProps } from "react";

export type SelectOption = {
  name: string;
  value: string | number;
  data?: any;
};

export interface SelectProps extends HTMLProps<HTMLSelectElement> {
  options: SelectOption[];
}

const Select = forwardRef(
  (
    { options, className, ...props }: SelectProps,
    ref: ForwardedRef<HTMLSelectElement>
  ) => {
    return (
      <select
        className={`w-full rounded text-sm text-gray-800 py-2 px-3 transition-[border] duration-300 border border-gray-300 focus:outline-none focus:border-indigo-500 cursor-pointer ${className}`}
        {...props}
        ref={ref}
      >
        {options.map((option, idx) => (
          <option key={`${option.name}_${idx}`} value={option.value}>
            {option.name}
          </option>
        ))}
      </select>
    );
  }
);

Select.displayName = "Select";

export default Select;
