import React, { useEffect } from "react";
import { NextPage } from "next";
import { useForm } from "react-hook-form";
import { signIn } from "@api/workers";
import { useRouter } from "next/router";
import { useAuth } from "../../context/auth.context";
import useAlertMessage, { AlertMessageType } from "../../hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage";

const LoginPage: NextPage = () => {
  const router = useRouter();
  const { setAuthInfo } = useAuth();
  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();
  const { register, handleSubmit } = useForm({
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const onSubmit = (submitData: { username: string; password: string }) => {
    signIn(submitData)
      .then((res) => {
        setAuthInfo && setAuthInfo(res.data);
        localStorage.setItem("zebra.authed", "true");
        router.reload();
      })
      .catch(() => {
        showAlertMessage("Неверные данные", AlertMessageType.ERROR);
      });
  };

  return (
    <div className="flex h-screen flex-col justify-center bg-slate-100 py-12 sm:px-6 lg:px-8">
      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-700"
              >
                Логин
              </label>
              <div className="mt-1">
                <input
                  {...register("username", { required: true })}
                  type="text"
                  className="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                />
              </div>
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700"
              >
                Пароль
              </label>
              <div className="mt-1">
                <input
                  {...register("password", { required: true })}
                  type="password"
                  className="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md border border-transparent bg-primary py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-primary/70 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                Войти
              </button>
            </div>
          </form>
        </div>
      </div>
      {alertMessage && (
        <AlertMessage
          message={alertMessage.message}
          type={alertMessage.type}
          onClose={hideAlertMessage}
        />
      )}
    </div>
  );
};

export default LoginPage;
