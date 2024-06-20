import React, { FC, Fragment, useEffect, useState } from "react";
import { Disclosure, Menu, Transition } from "@headlessui/react";
import clsx from "clsx";
import {
  AdjustmentsHorizontalIcon,
  Bars3Icon,
  InformationCircleIcon,
  UserMinusIcon,
  UserPlusIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { useRouter } from "next/router";
import Link from "next/link";
import HeaderMenu from "@common/HeaderMenu/HeaderMenu";
import { checkShift } from "@api/shifts";
import { getDateAndTime } from "@utils/dateFormatter";
import { clearStorage } from "@utils/clearStorage";
import useLocalStorage from "@hooks/useLocalStorage";

const navigation = [
  { name: "Новый заказ", href: "/terminal/order" },
  { name: "Заказы", href: "/terminal/orders" },
  { name: "Архив чеков", href: "/terminal/archive" },
];

export type UserShiftData = {
  id: number;
  shopId: number;
  shopName: string;
  workerId: number;
  workerName: string;
  isClosed: boolean;
  createdAt: string;
};

const TerminalHeader: FC = () => {
  const router = useRouter();
  const [userData, setUserData] = useState<UserShiftData>();
  const [isAutomaticPrint, setIsAutomaticPrint] = useLocalStorage(
    "zebra.automaticPrint",
    false,
  );
  const [hasPagers, setHasPagers] = useLocalStorage("zebra.hasPagers", true);

  const handleAutomaticPrint = () => {
    setIsAutomaticPrint(!isAutomaticPrint);
    router.reload();
  };

  const handleHasPagers = () => {
    setHasPagers(!hasPagers);
    router.reload();
  };

  const terminalSettings = [
    {
      optionName: "Автоматическая печать чека",
      inputStateValue: isAutomaticPrint,
      inputStateHandler: handleAutomaticPrint,
    },
    {
      optionName: "Включить пейджеры",
      inputStateValue: hasPagers,
      inputStateHandler: handleHasPagers,
    },
  ];

  useEffect(() => {
    checkShift().then((res) => {
      const userShiftData: UserShiftData = {
        id: res.data.id,
        shopId: res.data.shop_id,
        isClosed: res.data.is_closed,
        createdAt: res.data.created_at,
        workerId: res.data.worker_id,
        workerName: res.data.worker,
        shopName: res.data.shop_name,
      };
      setUserData(userShiftData);
    });
  }, [router]);

  return (
    <Disclosure as="nav" className="bg-gray-800">
      {({ open }) => (
        <>
          <div className="px-5">
            <div className="flex h-16 items-center justify-between">
              <div className="w-full flex items-center">
                <div className="hidden w-full md:flex items-center justify-between">
                  <div className="flex items-baseline space-x-4">
                    {navigation.map((item) => (
                      <Link key={item.name} href={item.href}>
                        <a
                          className={clsx(
                            item.href === router.pathname
                              ? "bg-gray-900 text-white"
                              : "text-gray-300 hover:bg-gray-700 hover:text-white",
                            "px-3 py-2 rounded-md text-sm font-medium",
                          )}
                          aria-current={
                            item.href === router.pathname ? "page" : undefined
                          }
                        >
                          {item.name}
                        </a>
                      </Link>
                    ))}
                  </div>
                  <div className="flex items-center space-x-3">
                    <Menu as="div" className="relative inline-block text-left">
                      <div>
                        <Menu.Button className="inline-flex w-full justify-center rounded-md bg-black bg-opacity-20 px-4 py-2 text-sm font-medium text-white hover:bg-opacity-30 focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75">
                          {userData?.isClosed ? (
                            <UserMinusIcon className="w-5 h-5 text-red-500" />
                          ) : (
                            <UserPlusIcon className="w-5 h-5 text-green-500" />
                          )}
                        </Menu.Button>
                      </div>
                      <Transition
                        as={Fragment}
                        enter="transition ease-out duration-100"
                        enterFrom="transform opacity-0 scale-95"
                        enterTo="transform opacity-100 scale-100"
                        leave="transition ease-in duration-75"
                        leaveFrom="transform opacity-100 scale-100"
                        leaveTo="transform opacity-0 scale-95"
                      >
                        <Menu.Items className="absolute z-50 right-0 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                          <table className="text-xs my-3">
                            <tbody>
                              <tr>
                                <td className="px-2 py-2 text-gray-500 font-semibold">
                                  Заведение
                                </td>
                                <td className="px-2 py-2">
                                  {userData?.shopName}
                                </td>
                              </tr>
                              <tr>
                                <td className="px-2 py-2 text-gray-500 font-semibold">
                                  Пользователь
                                </td>
                                <td className="px-2 py-2">
                                  {userData?.workerName}
                                </td>
                              </tr>
                              <tr>
                                <td className="px-2 py-2 text-gray-500 font-semibold">
                                  Смена открыта
                                </td>
                                <td className="px-2 py-2">
                                  {userData?.isClosed ? "Нет" : "Да"}
                                </td>
                              </tr>
                              <tr>
                                <td className="px-2 py-2 text-gray-500 font-semibold">
                                  Время открытия смены
                                </td>

                                <td className="px-2 py-2">
                                  {
                                    //@ts-ignore
                                    getDateAndTime(userData?.createdAt)
                                  }
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </Menu.Items>
                      </Transition>
                    </Menu>
                    <Menu as="div" className="relative inline-block text-left">
                      <div>
                        <Menu.Button className="inline-flex w-full justify-center rounded-md bg-black bg-opacity-20 px-4 py-2 text-sm font-medium text-white hover:bg-opacity-30 focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75">
                          <AdjustmentsHorizontalIcon className="w-5 h-5 text-gray-500" />
                        </Menu.Button>
                      </div>
                      <Transition
                        as={Fragment}
                        enter="transition ease-out duration-100"
                        enterFrom="transform opacity-0 scale-95"
                        enterTo="transform opacity-100 scale-100"
                        leave="transition ease-in duration-75"
                        leaveFrom="transform opacity-100 scale-100"
                        leaveTo="transform opacity-0 scale-95"
                      >
                        <Menu.Items className="absolute z-50 right-0 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                          <table className="text-xs my-3">
                            <tbody>
                              {terminalSettings.map((setting, idx) => (
                                <tr key={idx}>
                                  <td className="px-2 py-2 text-gray-500 font-semibold">
                                    {setting.optionName}
                                  </td>
                                  <td className="px-2 py-2">
                                    <div>
                                      <label className="relative inline-flex items-center cursor-pointer">
                                        <input
                                          type="checkbox"
                                          value=""
                                          className="sr-only peer"
                                          checked={setting.inputStateValue}
                                          onChange={setting.inputStateHandler}
                                        />
                                        <div className="w-11 h-6 bg-gray-200 rounded-full peer peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                                      </label>
                                    </div>
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </Menu.Items>
                      </Transition>
                    </Menu>
                    <a
                      href="https://drive.google.com/file/d/124idL79xltgrFEwOu-vYHffhq7Xu6eR_/view?usp=sharing"
                      className="bg-black/20 hover:bg-black/30 px-4 py-2 rounded-md"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <InformationCircleIcon className="w-5 h-5 text-white" />
                    </a>
                    <HeaderMenu shiftStarted={!userData?.isClosed}/>
                    <button
                      onClick={() => {
                        clearStorage();
                        router.reload();
                      }}
                      className="bg-primary/70 hover:bg-primary text-gray-300 hover:text-white px-3 py-2 rounded-md font-medium text-sm"
                    >
                      {/*<ArrowPathIcon className="w-5 h-5" />*/}
                      Обновить меню
                    </button>
                  </div>
                </div>
              </div>
              <div className="-mr-2 flex md:hidden">
                {/* Mobile menu button */}
                <Disclosure.Button className="inline-flex items-center justify-center rounded-md bg-gray-800 p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                  <span className="sr-only">Open main menu</span>
                  {open ? (
                    <XMarkIcon className="block h-6 w-6" aria-hidden="true" />
                  ) : (
                    <Bars3Icon className="block h-6 w-6" aria-hidden="true" />
                  )}
                </Disclosure.Button>
              </div>
            </div>
          </div>

          <Disclosure.Panel className="md:hidden">
            <div className="space-y-1 px-2 pt-2 pb-3 sm:px-3">
              {navigation.map((item) => (
                <Link key={item.name} href={item.href}>
                  <Disclosure.Button
                    as="a"
                    className={clsx(
                      item.href === router.pathname
                        ? "bg-gray-900 text-white"
                        : "text-gray-300 hover:bg-gray-700 hover:text-white",
                      "block px-3 py-2 rounded-md text-base font-medium",
                    )}
                    aria-current={
                      item.href === router.pathname ? "page" : undefined
                    }
                  >
                    {item.name}
                  </Disclosure.Button>
                </Link>
              ))}
            </div>
          </Disclosure.Panel>
        </>
      )}
    </Disclosure>
  );
};

export default TerminalHeader;
