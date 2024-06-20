import { Disclosure } from "@headlessui/react";
import React, { FC, ReactNode, useEffect, useState } from "react";
import clsx from "clsx";

export type NavigationItem = {
  id: number;
  name: string;
  component: ReactNode;
};

export interface TerminalCabinetLayoutProps {
  navigation: NavigationItem[];
}

const TerminalCabinetLayout: FC<TerminalCabinetLayoutProps> = ({
  navigation,
}) => {
  const [active, setActive] = useState<NavigationItem | null>(null);

  useEffect(() => {
    navigation && setActive(navigation[0]);
  }, [navigation]);

  return (
    <div className="h-full w-full flex flex-col">
      <Disclosure as="nav" className="bg-gray-900 shadow-sm">
        {({ open }) => (
          <>
            <div className="w-full px-5">
              <div className="flex h-12 justify-between">
                <div className="flex">
                  <div className="flex space-x-8">
                    {active &&
                      navigation.map((item) => (
                        <button
                          onClick={() => {
                            setActive(item);
                          }}
                          key={item.name}
                          className={clsx(
                            item.id === active.id
                              ? "border-b-2 border-white text-white"
                              : "border-transparent text-gray-300 hover:text-white hover:border-gray-300",
                            "inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                          )}
                          aria-current={
                            item.id === active.id ? "page" : undefined
                          }
                        >
                          {item.name}
                        </button>
                      ))}
                  </div>
                </div>
              </div>
            </div>

            <Disclosure.Panel className="sm:hidden">
              <div className="space-y-1 pt-2 pb-3">
                {active &&
                  navigation.map((item) => (
                    <Disclosure.Button
                      key={item.name}
                      onClick={() => {
                        setActive(item);
                      }}
                      className={clsx(
                        item.id === active.id
                          ? "bg-indigo-50 border-indigo-500 text-indigo-700"
                          : "border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800",
                        "block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                      )}
                      aria-current={item.id === active.id ? "page" : undefined}
                    >
                      {item.name}
                    </Disclosure.Button>
                  ))}
              </div>
            </Disclosure.Panel>
          </>
        )}
      </Disclosure>

      <div className="overflow-auto min-h-full pb-10">
        {active && active.component}
      </div>
    </div>
  );
};

export default TerminalCabinetLayout;
