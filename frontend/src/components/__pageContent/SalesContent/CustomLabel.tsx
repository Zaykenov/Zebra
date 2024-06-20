import React, { FC } from "react";
import { formatNumber } from "@utils/formatNumber";

const CustomLabel: FC<{
  x?: number | string;
  y?: number | string;
  width?: number | string;
  height?: number | string;
  value?: string | number;
  isCurrency?: boolean;
}> = ({ x, y, width, value, isCurrency }) => {
  return (
    <g>
      <text
        x={(x as number) + (width as number) - 50}
        y={(y as number) - 15}
        textAnchor="middle"
        dominantBaseline="middle"
        fill="#000"
        fontSize="small"
        fontWeight="bold"
      >
        {formatNumber(value as number, isCurrency, isCurrency)}
      </text>
    </g>
  );
};

export default CustomLabel;
