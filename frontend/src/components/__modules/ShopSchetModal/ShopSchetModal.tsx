import React, { Dispatch, FC, SetStateAction } from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import AccountForm from "../AccountForm";

export interface ShopSchetModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  type: "cash" | "card";
  onCreate: (id: number) => void;
}

const ShopSchetModal: FC<ShopSchetModalProps> = ({
  isOpen,
  setIsOpen,
  type,
  onCreate,
}) => {
  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title={type === "cash" ? "Новый наличный счет" : "Новый безналичный счет"}
    >
      <div className="flex flex-col items-start">
        <AccountForm width="w-full" type={type} onCreate={onCreate} />
      </div>
    </ModalLayout>
  );
};

export default ShopSchetModal;
