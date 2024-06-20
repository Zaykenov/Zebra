import React, { FC, useEffect, useMemo, useState } from "react";
import ShopSelect from "../../__common/FilterSelect/ShopSelect";
import { dateToString } from "@api/check";
import { formatNumber } from "@utils/formatNumber";
import clsx from "clsx";
import { ResponsiveLine } from "@nivo/line";
import MainLayout from "@layouts/MainLayout";
import { useFilter } from "@context/filter.context";
import useMasterRole from "@hooks/useMasterRole";
import {
  getDaysOfWeekStats,
  getHourlyStats,
  getPaymentsStats,
  getSalesEveryDay,
  getSalesEveryMonth,
  getSalesEveryWeek,
  getSalesToday,
} from "@api/stats";
import { btns, statsMapping, weekDaysLabels } from "./constants";
import BarChartBox from "./BarChartBox";
import { Period, Stats } from "./types";

const SalesContent: FC = () => {
  const { queryOptions } = useFilter();

  const isMaster = useMasterRole();

  const [todayStats, setTodayStats] = useState<any>(null);
  const [totalStats, setTotalStats] = useState<any>(null);

  const [shop, setShop] = useState<number | null>(null);

  const [period, setPeriod] = useState<Period>(Period.DAY);
  const [statsType, setStatsType] = useState<Stats>(Stats.REVENUE);
  const [tooltipPayload, setTooltipPayload] = useState<{
    name: string;
    isCurrency: boolean;
  } | null>({
    name: "Прибыль",
    isCurrency: true,
  });
  const [statsData, setStatsData] = useState<any[]>([]);

  const [paymentsStats, setPaymentsStats] = useState<
    { name: string; value: number; value2: number }[]
  >([]);
  const [weekDayStats, setWeekDayStats] = useState<any[]>([]);
  const [hourlyStats, setHourlyStats] = useState<any[]>([]);

  useEffect(() => {
    getSalesToday(shop ? { shop } : undefined).then((res) => {
      setTodayStats(res.data);
    });
  }, [shop]);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getSalesEveryMonth(queryOptions).then((res) => {
      setTotalStats({
        revenue: formatNumber(res.data.total_revenue, true, true),
        profit: formatNumber(res.data.total_profit, true, true),
        checks: formatNumber(res.data.total_checks, false, false),
        visitors: formatNumber(res.data.total_visitors, false, false),
        avg_check: formatNumber(res.data.total_avg_check, true, true),
      });
    });
    getPaymentsStats(queryOptions).then((res) => {
      setPaymentsStats([
        {
          name: "Карточка",
          value: res.data.total_card,
          value2: res.data.total_check_card,
        },
        {
          name: "Наличные",
          value: res.data.total_cash,
          value2: res.data.total_check_cash,
        },
      ]);
    });
    getDaysOfWeekStats(queryOptions).then((res) => {
      setWeekDayStats(
        Object.entries(res.data).map(([key, value]) => ({
          // @ts-ignore
          name: weekDaysLabels[key],
          ...(value as object),
        }))
      );
    });
    getHourlyStats(queryOptions).then((res) => {
      setHourlyStats(
        res.data.map((stat: any, idx: number) => ({
          name: `${idx}`,
          ...stat,
        }))
      );
    });
  }, [queryOptions]);

  const weekDaySelectedData = useMemo(
    () =>
      weekDayStats.map((stat) => ({
        name: stat.name,
        value: stat[statsType],
      })),
    [statsType, weekDayStats]
  );

  const hourlySelectedData = useMemo(
    () =>
      hourlyStats.map((stat) => ({
        name: stat.name,
        value: stat[statsType],
      })),
    [statsType, hourlyStats]
  );

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    if (period === Period.DAY) {
      getSalesEveryDay(queryOptions).then((res) => {
        const data = {
          id: statsType,
          data: res.data.statistics.map((stat: any) => ({
            x: dateToString(stat.from, false, true),
            y: stat[statsType],
          })),
        };
        setStatsData([data]);
      });
    } else if (period === Period.WEEK) {
      getSalesEveryWeek(queryOptions).then((res) => {
        const data = {
          id: statsType,
          data: res.data.statistics.map((stat: any) => ({
            x: `${dateToString(stat.from, false, true)} - ${dateToString(
              stat.to,
              false,
              true
            )}`,
            y: stat[statsType],
          })),
        };
        setStatsData([data]);
      });
    } else {
      getSalesEveryMonth(queryOptions).then((res) => {
        const data = {
          id: statsType,
          data: res.data.statistics.map((stat: any) => ({
            x: `${dateToString(stat.from, false, true)} - ${dateToString(
              stat.to,
              false,
              true
            )}`,
            y: stat[statsType],
          })),
        };
        setStatsData([data]);
      });
    }
  }, [period, statsType, queryOptions]);

  return (
    <MainLayout
      title="Статистика продаж"
      dateFilter
      customBtns={
        isMaster
          ? []
          : [
              <ShopSelect
                onChange={(value) => {
                  setShop(value[0]);
                }}
              />,
            ]
      }
    >
      <div className="px-5 py-3 w-full flex flex-col overflow-x-hidden">
        <div className="px-6 py-5 mb-10 bg-primary/10 border border-primary rounded flex items-center">
          <div className="text-lg font-semibold leading-5 mr-36">
            Сегодня,
            <br />
            {dateToString(new Date(Date.now()).toISOString(), false, true)}
          </div>
          <div className="grow flex items-center justify-around">
            <div className="flex flex-col">
              {todayStats && (
                <span className="text-lg font-bold">
                  {formatNumber(todayStats.revenue, true, true)}
                  <span
                    className={`text-xs ${
                      todayStats.percent_revenue >= 0
                        ? "text-green-500"
                        : "text-red-500"
                    }`}
                  >
                    {`${todayStats.percent_revenue >= 0 ? "+" : ""}${
                      todayStats.percent_revenue
                    }%`}
                  </span>
                </span>
              )}
              <span className="text-xs">выручка</span>
            </div>

            <div className="flex flex-col">
              {todayStats && (
                <span className="text-lg font-bold">
                  {formatNumber(todayStats.profit, true, true)}
                  <span
                    className={`text-xs ${
                      todayStats.percent_profit >= 0
                        ? "text-green-500"
                        : "text-red-500"
                    }`}
                  >
                    {`${todayStats.percent_profit >= 0 ? "+" : ""}${
                      todayStats.percent_profit
                    }%`}
                  </span>
                </span>
              )}
              <span className="text-xs">прибыль</span>
            </div>
            {todayStats && (
              <div className="flex flex-col">
                <span className="text-lg font-bold">
                  {todayStats.checks}
                  <span
                    className={`text-xs ${
                      todayStats.percent_checks >= 0
                        ? "text-green-500"
                        : "text-red-500"
                    }`}
                  >
                    {`${todayStats.percent_checks >= 0 ? "+" : ""}${
                      todayStats.percent_checks
                    }%`}
                  </span>
                </span>
                <span className="text-xs">
                  {todayStats.checks
                    ? todayStats.checks % 10 === 1
                      ? "чек"
                      : todayStats.checks % 10 < 5
                      ? "чека"
                      : "чеков"
                    : "чеков"}
                </span>
              </div>
            )}
            {todayStats && (
              <div className="flex flex-col">
                <span className="text-lg font-bold">
                  {todayStats.visitors}
                  <span
                    className={`text-xs ${
                      todayStats.percent_visitors >= 0
                        ? "text-green-500"
                        : "text-red-500"
                    }`}
                  >
                    {`${todayStats.percent_visitors >= 0 ? "+" : ""}${
                      todayStats.percent_visitors
                    }%`}
                  </span>
                </span>
                <span className="text-xs">
                  {todayStats.visitors
                    ? todayStats.visitors > 1
                      ? "посетителей"
                      : "посетитель"
                    : "посетителей"}
                </span>
              </div>
            )}
            {todayStats && (
              <div className="flex flex-col">
                <span className="text-lg font-bold">
                  {formatNumber(todayStats.avg_check, true, true)}
                  <span
                    className={`text-xs ${
                      todayStats.percent_avg_check >= 0
                        ? "text-green-500"
                        : "text-red-500"
                    }`}
                  >
                    {`${todayStats.percent_avg_check >= 0 ? "+" : ""}${
                      todayStats.percent_avg_check
                    }%`}
                  </span>
                </span>
                <span className="text-xs">средний чек</span>
              </div>
            )}
          </div>
        </div>
        <div className="flex flex-col shadow-md mb-10">
          <div className="flex items-center space-x-3 mb-4">
            <span className="text-lg font-semibold">
              {statsMapping[statsType]}
            </span>
            <span className="isolate inline-flex items-center rounded-md shadow-sm text-xs mt-0.5">
              <button
                type="button"
                onClick={() => {
                  setPeriod(Period.DAY);
                }}
                className={clsx([
                  "relative inline-flex items-center rounded-l-md border border-primary px-2 py-1 text-gray-700 hover:underline focus:z-10 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary",
                  period === Period.DAY
                    ? "bg-primary text-gray-50"
                    : "bg-white text-primary",
                ])}
              >
                День
              </button>
              <button
                type="button"
                onClick={() => {
                  setPeriod(Period.WEEK);
                }}
                className={clsx([
                  "relative -ml-px inline-flex items-center border border-primary px-2 py-1 text-gray-700 hover:underline focus:z-10 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary",
                  period === Period.WEEK
                    ? "bg-primary text-gray-50"
                    : "bg-white text-primary font-medium",
                ])}
              >
                Неделя
              </button>
              <button
                type="button"
                onClick={() => {
                  setPeriod(Period.MONTH);
                }}
                className={clsx([
                  "relative -ml-px inline-flex items-center rounded-r-md border border-primary bg-white px-2 py-1 text-gray-700 hover:underline focus:z-10 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary",
                  period === Period.MONTH
                    ? "bg-primary text-gray-50"
                    : "bg-white text-primary font-medium",
                ])}
              >
                Месяц
              </button>
            </span>
          </div>
          <div className="flex h-[400px]">
            <ResponsiveLine
              data={statsData}
              margin={{ top: 10, right: 110, bottom: 50, left: 60 }}
              xScale={{ type: "point" }}
              yScale={{
                type: "linear",
                min: 0,
                max: "auto",
                stacked: true,
                reverse: false,
              }}
              enableGridX={false}
              enableSlices="x"
              sliceTooltip={({ slice }) => {
                return (
                  <div className="bg-white border border-gray-500 px-3 py-2 flex flex-col">
                    <span className="text-xs font-medium leading-3">
                      {slice.points[0]?.data.xFormatted}
                    </span>
                    <span className="font-bold text-sm">
                      {formatNumber(
                        parseFloat(
                          slice.points[0]?.data.yFormatted as string
                        ) as number
                      )}
                    </span>
                  </div>
                );
              }}
              axisBottom={{
                ...{
                  tickValues:
                    period === Period.DAY
                      ? statsData[0]?.data
                          .filter((_: any, idx: number) => idx % 5 === 0)
                          .map((stat: any) => stat.x)
                      : "",
                },
                tickSize: 5,
                tickPadding: 12,
              }}
              yFormat=" >-.2f"
              colors={{
                scheme: "accent",
              }}
              pointSize={5}
              pointColor={{ theme: "background" }}
              pointBorderWidth={2}
              pointBorderColor={{ from: "serieColor" }}
              pointLabelYOffset={-12}
              useMesh={true}
              enableArea={true}
              animate={true}
            />
          </div>
          <div className="flex divide-x divide-gray-200">
            {btns.map((btn) => (
              <button
                type="button"
                onClick={() => {
                  setStatsType(btn.value);
                  setTooltipPayload({
                    name: btn.text,
                    isCurrency: btn.isCurrency,
                  });
                }}
                className={clsx([
                  "grow py-4 flex flex-col items-center justify-center",
                  statsType === btn.value
                    ? "bg-white"
                    : "bg-gray-50 shadow-inner border-t-2 !border-t-gray-400",
                ])}
              >
                <span className="text-lg font-bold">
                  {totalStats && totalStats[btn.value]}
                </span>
                <span className="text-xs">{btn.text}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="grid grid-cols-2 gap-10">
          <BarChartBox
            title="Методы оплаты"
            vertical
            switchable
            data={paymentsStats}
          />
          <BarChartBox
            title="По дням недели"
            data={weekDaySelectedData}
            tooltipPayload={tooltipPayload}
          />
          <BarChartBox
            title="По времени"
            data={hourlySelectedData}
            className="col-span-2"
            tooltipPayload={tooltipPayload}
          />
        </div>
      </div>
    </MainLayout>
  );
};

export default SalesContent;
