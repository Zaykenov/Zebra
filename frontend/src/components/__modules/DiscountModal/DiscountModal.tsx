import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import clsx from "clsx";

export interface DiscountModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  setData: (data: number) => void;
  onClose?: () => void;
  shiftData: any;
  data: number;
}

const DiscountModal: FC<DiscountModalProps> = ({
  isOpen,
  setIsOpen,
  setData,
  onClose,
  shiftData,
  data,
}) => {
  const [discounts, setDiscounts] = useState<number[]>([5, 10, 15, 20])

  useEffect(()=>{
    if(shiftData.shop_id === 11)
    setDiscounts(discounts.filter((discount)=>discount > 10))
  }, [shiftData])

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onClose={onClose}
      title="Добавить скидку"
    >
      <form className="flex flex-col">
        <ul className="w-full flex justify-center flex-wrap gap-2 mb-4">
          {discounts.map((discount) => (
            <li
              key={discount}
              className={clsx([
                "pt-0.5 pb-1 pl-3 pr-2 flex items-center space-x-2 cursor-pointer rounded-3xl text-white",
                data === discount
                  ? "bg-primary ring ring-teal-600"
                  : "bg-primary/75",
              ])}
              onClick={() => {
                data === discount ? setData(0) : setData(discount);
              }}
            >
              <span className="text-sm font-medium">{discount}%</span>
            </li>
          ))}
        </ul>
        <button
          type="button"
          onClick={() => {
            onClose && onClose();
            setIsOpen(false);
          }}
          className="mt-3 pt-1 pb-1.5 text-sm flex items-center justify-center bg-primary hover:opacity-80 text-white rounded"
        >
          Сохранить
        </button>
      </form>
    </ModalLayout>
  );
};

export default DiscountModal;
