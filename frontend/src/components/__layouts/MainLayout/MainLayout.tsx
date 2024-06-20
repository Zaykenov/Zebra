import React, {
  FC,
  forwardRef,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import Link from "next/link";
import { ChevronLeftIcon } from "@heroicons/react/24/outline";
import { useRouter } from "next/router";
import DatePicker from "react-datepicker";
import clsx from "clsx";
import { dateToString } from "@api/check";
import { Select } from "antd";
import { getMasterShops } from "@api/master";
import { useFilter } from "@context/index";
import Search from "@common/Search";
import { MainLayoutProps } from "@layouts/MainLayout/types";
import { filterComponents } from "@layouts/MainLayout/constants";
import Pagination from "@modules/Pagination";
import useMasterRole from "@hooks/useMasterRole";

const MainLayout: FC<MainLayoutProps> = ({
  title,
  addHref,
  backBtn = false,
  dateFilter = false,
  searchFilter = false,
  pagination = false,
  children,
  customBtns,
  content,
  withoutOverflow = false,
  excelDownloadButton,
  filterOptions,
}) => {
  const router = useRouter();

  const isMaster = useMasterRole();

  const { dateRange, handleDateChange, totalResults } = useFilter();

  const [shopOptions, setShopOptions] = useState<
    {
      label: string;
      value: number;
      data: {
        token: string;
      };
    }[]
  >([]);
  const [selectedShop, setSelectedShop] = useState<null | {
    label: string;
    value: number;
    data?: {
      token: string;
    };
  }>(null);

  const onChange = useCallback(
    (dates: [Date | null, Date | null]) => {
      const [startDate, endDate] = dates;
      handleDateChange({
        startDate,
        endDate,
      });
    },
    [handleDateChange]
  );

  useEffect(() => {
    if (isMaster === null) return;
    if (isMaster) {
      getMasterShops().then((res) => {
        setShopOptions([
          {
            label: "Все заведения",
            value: 0,
          },
          ...res.data.map((shop: any) => ({
            label: shop.name,
            value: shop.id,
            data: {
              token: shop.token,
            },
          })),
        ]);
      });
      const selectedShopString = localStorage.getItem(
        "zebra.master.selectedShop"
      );
      if (!selectedShopString) {
        localStorage.setItem(
          "zebra.master.mainToken",
          localStorage.getItem("zebra.authToken") || ""
        );
        return;
      }
      setSelectedShop(JSON.parse(selectedShopString));
    }
  }, [isMaster]);

  return (
    <div className="flex flex-col h-screen w-full">
      <header className="h-[70px] border-b border-gray-200 flex items-center justify-between px-5">
        <div className="text-lg font-semibold flex items-center space-x-4">
          {backBtn && (
            <button
              onClick={() => {
                router.back();
              }}
              className="hover:bg-gray-200 rounded-md"
            >
              <ChevronLeftIcon className="text-indigo-500 h-8 w-8 font-bold" />
            </button>
          )}
          <span>{title}</span>
        </div>
        <div className="flex items-center space-x-3">
          {isMaster && (
            <div className="w-48">
              <Select
                className="w-full"
                value={selectedShop ? selectedShop.value : 0}
                onChange={(_, option) => {
                  if (window !== undefined) {
                    localStorage.setItem(
                      "zebra.master.selectedShop",
                      JSON.stringify(option)
                    );
                    // @ts-ignore
                    if (!option.data)
                      localStorage.setItem(
                        "zebra.authToken",
                        localStorage.getItem("zebra.master.mainToken") || ""
                      );
                    else
                      localStorage.setItem(
                        "zebra.authToken",
                        // @ts-ignore
                        option.data.token
                      );
                    router.reload();
                  }
                }}
                options={shopOptions}
                showSearch
                filterOption={(input, option) =>
                  (option?.label ?? "")
                    .trim()
                    .toLowerCase()
                    .includes(input.trim().toLowerCase())
                }
              />
            </div>
          )}
          {customBtns && (
            <div className="flex items-center space-x-3">{customBtns}</div>
          )}
          {content}
          {dateFilter && (
            <div className="">
              <DatePicker
                locale="ru"
                renderCustomHeader={({
                  monthDate,
                  customHeaderCount,
                  decreaseMonth,
                  increaseMonth,
                }) => (
                  <div>
                    <button
                      aria-label="Previous Month"
                      className={clsx([
                        "react-datepicker__navigation react-datepicker__navigation--previous",
                        customHeaderCount === 1 && "hidden invisible",
                      ])}
                      onClick={decreaseMonth}
                    >
                      <span
                        className={
                          "react-datepicker__navigation-icon react-datepicker__navigation-icon--previous mt-2"
                        }
                      >
                        {"<"}
                      </span>
                    </button>
                    <span className="font-medium font-inter text-sm capitalize">
                      {monthDate.toLocaleDateString("default", {
                        month: "long",
                      })}{" "}
                      {monthDate.getFullYear()}
                    </span>
                    <button
                      aria-label="Next Month"
                      className={clsx([
                        "react-datepicker__navigation react-datepicker__navigation--next",
                        customHeaderCount === 0 && "hidden invisible",
                      ])}
                      onClick={increaseMonth}
                    >
                      <span
                        className={
                          "react-datepicker__navigation-icon react-datepicker__navigation-icon--next mt-2"
                        }
                      >
                        {">"}
                      </span>
                    </button>
                  </div>
                )}
                selected={dateRange?.startDate}
                onChange={onChange}
                startDate={dateRange?.startDate}
                endDate={dateRange?.endDate}
                monthsShown={2}
                customInput={<CustomDateInput />}
                excludeDateIntervals={[
                  {
                    start: new Date(),
                    end: new Date(
                      new Date().setFullYear(new Date().getFullYear() + 1)
                    ),
                  },
                ]}
                selectsRange
              />
            </div>
          )}
          {!!excelDownloadButton && (
            <button
              onClick={excelDownloadButton}
              className="text-white shadow-md pt-1.5 pb-2 px-3 bg-primary text-sm font-semibold rounded-md hover:bg-teal-600"
            >
              Экспорт в Excel
            </button>
          )}
          {!!addHref && (
            <Link href={addHref}>
              <button className="text-white shadow-md pt-1.5 pb-2 px-3 bg-primary text-sm font-semibold rounded-md hover:bg-teal-600">
                Добавить
              </button>
            </Link>
          )}
        </div>
      </header>
      <main
        className={clsx(["flex-1 w-full", !withoutOverflow && "overflow-auto"])}
      >
        {!!filterOptions && (
          <div className="flex items-center space-x-5 p-3">
            {searchFilter && <Search />}
            <div className="flex items-center space-x-3">
              {filterOptions.map(
                (filterOption) => filterComponents[filterOption]
              )}
            </div>
          </div>
        )}
        {children}
        {pagination && <Pagination detailed />}
      </main>
    </div>
  );
};

const CustomDateInput = forwardRef(
  ({ value, onClick }: { value?: any; onClick?: any }, ref) => {
    const datesFromValue: string[] = useMemo(
      () =>
        !!value
          ? value
              .replace(/\s/g, "")
              .split("-")
              .map((dateString: string) =>
                dateToString(dateString, false, true)
              )
          : [],
      [value]
    );

    const fromDateString = useMemo(
      () => datesFromValue[0] || "",
      [datesFromValue]
    );
    const toDateString = useMemo(
      () => datesFromValue[1] || "",
      [datesFromValue]
    );

    return (
      <button
        className="border-2 border-emerald-500 hover:border-emerald-700 px-3 py-0.5 rounded text-emerald-700 font-medium"
        onClick={onClick}
        // @ts-ignore
        ref={ref}
      >
        {fromDateString} -{" "}
        {!toDateString.toLowerCase().includes("nan") ? toDateString : ""}
      </button>
    );
  }
);

export default MainLayout;
