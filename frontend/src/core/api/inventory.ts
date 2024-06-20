import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type InventoryData = {
  id?: number;
  inventarization_id?: number;
  sklad_id?: number;
  sklad?: string;
  time: string;
  date?: Date;
  type: string;
  result?: number;
  status: string;
  items: InventoryItem[];
};

export type InventoryItem = {
  id?: number;
  inventarization_id?: number;
  item_id: number;
  sklad_id?: number;
  status?: string;
  time: string;
  type: string;
  start_quantity?: number;
  expenses?: number;
  income?: number;
  removed?: number;
  removed_sum?: number;
  plan_quantity?: number;
  fact_quantity?: number | string;
  difference?: number;
  difference_sum?: number;
};

export const getAllInventory = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "sklad/inventarization/getAll" + queryOptionsToString(queryOptions, false),
  );
  return res.data;
};

export const getInventoryById = async (id: string | number) => {
  const res = await backendApiInstance.get(`sklad/inventarization/get/${id}`);
  return res.data;
};

export const createInventory = async (data: InventoryData) => {
  const res = await backendApiInstance.post(
    "sklad/inventarization/create",
    data,
  );
  return res.data;
};

export const updateInventory = async (data: InventoryData) => {
  const res = await backendApiInstance.post(
    "sklad/inventarization/update",
    data,
  );
  return res.data;
};

export const updateInventoryParams = async (data: InventoryData) => {
  const res = await backendApiInstance.post(
    "sklad/inventarization/updateParams",
    data,
  );
  return res.data;
};

export const deleteInventory = async (data: { id: number | string }) => {
  const res = await backendApiInstance.post(
    `sklad/inventarization/delete`,
    data,
  );
  return res.data;
};

export const deleteInventoryItem = async (data: { id: number | string }) => {
  const res = await backendApiInstance.post(
    `sklad/inventarization/deleteItem`,
    data,
  );
  return res.data;
};

export const getInventoryIncomeDetails = async (id: number) => {
  const res = await backendApiInstance.get(
    `sklad/inventarization/getDetails/income/${id}`,
  );
  return res.data;
};

export const getInventoryWasteDetails = async (id: number) => {
  const res = await backendApiInstance.get(
    `sklad/inventarization/getDetails/spisanie/${id}`,
  );
  return res.data;
};

export const getInventoryExpenseDetails = async (id: number) => {
  const res = await backendApiInstance.get(
    `sklad/inventarization/getDetails/expence/${id}`,
  );
  return res.data;
};
