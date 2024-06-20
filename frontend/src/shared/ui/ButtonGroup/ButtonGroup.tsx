import React, { FC } from "react";
import clsx from "clsx";

export interface ButtonGroupProps {
  buttons: {
    label: string;
    value: string;
  }[];
  value: string;
  onBtnClick: (value: string) => void;
}

const ButtonGroup: FC<ButtonGroupProps> = ({ buttons, onBtnClick, value }) => {
  return (
    <span className="isolate inline-flex rounded-md shadow-sm">
      {buttons.map((btn) => (
        <button
          type="button"
          onClick={() => {
            onBtnClick(btn.value);
          }}
          className={clsx([
            "relative inline-flex items-center border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 focus:z-10",
            btn.value === value && "bg-indigo-100 border-indigo-500",
          ])}
        >
          {btn.label}
        </button>
      ))}
    </span>
  );
};

export default ButtonGroup;
