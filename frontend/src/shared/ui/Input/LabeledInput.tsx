import React, {
  FC,
  ForwardedRef,
  forwardRef,
  HTMLProps,
  useEffect,
} from "react";
import Input from "./Input";
import { EyeIcon, EyeSlashIcon } from "@heroicons/react/24/outline";

export interface LabeledInputProps extends HTMLProps<HTMLInputElement> {
  label: string;
  fieldClass?: string;
  labelClass?: string;
  showPassword?: boolean;
  handlePassword?: () => void;
}

const LabeledInput = forwardRef(
  (
    {
      label,
      fieldClass = "",
      labelClass = "",
      className = "",
      showPassword,
      handlePassword,
      ...props
    }: LabeledInputProps,
    ref: ForwardedRef<HTMLInputElement>
  ) => {
    return (
      <div className={`w-full flex items-center ${fieldClass}`}>
        <label className={`w-40 mr-4 ${labelClass}`}>{label}</label>
        <div className="flex-1 relative">
          <Input className={`w-full ${className}`} {...props} ref={ref} />
          {handlePassword && (
            <div className="absolute z-10 inset-y-0 right-0 flex items-center justify-center px-2">
              <button
                type="button"
                className="p-1 rounded-md"
                onClick={handlePassword}
              >
                {showPassword ? (
                  <EyeSlashIcon />
                ) : (
                  <EyeIcon className="h-4 w-4" />
                )}
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }
);

LabeledInput.displayName = "LabeledInput";

export default LabeledInput;
