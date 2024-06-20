import { Menu, Transition } from "@headlessui/react";
import React, { FC, Fragment } from "react";
import { Bars3Icon } from "@heroicons/react/24/outline";
import { useRouter } from "next/router";
import { useAuth } from "../../../context/auth.context";
import Link from "next/link";

const navigation = [
  { name: "Поставка", href: "/terminal/supply" },
  { name: "Списание", href: "/terminal/waste" },
  { name: "Начать смену", href: "/terminal/shift" },
  { name: "Закрыть смену", href: "/terminal/close-shift" },
  { name: "Инкассация", href: "/terminal/collection" },
  { name: "Внесение", href: "/terminal/income" },
  { name: "Инвентаризация", href: "/terminal/inventory" },
];

interface HeaderMenuProps {
  shiftStarted: boolean | undefined
}

const HeaderMenu: FC<HeaderMenuProps> = ({shiftStarted}) => {
  const router = useRouter();

  const { logOut } = useAuth();

  return (
    <div className="">
      <Menu as="div" className="relative inline-block text-left">
        <div>
          <Menu.Button className="inline-flex w-full justify-center rounded-md bg-black bg-opacity-20 px-4 py-2 text-sm font-medium text-white hover:bg-opacity-30 focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75">
            <Bars3Icon className="w-5 h-5" />
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
            <div className="px-1 py-1 ">
              {navigation.map((item) => {
                if((!shiftStarted && item.name!=="Закрыть смену") || shiftStarted && item.name!=="Начать смену")
                return (
                  <Menu.Item key={item.name}>
                    <Link href={item.href}>
                      <a
                        className={`${
                          item.href === router.pathname
                            ? "bg-gray-900 text-white font-medium"
                            : "text-gray-900 hover:bg-gray-200 "
                        } group flex w-full items-center rounded-md px-2 py-2`}
                      >
                        {item.name}
                      </a>
                    </Link>
                  </Menu.Item>
                )
              })}
            </div>
            <div className="px-1 py-1">
              <button
                onClick={() => {
                  logOut && logOut();
                  router.push("/login");
                }}
                className="w-full bg-red-500/80 text-sm text-white hover:bg-red-600/80 hover:text-white px-3 py-2 rounded-md font-medium"
              >
                Выйти
              </button>
            </div>
          </Menu.Items>
        </Transition>
      </Menu>
    </div>
  );
};

export default HeaderMenu;
