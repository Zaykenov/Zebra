import React, { FC, ForwardedRef, forwardRef, HTMLProps } from "react";

export interface CheckboxProps extends HTMLProps<HTMLInputElement> {
  name: string;
}

const Checkbox = forwardRef(
  (
    { label, className, ...props }: CheckboxProps,
    ref: ForwardedRef<HTMLInputElement>
  ) => {
    return (
      <div className="relative flex items-start">
        <div className="flex h-5 items-center">
          <input
            ref={ref}
            type="checkbox"
            className={`h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 ${className}`}
            {...props}
          />
        </div>
        <div className="ml-3 text-sm">
          <label htmlFor="comments" className="font-medium text-gray-700">
            {label}
          </label>
        </div>
      </div>
    );
  }
);

Checkbox.displayName = "Checkbox";

export default Checkbox;
