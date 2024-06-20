import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export const getAllTrafficReports = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "sklad/getTrafficReport" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};
