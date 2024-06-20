import React, { FC } from "react";
import { Popover } from "@headlessui/react";
import { formatNumber } from "@utils/formatNumber";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import { DetailsPopoverProps } from "./types";

const DetailsPopover: FC<DetailsPopoverProps> = ({ value, row, details }) => {
  return (
    <Popover className="relative">
      <Popover.Button
        onClick={() => {}}
        className="flex items-center space-x-3"
      >
        <span>
          {/* @ts-ignore */}
          {formatNumber(value, false, false)} {row.original.measure}
        </span>
        <div className="p-0.5 border border-indigo-300 rounded">
          <ChevronDownIcon className="w-2 h-2 text-indigo-300" />
        </div>
      </Popover.Button>

      <Popover.Panel className="absolute z-10 mt-5 min-w-[400px] -ml-40 p-4 border border-gray-400 shadow-lg rounded bg-white">
        <div className="grid grid-cols-3 gap-2 text-xs">
          {details.postavka !== undefined && (
            <>
              <div>Поставка</div>
              <div>
                {formatNumber(details.postavka, false, false)} {details.measure}
              </div>
              <div>{formatNumber(details.postavka_cost || 0, true, false)}</div>
            </>
          )}
          {details.sales !== undefined && (
            <>
              <div>Продажи</div>
              <div>
                {formatNumber(details.sales, false, false)} {details.measure}
              </div>
              <div>---</div>
            </>
          )}
          <>
            <div>Инвентаризация</div>
            {details.inventarization !== undefined ? (
              <div>
                {formatNumber(details.inventarization, false, false)}{" "}
                {details.measure}
              </div>
            ) : (
              <div>---</div>
            )}
            <div>---</div>
          </>
          <>
            <div>Перемещение</div>
            {details.transfer !== undefined ? (
              <div>
                {formatNumber(details.transfer, false, false)} {details.measure}
              </div>
            ) : (
              <div>---</div>
            )}
            <div>---</div>
          </>
        </div>
      </Popover.Panel>
    </Popover>
  );
};

export default DetailsPopover;
