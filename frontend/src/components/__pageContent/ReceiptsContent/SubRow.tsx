import React, { FC, ReactNode, useState } from "react";
import { Row } from "@tanstack/react-table";
import { Tab } from "@headlessui/react";
import clsx from "clsx";
import CheckOverview from "./CheckOverview";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [tabs] = useState<{ name: string; component: ReactNode }[]>([
    {
      name: "Счет",
      component: <CheckOverview rowData={rowData} />,
    },
    {
      name: "Списания",
      component:
        rowData.original.status === "Возврат" ? (
          <></>
        ) : (
          <CheckOverview rowData={rowData} withIngredients />
        ),
    },
  ]);

  return (
    <div className="px-4">
      <Tab.Group>
        <Tab.List className="flex space-x-1 rounded-xl mt-1 rounded-t">
          {tabs.map((tab) => (
            <Tab
              key={tab.name}
              className={({ selected }) =>
                clsx(
                  "p-2 text-sm hover:underline",
                  selected ? "bg-white" : "text-blue-500 hover:text-black",
                )
              }
            >
              {tab.name}
            </Tab>
          ))}
        </Tab.List>
        <Tab.Panels className="bg-white p-6 mb-1">
          {tabs.map((tab, idx) => (
            <Tab.Panel key={idx} className="">
              {tab.component}
            </Tab.Panel>
          ))}
        </Tab.Panels>
      </Tab.Group>
    </div>
  );
};

export default SubRow;
