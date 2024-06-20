import React, { FC } from "react";
import ButtonGroup, { ButtonGroupProps } from "./ButtonGroup";
import Select from "../Select/Select";

export interface LabeledButtonGroupProps extends ButtonGroupProps {
  label: string;
  fieldClass?: string;
  labelClass?: string;
}

const LabeledButtonGroup: FC<LabeledButtonGroupProps> = ({
  buttons,
  onBtnClick,
  value,
  label,
  fieldClass = "",
  labelClass = "",
}) => {
  return (
    <div className={`w-full flex items-center ${fieldClass}`}>
      <label className={`w-40 mr-4 ${labelClass}`}>{label}</label>
      <ButtonGroup buttons={buttons} onBtnClick={onBtnClick} value={value} />
    </div>
  );
};

export default LabeledButtonGroup;
