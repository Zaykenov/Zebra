import React, { useCallback, useState } from "react";
import { NextPage } from "next";
import TerminalLayout from "@layouts/TerminalLayout";
import { Input } from "@shared/ui/Input";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import {
  closeShift,
  ShiftCategory,
  ShiftPostData,
} from "@api/shifts";
import { cleanCache } from "@utils/cleanCache";
import axios, { AxiosError } from "axios";
import AlertMessage from "@common/AlertMessage";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import { IdempotencyKey, sendIdempotencyKeys } from "@api/idempotencyKey";

const TerminalCloseShiftPage: NextPage = () => {
  const router = useRouter();

  const [loading, setLoading] = useState(false);
  const {alertMessage, showAlertMessage, hideAlertMessage} = useAlertMessage()
  const { handleSubmit, register } = useForm<ShiftPostData>({
    defaultValues: {
      category: ShiftCategory.OPEN,
      sum: "",
      comment: "",
    },
  });

  const onSubmit = useCallback(
    (submitData: ShiftPostData) => {
      setLoading(true);
      let keys: IdempotencyKey[] = JSON.parse(localStorage.getItem("zebra.idempotencyKeys") || "[]")
      console.log(keys)
      localStorage.removeItem("zebra.idempotencyKeys")
      sendIdempotencyKeys(keys).then((res)=>{
        console.log(res)
      })
      closeShift({
        ...submitData,
      })
        .then(() => {
          cleanCache();
          return router.push("/terminal/order");
        })
        .catch((e: Error | AxiosError) => {
          setLoading(false)
          if (axios.isAxiosError(e)) {
            if (e.response?.status === 444) {
              showAlertMessage('Закройте все чеки прежде чем закрывать смену', AlertMessageType.ERROR)
              router.push("/terminal/shift");
            } showAlertMessage(e.message, AlertMessageType.ERROR);
          } showAlertMessage(e.message, AlertMessageType.ERROR);
        }
        );
    },
    [router]
  );

  return (
    <TerminalLayout>
      <div className="w-full flex items-center justify-center">
        <div className="w-full max-w-2xl p-4 rounded bg-gray-100 shadow-2xl border border-gray-300 flex flex-col">
          <div className="text-lg font-medium mb-4">Закрыть смену</div>
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="flex flex-col space-y-5"
          >
            <Input
              {...register("sum")}
              placeholder="Остаток наличных после инкассации"
              className="w-2/3"
            />
            <Input
              {...register("comment")}
              placeholder="Комментарий"
              className="w-2/3"
            />
            <div className="pt-5 border-t border-gray-200">
              <div className="flex items-center justify-between">
                <button
                  disabled={loading}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    router.back();
                  }}
                  className="px-4 py-2 bg-transparent hover:bg-gray-300 rounded text-gray-500 hover:text-gray-900"
                >
                  Отмена
                </button>
                <button
                  disabled={loading}
                  type="submit"
                  className="disabled:bg-primary/50 px-8 py-2 bg-primary hover:opacity-80 rounded text-white font-medium"
                >
                  Закончить
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
      {alertMessage && (
        <AlertMessage message={alertMessage.message} type={alertMessage.type} onClose={hideAlertMessage}/>
      )}
    </TerminalLayout>
  );
};

export default TerminalCloseShiftPage;
