import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export const getSalesToday = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/today" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getSalesEveryDay = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/everyDay" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getSalesEveryWeek = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/everyWeek" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getSalesEveryMonth = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/everyMonth" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getWorkersStats = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/workers" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getPaymentsStats = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/payments" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getDaysOfWeekStats = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/daysOfTheWeek" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getHourlyStats = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/statByHour" + queryOptionsToString(queryOptions)
  );
  return res.data;
};
