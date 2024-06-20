import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export const getAllOstatki = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "sklad/ostatki" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};
