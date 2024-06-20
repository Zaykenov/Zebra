import React, {
  Dispatch,
  FC,
  Fragment,
  ReactNode,
  SetStateAction,
} from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XMarkIcon } from "@heroicons/react/24/outline";
import clsx from "clsx";

export interface ModalLayoutProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  onClose?: () => void;
  title: string;
  children: ReactNode;
  disableClose?: boolean;
  maxWidth?: string;
  footer?: ReactNode;
  fullScreen?: boolean;
}

const ModalLayout: FC<ModalLayoutProps> = ({
  isOpen,
  setIsOpen,
  onClose,
  title,
  children,
  disableClose = false,
  maxWidth = "max-w-xl",
  footer,
  fullScreen = false,
}) => {
  function closeModal() {
    if (disableClose) return;
    setIsOpen(false);
    onClose && onClose();
  }

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={closeModal}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div
            className={clsx([
              "flex items-center justify-center p-4 text-center",
              fullScreen ? "h-screen" : "min-h-full",
            ])}
          >
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel
                className={clsx([
                  "w-full flex flex-col transform h-full overflow-hidden rounded bg-white text-left align-middle shadow-xl transition-all",
                  maxWidth,
                ])}
              >
                <Dialog.Title
                  as="div"
                  className="w-full flex items-center justify-between px-6 py-4 border-b"
                >
                  <span className="text-lg font-medium leading-6 text-gray-900">
                    {title}
                  </span>
                  <button
                    onClick={closeModal}
                    className="p-1 rounded-md hover:bg-gray-200"
                  >
                    <XMarkIcon className="w-10 h-10" />
                  </button>
                </Dialog.Title>
                <div className="px-6 py-4 flex flex-col h-full overflow-y-auto overflow-x-hidden">
                  {children}
                </div>
                {footer && <div className="">{footer}</div>}
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};

export default ModalLayout;
