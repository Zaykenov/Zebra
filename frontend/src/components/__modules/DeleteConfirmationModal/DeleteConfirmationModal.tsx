import React, { Dispatch, FC, ReactNode, SetStateAction, useEffect, useState } from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import { Row } from "react-table";
import { getDeleteResourceDetails } from "@api/index";

export interface DeleteConfirmationModalProps {
  isOpen: boolean;
  isDetailed?: boolean;
  row?: Row<any>;
  path?: string;
  deleteConfirmationText?: string;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  onDelete: (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => void;
}

const DeleteConfirmationModal: FC<DeleteConfirmationModalProps> = ({
  isOpen,
  row,
  path,
  isDetailed,
  setIsOpen,
  onDelete,
  deleteConfirmationText
}) => {

  const [deleteResourceName, setDeleteResourceName] = useState<string>();
  const [link, setLink] = useState<string>();
  const [message, setMessage] = useState<string>('');

  useEffect(()=>{
    if(isOpen && isDetailed && path){
      getDeleteResourceDetails(row?.original.id, path).then((res)=>{
        if (res) {
          setDeleteResourceName(res.name)
          setLink(res.link)
          setMessage(res.message)
        }
      })
    }
  }, [isOpen, isDetailed])

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Подтверждение удаления"
    >
      <div className="relative p-4 text-center bg-white sm:p-5">
        <div className="mb-4">
          {isDetailed ? (
            <div>
              <p className="mb-4 text-gray-500 dark:text-gray-300">
                {message}
              </p>
              <p className="mb-4 text-blue-500 dark:text-red-300">
                <a href={link}>{deleteResourceName}</a>
              </p>
            </div>
          ) : (
            <>
              <svg className="text-gray-400 dark:text-gray-500 w-11 h-11 mb-3.5 mx-auto">
                <path
                  fillRule="evenodd"
                  d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                ></path>
              </svg>
              <p className="mb-4 text-gray-500 dark:text-gray-300">
                {deleteConfirmationText ? deleteConfirmationText: "Вы уверены, что хотите удалить это?"}
              </p>
            </>
          )}
        </div>
        <div className="flex justify-center items-center space-x-4">
          <button
            type="button"
            className="py-2 px-3 text-sm font-medium text-gray-500 bg-white rounded-lg border border-gray-200 hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-primary-300 hover:text-gray-900 focus:z-10 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-500 dark:hover:text-white dark:hover:bg-gray-600 dark:focus:ring-gray-600"
            onClick={() => setIsOpen(false)}
          >
            Нет, отмена
          </button>
          <button
            type="submit"
            className="py-2 px-3 text-sm font-medium text-center text-white bg-red-600 rounded-lg hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-red-300 dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900"
            onClick={(e) => {
              onDelete(e);
            }}
          >
            Да, продолжить
          </button>
        </div>
      </div>
    </ModalLayout>
  );
};

export default DeleteConfirmationModal;
