import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import { UserShiftData } from "@modules/TerminalHeader/TerminalHeader";
import { checkShift } from "@api/shifts";

export interface InfoModalProps {
  isOpen: boolean;
  shiftData: any
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  onClose?: () => void;
}

const InfoModal: FC<InfoModalProps> = ({ isOpen, shiftData, setIsOpen, onClose }) => {
  const [data, setData] = useState<UserShiftData | null>(null);

  useEffect(() => {
    checkShift().then((res) => {
      setData({ ...res.data, shopName: res.data.shop_name });
    });
  }, []);

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onClose={onClose}
      title="Информация о заведении"
    >
      <div className="flex flex-col items-center">
        <h2 className="text-xl font-medium text-red-600 mb-4">Внимание!</h2>
        <p className="text-lg text-center mb-4">
          Вы авторизовались в следующем заведении:
        </p>
        <span className="text-2xl font-bold mb-8">{shiftData && shiftData.shop_name}</span>
        <button
          type="button"
          onClick={() => {
            setIsOpen(false);
            onClose && onClose();
          }}
          className="w-full flex flex-col items-center justify-center text-center py-2 bg-primary text-white rounded-md"
        >
          Понятно
        </button>
      </div>
    </ModalLayout>
  );
};

export default InfoModal;
