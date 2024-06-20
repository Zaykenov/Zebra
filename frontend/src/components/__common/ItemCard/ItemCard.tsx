import React, { FC, ReactNode, useEffect, useState } from "react";
import clsx from "clsx";
import { EllipsisHorizontalIcon } from "@heroicons/react/24/outline";
import { stringToColor } from "@utils/stringToColor";

export interface ItemCardProps {
  name: string;
  cover?: string;
  onSelect: (e?: any) => void;
  price?: number;
  hasModal?: boolean;
  className?: string;
  height?: string;
  quantity?: number;
  onRemove?: (e?: any) => void;
}

const ItemCard: FC<ItemCardProps> = ({
  name,
  cover,
  price,
  hasModal = false,
  onSelect,
  className,
  height = "h-[130px]",
  quantity,
  onRemove,
}) => {
  return (
    <button
      type="button"
      className={clsx(["w-full overflow-visible", className])}
      onClick={(e) => {
        e.stopPropagation();
        e.preventDefault();
        onSelect(e);
      }}
    >
      <div
        className={clsx([`relative shadow-md rounded flex flex-col`, height])}
      >
        {!!quantity && (
          <div className="absolute inset-0 bg-black/50 flex flex-col items-center justify-center text-3xl text-white font-bold">
            <button
              type="button"
              className="absolute z-10 rounded-full leading-[90%] w-6 h-6 bg-red-500 top-0 right-0 -mt-3 -mr-3 pb-1 font-bold text-xl"
              onClick={onRemove}
            >
              -
            </button>
            <span>{quantity}</span>
          </div>
        )}
        <div
          className={clsx([
            "grow flex items-center justify-center uppercase font-medium text-white text-3xl",
          ])}
          style={
            cover && cover !== "image.png"
              ? {
                  background: `url('https://zebra-crm.kz:8029/itemImage/${cover}') no-repeat center`,
                }
              : {
                  backgroundColor: stringToColor(name),
                }
          }
        >
          {(!cover || cover === "image.png") &&
            name
              .split(" ")
              .filter((_, i) => i < 2)
              .map((str) => str.charAt(0))
              .join("")}
        </div>
        <div className="py-1 px-2 bg-white flex justify-between items-center">
          <span className="text-sm text-gray-600 font-semibold text-left">
            {name}
          </span>
          {price && (
            <span className="text-xs text-gray-400 font-bold whitespace-nowrap">
              {price} ₸
            </span>
          )}
        </div>
      </div>
    </button>
  );
};

export default ItemCard;
