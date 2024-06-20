import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import clsx from "clsx";
import {
  Bars3Icon,
  ChartBarIcon,
  ChevronLeftIcon,
  ClipboardDocumentListIcon,
  CurrencyDollarIcon,
  LockClosedIcon,
  Square3Stack3DIcon,
} from "@heroicons/react/24/outline";
import { Disclosure, Transition } from "@headlessui/react";
import Link from "next/link";
import { useRouter } from "next/router";
import { useAuth } from "../../../context/auth.context";
import { UsersIcon } from "@heroicons/react/24/outline";
export interface SidebarProps {
  sidebarOpen: boolean;
  setSidebarOpen: Dispatch<SetStateAction<boolean>>;
  width: string;
}

export enum NavSection {
  STATS = "statistics",
  FINANCE = "finances",
  MENU = "menu",
  STORAGE = "storage",
  ACCESS = "access",
  MARKETING = "marketing",
}

const navigation = [
  {
    name: NavSection.STATS,
    text: "Статистика",
    icon: ChartBarIcon,
    links: [
      { name: "Продажи", href: "/sales" },
      { name: "Сотрудники", href: "/cashiers" },
      { name: "Чеки", href: "/receipts" },
      { name: "Оплаты", href: "/payments" },
      { name: "ABC анализ", href: "/abc" },
    ],
    list: ["/sales", "/cashiers", "/receipts", "/payments", "/abc"],
    height: "h-44",
  },
  {
    name: NavSection.FINANCE,
    text: "Финансы",
    icon: CurrencyDollarIcon,
    links: [
      { name: "Транзакции", href: "/transactions" },
      { name: "Кассовые смены", href: "/shifts" },
      { name: "Счета", href: "/accounts" },
    ],
    list: ["/transactions", "/shifts", "/accounts"],
    height: "h-24",
  },
  {
    name: NavSection.MENU,
    text: "Меню",
    icon: ClipboardDocumentListIcon,
    links: [
      { name: "Товары", href: "/menu" },
      { name: "Тех. карты", href: "/dishes" },
      { name: "Наборы модификаторов", href: "/nabors" },
      { name: "Ингредиенты", href: "/ingredients" },
      { name: "Категории товаров и тех. карт", href: "/categories_products" },
      { name: "Категории ингредиентов", href: "/categories_ingredients" },
    ],
    list: [
      "/menu",
      "/dishes",
      "/nabors",
      "/ingredients",
      "/categories_products",
      "/categories_ingredients",
    ],
    height: "h-68",
  },
  {
    name: NavSection.STORAGE,
    text: "Склад",
    icon: Square3Stack3DIcon,
    links: [
      { name: "Остатки", href: "/calculations" },
      { name: "Поставки", href: "/supply" },
      { name: "Перемещения", href: "/transfer" },
      { name: "Списания", href: "/waste" },
      { name: "Отчет по движению", href: "/reports" },
      { name: "Группировка", href: "/groups" },
      { name: "Инвентаризации", href: "/inventory" },
      { name: "Поставщики", href: "/suppliers" },
      { name: "Склады", href: "/storages" },
    ],
    list: [
      "/calculations",
      "/supply",
      "/transfer",
      "/waste",
      "/reports",
      "/groups",
      "/inventory",
      "/suppliers",
      "/storages",
    ],
    height: "h-58",
  },
  {
    name: NavSection.ACCESS,
    text: "Доступ",
    icon: LockClosedIcon,
    links: [
      { name: "Сотрудники", href: "/workers" },
      { name: "Заведения", href: "/shops" },
    ],
    list: ["/workers", "/shops"],
    height: "h-20",
  },
  {
    name: NavSection.MARKETING,
    text: "Маркетинг",
    icon: UsersIcon,
    links: [
      { name: "Клиенты", href: "/clients" },
      { name: "Отзывы", href: "/feedbacks" },
    ],
    list: ["/users", "/feedbacks"],
    height: "h-20",
  },
];

const Sidebar: FC<SidebarProps> = ({ sidebarOpen, setSidebarOpen, width }) => {
  const router = useRouter();

  const { logOut } = useAuth();

  const [section, setSection] = useState<NavSection | null>(null);

  useEffect(() => {
    const activeSection = navigation.filter((item) =>
      item.list.includes("/" + router.pathname.split("/")[1]),
    )[0];
    setSection(activeSection?.name);
  }, [router]);

  return (
    <div
      className={clsx([
        width,
        "overflow-hidden fixed inset-y-0 flex flex-col transition-width duration-200 font-inter",
      ])}
    >
      <div className="flex flex-grow flex-col overflow-y-auto overflow-x-hidden border-r border-gray-200 bg-menu">
        <div className="flex flex-shrink-0 items-center px-4 border-b border-gray-200 h-[70px]">
          <button onClick={() => setSidebarOpen((prevState) => !prevState)}>
            <div className="h-5 w-5 mr-3 text-gray-500 hover:text-gray-800 relative">
              <ChevronLeftIcon
                className={clsx([
                  sidebarOpen ? "opacity-100 scale-100" : "opacity-0 scale-0",
                  "absolute inset-0 transform transition-all duration-300 ",
                ])}
              />
              <Bars3Icon
                className={clsx([
                  sidebarOpen ? "opacity-0 scale-0" : "opacity-100 scale-125",
                  "absolute inset-0 transform transition-all duration-300",
                ])}
              />
            </div>
          </button>
          <div className="flex items-center text-lg uppercase text-gray-600 font-inter font-black tracking-wider">
            <span className="whitespace-nowrap">Zebra CRM</span>
          </div>
        </div>
        <div className="pt-2 flex flex-grow flex-col">
          <nav className="flex-1 px-2 pb-4">
            {navigation.map((item) => (
              <Disclosure key={item.name}>
                <Disclosure.Button
                  onClick={() =>
                    section === item.name
                      ? setSection(null)
                      : setSection(item.name)
                  }
                  className={clsx(
                    "text-gray-600 hover:text-gray-900",
                    "group w-full flex items-center p-2 font-medium text-[15px]",
                  )}
                >
                  <item.icon
                    className={clsx("mr-3 flex-shrink-0 h-5 w-5")}
                    aria-hidden="true"
                  />
                  {item.text}
                </Disclosure.Button>
                <Transition
                  show={sidebarOpen && section === item.name}
                  enter="transition-height duration-300 ease-out"
                  enterFrom="h-0"
                  enterTo={item.height}
                  leave="transition-height duration-200 ease-out"
                  leaveFrom={item.height}
                  leaveTo="h-0"
                  className="overflow-hidden"
                >
                  <Disclosure.Panel>
                    <ul>
                      {item.links.map((link, idx) => (
                        <li key={idx} className="">
                          <Link href={link.href}>
                            <a
                              className={clsx([
                                `/${router.pathname.split("/")[1]}` ===
                                link.href
                                  ? "bg-indigo-100 text-gray-700 font-semibold"
                                  : "text-blue-700 hover:text-blue-900",
                                "pl-10 pr-5 py-1.5 block text-[15px] ",
                              ])}
                            >
                              {link.name}
                            </a>
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </Disclosure.Panel>
                </Transition>
              </Disclosure>
            ))}
          </nav>
          {/* <button
            onClick={() => {
              logOut && logOut();
              router.push("/guide");
            }}
            className="py-2 border-gray-300 border-t text-green-600 font-small"
          >
            Инструкция по использованию
          </button> */}
          <button
            onClick={() => {
              logOut && logOut();
              router.push("/login");
            }}
            className="py-2 border-t border-gray-300 text-red-600 font-medium"
          >
            Выйти
          </button>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
