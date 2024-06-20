import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import { getCheckForPrint } from "@api/check";
import { useRouter } from "next/router";
import Image from "next/image";
import useTimer from "@hooks/useTimer";
import { formatSeconds } from "@utils/dateFormatter";
import { DocumentTextIcon } from "@heroicons/react/24/solid";

const ReceiptPage: NextPage = () => {
  const router = useRouter();
  const [receiptReady, setReceiptReady] = useState<boolean>(false);
  const [redirectUrl, setRedirectUrl] = useState<string>("");
  const { seconds, isExpired } = useTimer(300);

  const checkReceiptStatus = async () => {
    try {
      const id = router.query.id;
      if (!id) return;
      const response = await getCheckForPrint({ id: parseInt(id as string) });
      if (
        !(
          response.data.tisCheckUrl === "" ||
          response.data.tisCheckUrl === undefined
        )
      ) {
        setReceiptReady(true);
        setRedirectUrl(response.data.tisCheckUrl);
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const handleRedirect = () => {
    window.location.href = redirectUrl;
  };

  useEffect(() => {
    if (redirectUrl !== "") {
      handleRedirect();
    }
  }, [receiptReady]);

  useEffect(() => {
    checkReceiptStatus();
    const intervalId = setInterval(checkReceiptStatus, 10000);
    return () => clearInterval(intervalId);
  }, [router]);

  return (
    <div className="overflow-auto h-screen">
      <div className="mt-20 relative p-4 text-center bg-white sm:p-5">
        <div className="mt-20 flex flex-col items-center">
          {!receiptReady ? (
            <div>
              <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                Ваш фискальный чек все еще обрабатывается.
              </p>
              <Image src={"/images/jdun-meme.gif"} height={200} width={200} />
              {!isExpired && (
                <div>
                  <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                    Максимальное время ожидания{" "}
                    <strong>{formatSeconds(seconds)}</strong>
                  </p>
                  <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                    По истечению времени вас перекинет на страницу чека{" "}
                  </p>
                  <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                    Спасибо за ожидание
                  </p>
                </div>
              )}
              <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl"></p>
            </div>
          ) : (
            <div className="flex flex-col items-center">
              <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                Чек готов. Если вас не перевело на страницу чека, нажмите на
                иконку
              </p>
              <DocumentTextIcon
                width={200}
                height={200}
                onClick={handleRedirect}
              />
              <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                Или перейдите по ссылке:
              </p>
              <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
                {redirectUrl ? redirectUrl : "link fetch error"}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ReceiptPage;
