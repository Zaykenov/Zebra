import React, {
  Dispatch,
  FC,
  SetStateAction,
  useState,
} from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import { Input } from "@shared/ui/Input";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage/AlertMessage";
import { UserMobileData, getUserByQR } from "@api/mobile";

export interface CommentCheckModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  setDiscount: (data: number) => void;
  setData: (data: UserMobileData) => void;
}

const CodeModal: FC<CommentCheckModalProps> = ({
  isOpen,
  setIsOpen,
  setDiscount,
  setData
}) => {
  const {alertMessage, showAlertMessage, hideAlertMessage} = useAlertMessage();
  const [code, setCode] = useState<string>("");

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Ввести код пользователя"
    >
      <div className="flex flex-col">
        <div className="w-full flex space-x-3 mb-2">
          <Input
            className=""
            type="text"
            value={code}
            onInput={(e) => {
              setCode((e.target as HTMLInputElement).value);
            }}
            name="text"
          />
        </div>
        <button
          disabled={code.length !== 4}
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            try{
              getUserByQR(code).then((res)=>{
                console.log(res)
                setData(res)
                setDiscount(res.discount*100)
                setIsOpen(false)
              })
            }catch(e){
              alert("Код не найден или больше не действителен")
              // showAlertMessage("Код не найден или больше не действителен", AlertMessageType.ERROR)
            }
          }}
          className="mt-3 pt-1 pb-1.5 text-sm flex items-center justify-center bg-primary disabled:bg-gray-300 disabled:cursor-not-allowed hover:opacity-80 text-white rounded"
        >
          Сохранить
        </button>
      </div>
      {alertMessage && (
        <AlertMessage message={alertMessage.message} type={alertMessage.type} onClose={hideAlertMessage}/>
      )}
    </ModalLayout>
  );
};

export default CodeModal;
