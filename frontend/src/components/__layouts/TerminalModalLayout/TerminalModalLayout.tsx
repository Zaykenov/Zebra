import React, {
    Dispatch,
    FC,
    Fragment,
    ReactNode,
    SetStateAction,
  } from "react";
  import { Dialog, Transition } from "@headlessui/react";
  import { XMarkIcon } from "@heroicons/react/24/outline";
  
  export interface TerminalModalLayoutProps {
    isOpen: boolean;
    setIsOpen: Dispatch<SetStateAction<boolean>>;
    onClose?: () => void;
    title: string;
    children: ReactNode;
    disableClose?: boolean;
  }
  
  const TerminalModalLayout: FC<TerminalModalLayoutProps> = ({
    isOpen,
    setIsOpen,
    onClose,
    title,
    children,
    disableClose = false,
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
            <div className="flex h-full items-center justify-center p-4 text-center">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel className="w-3/5 max-w-xl transform overflow-hidden rounded bg-white text-left align-middle shadow-xl transition-all"
                style={{ width: "95%", height: "100%", maxWidth: "none", maxHeight: "none" }}>
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
                      <XMarkIcon className="w-12 h-12" />
                    </button>
                  </Dialog.Title>
                  <div className="px-6 py-4">{children}</div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>
    );
  };
  
  export default TerminalModalLayout;
  