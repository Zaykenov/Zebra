import React, { FC } from "react";

const CustomTooltip: FC<{
  active?: boolean;
  payload?: { name: string; value: string } | null;
  label?: string;
}> = ({ active, payload }) => {
  if (!active || !payload) return <></>;
  return (
    <div className="p-2 flex flex-col space-y-1 border border-gray-400 bg-white rounded-md">
      <span className="capitalize text-xs font-light">{payload.name}</span>
      <span>{payload.value}</span>
    </div>
  );
};

export default CustomTooltip;
