import React, { FC, ForwardedRef, forwardRef, HTMLProps } from "react";
import Select from "./Select";

export interface LabeledSelectProps extends HTMLProps<HTMLSelectElement> {
  label: string;
  options: {
    name: string;
    value: string | number;
    data?: any;
  }[];
  fieldClass?: string;
  labelClass?: string;
}

const LabeledSelect = forwardRef(
  (
    {
      label,
      options,
      fieldClass = "",
      labelClass = "",
      className = "",
      ...props
    }: LabeledSelectProps,
    ref: ForwardedRef<HTMLSelectElement>
  ) => {
    return (
      <div className={`w-full flex items-center ${fieldClass}`}>
        <label className={`w-40 mr-4 ${labelClass}`}>{label}</label>
        <Select
          options={options}
          className={`flex-1 ${className}`}
          {...props}
          ref={ref}
        />
      </div>
    );
  }
);

LabeledSelect.displayName = "LabeledSelect";

export default LabeledSelect;
