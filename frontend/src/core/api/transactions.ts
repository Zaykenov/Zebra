import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type TransactionData = {
  id?: number;
  worker_id?: number;
  shift_id?: number;
  schet_id: number;
  status?: string;
  time?: string;
  date?: Date | null;
  category: string;
  sum: number;
  comment: string;
};

export enum TransactionCategory {
  OPEN_SHIFT = "openShift",
  CLOSE_SHIFT = "closeShift",
  POSTAVKA = "postavka",
  COLLECTION = "collection",
  INCOME = "income",
}

export const mapTransactionCategoryToString = (category: TransactionCategory) =>
  ({
    [TransactionCategory.OPEN_SHIFT]: "Открытие смены",
    [TransactionCategory.CLOSE_SHIFT]: "Закрытие смены",
    [TransactionCategory.POSTAVKA]: "Поставки",
    [TransactionCategory.COLLECTION]: "Инкассация",
    [TransactionCategory.INCOME]: "Внесение",
  }[category]);

export const getAllTransactions = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "transaction/getAll" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getTransaction = async (id: string) => {
  const res = await backendApiInstance.get(`transaction/get/${id}`);
  return res.data;
};

export const createTransaction = async (data: TransactionData) => {
  const res = await backendApiInstance.post("transaction/create", data);
  return res.data;
};

export const updateTransaction = async (data: TransactionData) => {
  const res = await backendApiInstance.post("transaction/update", data);
  return res.data;
};

export const deleteTransaction = async (data: { id: number }) => {
  const res = await backendApiInstance.post(`transaction/delete/${data.id}`);
  return res.data;
};
