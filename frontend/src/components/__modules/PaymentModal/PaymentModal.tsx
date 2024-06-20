import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import clsx from "clsx";
import { Oval } from "react-loader-spinner";
import { useRouter } from "next/router";
import axios from "axios";
import { PaymentMethod, printCheck } from "@api/check";
import TerminalModalLayout from "@layouts/TerminalModalLayout/TerminalModalLayout";
import TerminalInput from "../TerminalInput";
import AlertMessage from "@common/AlertMessage";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import { IdempotencyKey } from "@api/idempotencyKey";
import { useAppSelector } from "@hooks/useAppSelector";
import { processPrintObject } from "@utils/processPrintObject";

// export type OnSubmitResponse = {
//   registeredCheck: Promise<any>;
//   printData: any
// }

export interface PaymentModalProps {
  isOpen: boolean;
  shiftData: any;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  paymentDone: boolean;
  setPaymentDone: Dispatch<SetStateAction<boolean>>;
  total: number;
  onClose: () => void;
  onSubmit: (
    payment: PaymentMethod,
    card: number,
    cash: number,
    pager: number,
  ) => Promise<any>;
  discount: number;
  clearStates: boolean;
}

const pagers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15];

const PaymentModal: FC<PaymentModalProps> = ({
  isOpen,
  setIsOpen,
  shiftData,
  total,
  onClose,
  onSubmit,
  clearStates,
  paymentDone,
  setPaymentDone,
}) => {
  const idempotencyState = useAppSelector((state) => state.idempotency);

  const router = useRouter();
  const [timer] = useState<NodeJS.Timer>();
  const [payment, setPayment] = useState<PaymentMethod>(PaymentMethod.CARD);
  const [cash, setCash] = useState<number>(0);
  const [card, setCard] = useState<number>(total);
  const [link, setLink] = useState<string | null>(null);
  const [printData, setPrintData] = useState<string>("");
  const [currentInput, setCurrentInput] = useState<string>("Card");
  const [loading, setLoading] = useState<boolean>(false);
  const [btnLoading, setBtnLoading] = useState<boolean>(false);
  const [disableClose, setDisableClose] = useState<boolean>(false);
  const [surplus, setSurplus] = useState<number>(0);
  const [isChangeAvailable, setIsChangeAvailable] = useState<boolean>(false);
  const [stopSwitching, setStopSwitching] = useState<boolean>(false);

  const [hasPagers, setHasPagers] = useState<boolean | undefined>();

  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const [selectedPager, setSelectedPager] = useState<number | null>(null);

  const handleSetTotalPriceClick = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
  ) => {
    currentInput === "Cash" ? setCash(total) : setCard(total);
  };

  const handleAddClick = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
  ) => {
    const buttonValue = parseFloat(e.currentTarget.textContent ?? "0");

    if (currentInput === "Cash") {
      if (cash === total) {
        setCash(buttonValue);
      } else {
        const newCash = parseFloat(`${cash}${buttonValue}`);
        setCash(newCash);
      }
    } else {
      if (card === total) {
        setCard(buttonValue);
      } else {
        const newCard = parseFloat(`${card}${buttonValue}`);
        setCard(newCard);
      }
    }
  };

  const handleDeleteClick = () => {
    if (currentInput === "Cash") {
      let currentValue = cash.toString();
      if (currentValue.length > 0) {
        const newValue = parseFloat(currentValue.slice(0, -1));
        setCash(!isNaN(newValue) ? newValue : 0);
      }
    } else {
      let currentValue = card.toString();
      if (currentValue.length > 0) {
        const newValue = parseFloat(currentValue.slice(0, -1));
        setCard(!isNaN(newValue) ? newValue : 0);
      }
    }
  };

  const handlePaymentSubmit = async () => {
    try {
      setBtnLoading(true);
      setLoading(true);
      setDisableClose(true);
      let keys: IdempotencyKey[] = JSON.parse(
        localStorage.getItem("zebra.idempotencyKeys") || "[]",
      );
      keys.push({
        idempotency_key: idempotencyState.idempotencyKey,
        time: new Date().toISOString(),
      });
      localStorage.setItem("zebra.idempotencyKeys", JSON.stringify(keys));
      // prod without mobile support
      const res = await onSubmit(payment, card, cash, selectedPager || 1);
      const tisCheckUrl = res.data.link;
      const printObj = processPrintObject(res.data, true);
      const printData = JSON.stringify(printObj);
      setPrintData(printData);
      if (localStorage.getItem("zebra.automaticPrint") == "true")
        printCheck(printData);
      setLink(tisCheckUrl);

      // prod with mobile
      // const {registeredCheck, printData} = onSubmit(payment, card, cash, selectedPager || 1);
      // setPrintData(printData);
      // registeredCheck.then((res)=>{
      //   const tisCheckUrl = shiftData.shop_name == "Ландмарк" ? res.data.link : ""
      //   setLink(tisCheckUrl)
      // })
      // if (localStorage.getItem("zebra.automaticPrint") == "true") printCheck(printData);
    } catch (e: any) {
      if (axios.isAxiosError(e)) {
        if (e.response?.status === 666) {
          showAlertMessage(`НАЧНИТЕ СМЕНУ!`, AlertMessageType.WARNING);
          setTimeout(() => {
            return router.push("/terminal/shift");
          }, 2000);
        } else {
          showAlertMessage(
            `Оплата не прошла. Повторите еще раз. Детали: ${e.message}`,
            AlertMessageType.ERROR,
          );
        }
      } else console.log(e.message);
    } finally {
      setLoading(false);
      setDisableClose(false);
      setBtnLoading(false);
      setPaymentDone(true);
    }
  };

  useEffect(() => {
    if (cash + card > total && cash > card) {
      setSurplus(total - card - cash);
    }
  }, [cash, card, total]);

  useEffect(() => {
    setCard(total);
  }, [total]);

  useEffect(() => {
    if (cash < 0) {
      setCash(0);
    }

    if (card < 0) {
      setCard(0);
    }

    if (cash === 0) {
      setPayment(PaymentMethod.CARD);
    } else if (card === 0) {
      setPayment(PaymentMethod.CASH);
    } else {
      setPayment(PaymentMethod.MIXED);
    }
    let legitSurplus = cash + card - total;
    legitSurplus < 0 ? setSurplus(0) : setSurplus(legitSurplus);
    setIsChangeAvailable(total - cash > -20000);
  }, [cash, card, currentInput]);

  useEffect(() => {
    setStopSwitching(true);
  }, [handleAddClick]);

  useEffect(() => {
    const hasPagersLocalStorage =
      localStorage.getItem("zebra.hasPagers") === "true";
    setHasPagers(hasPagersLocalStorage);
    if (!hasPagersLocalStorage) {
      setSelectedPager(1);
    }
  }, []);

  useEffect(() => {
    setPaymentDone(false);
    setPrintData("");
    setLink(null);
  }, [clearStates]);

  return (
    <TerminalModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onClose={() => {
        setCash(0);
        setCard(0);
        if (paymentDone) {
          localStorage.removeItem("zebra.activeCheck");
          onClose();
        }
      }}
      title={link ? "Чек" : `К оплате: ${total} ₸`}
      disableClose={disableClose}
    >
      <div className="flex w-full h-screen">
        <TerminalInput
          total={total}
          handleAddClick={handleAddClick}
          handleDeleteClick={handleDeleteClick}
          handleSetTotalPriceClick={handleSetTotalPriceClick}
        />
        {alertMessage && (
          <AlertMessage
            message={alertMessage.message}
            type={alertMessage.type}
            onClose={hideAlertMessage}
          />
        )}
        <div
          className="w-3/5 h-full bg-neutral-100 flex justify-center content-center"
          style={{ overflow: "hidden" }}
        >
          <div
            style={{ width: "80%", height: "80%" }}
            className="bg-white-500 mt-10"
          >
            <div className="relative flex w-full flex-wrap items-center mb-3">
              <button
                className="absolute inset-0 bg-transparent z-20"
                onClick={() => {
                  setCurrentInput("Card");
                  if (cash == total) {
                    setCash(0);
                    setCard(total);
                  }
                }}
              />
              <span className="z-10 flex flex-col items-center justify-center h-full leading-snug font-normal absolute text-center text-slate-300 absolute bg-transparent rounded text-base items-center justify-center w-8 pl-3 py-3">
                <svg
                  style={{ color: "rgb(78, 170, 29)" }}
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  className="h-6 w-6"
                  viewBox="0 0 16 16"
                >
                  <path
                    d="M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V4zm2-1a1 1 0 0 0-1 1v1h14V4a1 1 0 0 0-1-1H2zm13 4H1v5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V7z"
                    fill="#4eaa1d"
                  ></path>{" "}
                  <path
                    d="M2 10a1 1 0 0 1 1-1h1a1 1 0 0 1 1 1v1a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-1z"
                    fill="#4eaa1d"
                  ></path>
                </svg>
              </span>
              <div
                style={{ display: "inline-block", position: "relative" }}
                className="w-full"
              >
                <input
                  type="text"
                  placeholder="Картой"
                  className={clsx([
                    currentInput === "Card"
                      ? "outline outline-primary border-primary"
                      : "border-slate-300",
                    "px-5 py-3 text-right placeholder-slate-300 text-slate-600 relative bg-white bg-white rounded text-lg border outline-none focus:outline-none focus:ring w-full pl-10",
                  ])}
                  value={card}
                  onChange={(e) => {
                    setCard(parseFloat((e.target as HTMLInputElement).value));
                  }}
                  disabled
                />
              </div>
            </div>
            <div className="relative flex w-full flex-wrap items-stretch mb-3 z-10">
              <button
                className="absolute inset-0 bg-transparent z-20"
                onClick={() => {
                  setCurrentInput("Cash");
                  if (card == total) {
                    setCard(0);
                    setCash(total);
                  }
                }}
              />
              <span className="z-10 flex flex-col items-center justify-center h-full leading-snug font-normal absolute text-center text-slate-300 absolute bg-transparent rounded text-base items-center justify-center w-8 pl-3 py-3">
                <svg
                  style={{ color: "rgb(82, 209, 35)" }}
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  className="h-6 w-6"
                  viewBox="0 0 16 16"
                >
                  <path
                    d="M8 10a2 2 0 1 0 0-4 2 2 0 0 0 0 4z"
                    fill="#52d123"
                  ></path>{" "}
                  <path
                    d="M0 4a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H1a1 1 0 0 1-1-1V4zm3 0a2 2 0 0 1-2 2v4a2 2 0 0 1 2 2h10a2 2 0 0 1 2-2V6a2 2 0 0 1-2-2H3z"
                    fill="#52d123"
                  ></path>{" "}
                </svg>
              </span>
              <div className="w-full relative">
                <input
                  type="text"
                  placeholder="Наличными"
                  className={clsx([
                    currentInput === "Cash"
                      ? "outline outline-primary border-primary"
                      : "border-slate-300",
                    "px-5 py-3 text-right placeholder-slate-300 text-slate-600 relative bg-white bg-white rounded text-lg border outline-none focus:outline-none focus:ring w-full pl-10",
                  ])}
                  value={cash}
                  onChange={(e) => {
                    setCash(parseFloat((e.target as HTMLInputElement).value));
                  }}
                  disabled
                />
              </div>
            </div>
            <div className="w-1/2 font-light">
              Сдача: <span className="font-bold">{surplus}</span>
            </div>
            {!isChangeAvailable && (
              <p className="text-red-600">Введена слишком большая сумма</p>
            )}
            {!loading && !link && (
              <div className="w-full flex flex-col items-start">
                <div className="flex flex-col items-start space-y-4">
                  <button
                    className="mt-5 bg-teal-500 hover:bg-primary disabled:bg-primary/50 text-white font-bold font-xl py-2 px-4 border-b-4 border-teal-700 hover:border-teal-500 rounded"
                    disabled={
                      btnLoading ||
                      total > cash + card ||
                      total - cash < -20000 ||
                      !selectedPager
                    }
                    onClick={handlePaymentSubmit}
                  >
                    Оплатить
                  </button>
                  {hasPagers && (
                    <span className="">
                      Выберите пейджер перед подтверждением оплаты:
                    </span>
                  )}
                </div>
                {hasPagers && (
                  <div className="w-full grid grid-cols-5 gap-2 mt-2">
                    {pagers.map((pager) => (
                      <button
                        className={clsx([
                          "h-14 text-white flex flex-col items-center justify-center rounded-md text-lg font-bold",
                          selectedPager === pager
                            ? "bg-primary ring ring-primary/80"
                            : "bg-neutral-400 hover:bg-primary/80",
                        ])}
                        onClick={() => {
                          setSelectedPager(pager);
                        }}
                      >
                        {pager}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
            {loading ? (
              <div className="w-full flex flex-col items-start space-y-3">
                <div className="w-full flex items-center justify-center h-[162px] pb-8">
                  <Oval
                    height={80}
                    width={80}
                    color="#3eb2b2"
                    wrapperStyle={{}}
                    wrapperClass=""
                    visible={true}
                    ariaLabel="oval-loading"
                    secondaryColor="#4acfcf"
                    strokeWidth={2}
                    strokeWidthSecondary={2}
                  />
                </div>
                <button
                  type="button"
                  onClick={() => {
                    clearInterval(timer);
                    setLoading(false);
                  }}
                  className="px-8 py-2 text-gray-400 text-lg flex items-center justify-center rounded-md border-2 border-gray-400 hover:bg-gray-400 hover:text-white"
                >
                  Отмена
                </button>
              </div>
            ) : link ? (
              <div className="w-full flex flex-col items-start">
                <a
                  href={link}
                  target="_blank"
                  rel="noreferrer"
                  className="text-indigo-500 text-sm hover:underline font-xl my-4"
                >
                  ПОСМОТРЕТЬ ЧЕК
                </a>
                <button
                  onClick={() => {
                    printCheck(printData).then((res) => {
                      // console.log(res);
                    });
                  }}
                  className="bg-teal-500 hover:bg-teal-400 text-white font-bold font-xl py-2 px-4 border-b-4 border-teal-700 hover:border-teal-500 rounded"
                  disabled={!printData}
                >
                  Распечатать
                </button>
              </div>
            ) : (
              <></>
            )}
          </div>
        </div>
      </div>
    </TerminalModalLayout>
  );
};

export default PaymentModal;
