import React, { FC, ReactNode, useState } from "react";
import Sidebar from "../../__modules/Sidebar";
import clsx from "clsx";
import FilterProvider from "../../../context/filter.context";

export interface PageLayoutProps {
  children?: ReactNode;
  defaultFilters?: {
    hasPagination?: boolean;
    hasDatePicker?: boolean;
  };
}

const openWidth = "w-60";
const shrinkWidth = "w-[50px]";

const openPadding = "pl-60";
const shrinkPadding = "pl-[50px]";

const PageLayout: FC<PageLayoutProps> = ({ children, defaultFilters }) => {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <FilterProvider defaultFilters={defaultFilters}>
      <Sidebar
        sidebarOpen={sidebarOpen}
        setSidebarOpen={setSidebarOpen}
        width={sidebarOpen ? openWidth : shrinkWidth}
      />
      <div
        className={clsx([
          sidebarOpen ? openPadding : shrinkPadding,
          "flex flex-1 flex-col transition-padding duration-300",
        ])}
      >
        {children}
      </div>
    </FilterProvider>
  );
};

export default PageLayout;
