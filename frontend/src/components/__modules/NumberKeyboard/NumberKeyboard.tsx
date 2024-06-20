import React, { FC } from "react";

export interface NumberKeyboardProps {
  onKeyClick: (val: string | number) => void;
  onDelete: () => void;
}

const NumberKeyboard: FC<NumberKeyboardProps> = ({ onKeyClick, onDelete }) => {
  return (
    <div
      className="w-full h-full bg-white-100 flex justify-center content-center"
      style={{ overflow: "hidden" }}
    >
      <div style={{ width: "80%", height: "80%" }} className="mt-5">
        <div className="grid grid-cols-6 gap-4 h-full font-mono text-black text-3xl text-center font-bold leading-10 text-center">
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(7)}
          >
            7
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(8)}
          >
            8
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(9)}
          >
            9
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(4)}
          >
            4
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(5)}
          >
            5
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(6)}
          >
            6
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(1)}
          >
            1
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(2)}
          >
            2
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(3)}
          >
            3
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(".")}
          >
            .
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={() => onKeyClick(0)}
          >
            0
          </button>
          <button
            className="p-4 rounded-lg shadow-md bg-white col-span-2 flex justify-center items-center"
            onClick={onDelete}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth="1.5"
              stroke="currentColor"
              className="w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 9.75L14.25 12m0 0l2.25 2.25M14.25 12l2.25-2.25M14.25 12L12 14.25m-2.58 4.92l-6.375-6.375a1.125 1.125 0 010-1.59L9.42 4.83c.211-.211.498-.33.796-.33H19.5a2.25 2.25 0 012.25 2.25v10.5a2.25 2.25 0 01-2.25 2.25h-9.284c-.298 0-.585-.119-.796-.33z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
};

export default NumberKeyboard;
