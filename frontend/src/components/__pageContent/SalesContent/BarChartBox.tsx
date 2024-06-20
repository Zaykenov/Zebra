import React, { FC, useState } from "react";
import clsx from "clsx";
import {
  Bar,
  BarChart,
  CartesianGrid,
  LabelList,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { formatNumber } from "@utils/formatNumber";
import { BarChartBoxProps } from "./types";
import CustomTooltip from "./CustomTooltip";
import CustomLabel from "./CustomLabel";

const BarChartBox: FC<BarChartBoxProps> = ({
  title,
  vertical = false,
  switchable = false,
  data,
  className = "",
  tooltipPayload,
}) => {
  const [value1Active, setValue1Active] = useState<boolean>(true);

  return (
    <div
      className={clsx([
        "px-8 py-4 flex flex-col rounded border border-gray-300",
        className,
      ])}
    >
      <div className="flex items-end justify-between">
        <span className="text-lg font-medium">{title}</span>
        {switchable && (
          <div className="flex items-center space-x-3">
            <button
              onClick={() => setValue1Active(true)}
              type="button"
              className={clsx([
                "text-xs",
                !value1Active
                  ? "text-blue-500 cursor-pointer"
                  : "cursor-default font-medium",
              ])}
            >
              Оборот
            </button>
            <button
              onClick={() => setValue1Active(false)}
              type="button"
              className={clsx([
                "text-xs",
                value1Active
                  ? "text-blue-500 cursor-pointer"
                  : "cursor-default font-medium",
              ])}
            >
              Чеки
            </button>
          </div>
        )}
      </div>
      <div className="flex flex-col">
        <ResponsiveContainer width="100%" height={150}>
          <BarChart
            className="w-full"
            data={
              value1Active
                ? data
                : data.map((item) => ({ ...item, value: item.value2 }))
            }
            margin={{ top: 20, left: 10 }}
            layout={vertical ? "vertical" : "horizontal"}
            maxBarSize={vertical ? 20 : 100}
          >
            {vertical ? (
              <>
                <XAxis type="number" hide />
                <YAxis
                  dataKey="name"
                  axisLine={false}
                  tickLine={false}
                  minTickGap={3}
                  type="category"
                  tick={{ fontSize: "12px" }}
                />
              </>
            ) : (
              <>
                <YAxis
                  type="number"
                  dataKey="value"
                  axisLine={false}
                  tickLine={false}
                  tick={{ fontSize: "12px" }}
                />
                <XAxis
                  dataKey="name"
                  axisLine={false}
                  tickLine={false}
                  type="category"
                  tick={{ fontSize: "12px", fontWeight: 600 }}
                />
              </>
            )}
            {!vertical && (
              <CartesianGrid strokeDasharray="4 1 1 1 1" vertical={false} />
            )}
            {!vertical && (
              <Tooltip
                content={({ active, payload, label }) => (
                  <CustomTooltip
                    active={active}
                    payload={
                      payload &&
                      payload[0] &&
                      tooltipPayload && {
                        name: tooltipPayload.name,
                        value: formatNumber(
                          payload[0].payload.value,
                          tooltipPayload.isCurrency,
                          tooltipPayload.isCurrency
                        ),
                      }
                    }
                    label={label}
                  />
                )}
              />
            )}
            <Bar
              dataKey="value"
              fill="#3EB2B2"
              background={{ fill: "#f3f4f6" }}
            >
              {vertical && (
                <LabelList
                  dataKey="value"
                  position="top"
                  className="text-sm text-primary font-inter font-medium"
                  content={(props) => (
                    <CustomLabel {...props} isCurrency={value1Active} />
                  )}
                />
              )}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default BarChartBox;
