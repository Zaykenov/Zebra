import { FC } from "react";

interface TerminalButtonProps {
    setIsModalOpen: (value: React.SetStateAction<boolean>) => void,
    buttonContent: string
}
const TerminalModalButton: FC<TerminalButtonProps> = ({
    setIsModalOpen,
    buttonContent
}) => {
 return (
    <button
        className="col-span-2 py-4 px-3 text-sm tablet:text-base border border-gray-300 rounded-md bg-neutral-100 hover:bg-neutral-200 shadow"
        onClick={() => {
            setIsModalOpen(true);
        }}
    >
        {buttonContent}
    </button>
 )
}

export default TerminalModalButton