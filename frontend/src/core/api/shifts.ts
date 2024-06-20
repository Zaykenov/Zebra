import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type ShiftTransactionData = {
  id: number;
  shift_id: number;
  worker_id: number;
  category: string;
  time: string;
  sum: number;
  comment: string;
};

export type ShiftData = {
  id: number;
  created_at: string;
  closed_at: string;
  begin_sum: number;
  end_sum_fact: number;
  end_sum_plan: number;
  expense: number;
  cash: number;
  card: number;
  collection: number;
  difference: number;
  transactions: ShiftTransactionData[];
  is_closed: boolean;
};

export enum ShiftCategory {
  OPEN = "openShift",
  CLOSE = "closeShift",
  COLLECTION = "collection",
  SUPPLY = "postavka",
  INCOME = "income",
}

export const shiftMapping = {
  [ShiftCategory.OPEN]: "Открытие смены",
  [ShiftCategory.CLOSE]: "Закрытие смены",
  [ShiftCategory.COLLECTION]: "Инкассация",
  [ShiftCategory.SUPPLY]: "Поставка",
  [ShiftCategory.INCOME]: "Внесение",
};

export type ShiftPostData = {
  schet_id?: number;
  category: ShiftCategory;
  sum: string | number;
  comment: string;
};

export const getAllShifts = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "shift/getAll" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const checkShift = async () => {
  const res = await backendApiInstance.get(`shift/check`);
  return res.data;
};

export const getShiftById = async (id: number) => {
  const res = await backendApiInstance.get(`shift/get/${id}`);
  return res.data;
};

export const openShift = async (data: any) => {
  const postData: ShiftPostData = {
    schet_id: data.schet_id,
    category: ShiftCategory.OPEN,
    sum: parseInt(data.sum),
    comment: data.comment,
  };
  const res = await backendApiInstance.post("transaction/create", postData);
  return res.data;
};

export const closeShift = async (data: any) => {
  const postData: ShiftPostData = {
    category: ShiftCategory.CLOSE,
    sum: parseInt(data.sum),
    comment: data.comment,
  };
  const res = await backendApiInstance.post("transaction/create", postData);
  return res.data;
};

export const createTransaction = async (data: any) => {
  const postData: ShiftPostData = {
    category: data.category,
    sum: parseInt(data.sum),
    comment: data.comment,
  };
  const res = await backendApiInstance.post("transaction/create", postData);
  return res.data;
};
