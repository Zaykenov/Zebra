import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type SupplyItem = {
  item_id: number;
  type: "ingredient" | "tovar";
  quantity: number;
  cost: number;
};

export type SupplyData = {
  dealer_id?: number;
  sklad_id?: number;
  schet_id?: number;
  time?: string;
  date?: Date;
  items: SupplyItem[];
};

const processFormData = (data: SupplyData) => {
  const postData = {
    ...data,
    time: data.date?.toISOString() || new Date().toISOString(),
  };
  return postData;
};

export const getAllSupplies = async (
  queryOptions?: QueryOptionsData,
  withDate: boolean = true
) => {
  const res = await backendApiInstance.get(
    "sklad/postavka/getAll/" + queryOptionsToString(queryOptions, withDate)
  );
  return res.data;
};

export const getSupply = async (id: string) => {
  const res = await backendApiInstance.get(`sklad/postavka/get/${id}`);
  return res.data;
};

export const createSupply = async (data: any) => {
  const res = await backendApiInstance.post(
    "sklad/postavka/create",
    processFormData(data)
  );
  return res.data;
};

export const createSupplyAsWorker = async (data: any) => {
  try {
    const res = await backendApiInstance.post(
      "sklad/postavka/createWorker",
      processFormData(data)
    );
    return res.data;
  } catch (e) {
    throw e;
  }
};

export const updateSupply = async (data: any) => {
  const res = await backendApiInstance.post("sklad/postavka/update", data);
  return res.data;
};

export const deleteSupply = async (data: { id: number }) => {
  const res = await backendApiInstance.post("sklad/postavka/delete", data);
  return res.data;
};
