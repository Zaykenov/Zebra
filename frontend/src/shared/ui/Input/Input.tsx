import React, {
  FC,
  ForwardedRef,
  forwardRef,
  HTMLProps,
  useEffect,
} from "react";

export interface InputProps extends HTMLProps<HTMLInputElement> {}

const Input = forwardRef(
  (
    { className, ...props }: InputProps,
    ref: ForwardedRef<HTMLInputElement>
  ) => {
    return (
      <input
        onClick={(e) => (e.target as HTMLInputElement).select()}
        ref={ref}
        className={`w-full rounded text-sm transition-[border] duration-300 text-gray-800 py-2 px-3 border border-gray-300 focus:outline-none focus:border-indigo-500 ${className}`}
        {...props}
      />
    );
  }
);

Input.displayName = "Input";

export default Input;
