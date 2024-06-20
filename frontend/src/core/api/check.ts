import { backendApiInstance } from "./index";
import axios from "axios";
import { QueryOptions, QueryOptionsData, queryOptionsToString } from "./index";

export enum PaymentMethod {
  CASH = "Наличными",
  CARD = "Картой",
  MIXED = "Смешанно",
}

export type CheckItemData = {
  id: number;
  name: string;
  count: number;
  price: number;
  hasDiscount: boolean;
  total: number;
  comments: string;
  type: "tovar" | "techCart";
  selectedModificators?: {
    id: number;
    nabor_id: number;
    name: string;
    price: number;
    quantity: number;
  }[];
};

export type CheckData = {
  id?: number | null;
  worker?: number;
  opened_at?: string;
  closed_at?: string;
  sum: number;
  discount: number;
  discount_percent: number;
  cost: number;
  payment: PaymentMethod;
  status?: "opened" | "closed";
  comment?: string;
  user_id?: number | null;
  tovarCheck: {
    tovar_id: number;
    id?: number;
    tovar_name?: string;
    quantity: number;
    cost: number;
    price: number;
    comment: string;
    Modifications: string;
  }[];
  techCartCheck: {
    tech_cart_id: number;
    id?: number;
    name?: string;
    quantity: number;
    cost: number;
    price: number;
    comment: string;
    modificators: {
      id: number;
      nabor_id: number;
      name: string;
      brutto: number;
    }[];
  }[];
};

export enum CheckStatus {
  OPEN = "opened",
  CLOSE = "closed",
}

export type ModificatorCheckData = {
  id: number;
  nabor_id: number;
  name?: string;
  cost?: number;
  brutto?: number;
};

export type TovarCheckData = {
  tovar_id: number;
  name?: string;
  quantity: number;
  cost?: number;
  price?: number;
  discount?: number;
  modifications: string;
  comments: string;
};

export type TechCartCheckData = {
  tech_cart_id: number;
  name?: string;
  quantity: number;
  discount?: number;
  cost?: number;
  price?: number;
  modificators: ModificatorCheckData[];
  comments: string;
};

export type CheckPostData = {
  id?: number;
  user_id?: number;
  worker?: number;
  opened_at?: string;
  closed_at?: string;
  sum?: number;
  cost?: number;
  status: CheckStatus;
  payment: PaymentMethod;
  cash: number;
  card: number;
  discount?: number;
  discount_percent: number;
  comment: string;
  tovarCheck: TovarCheckData[];
  techCartCheck: TechCartCheckData[];
};

export const dateToString = (
  dateStr: string,
  onlyTime: boolean,
  onlyDate?: boolean
) => {
  const date = new Date(dateStr);
  const nameOfMonth = date.toLocaleString("ru", {
    month: "long",
  });
  const day = `${date.getDate()} ${nameOfMonth}`;
  const time = date.toTimeString().slice(0, 5);
  return onlyTime ? time : onlyDate ? day : `${day}, ${time}`;
};

export const getAllChecks = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "check/getAll" + queryOptionsToString(queryOptions)
  );
  return queryOptions?.page ? res.data.data : res.data;
};

export const deleteCheckById = async (data: { id: number }) => {
  const res = await backendApiInstance.post("check/delete", data);
  return res.data;
};

export const getAllWorkerChecks = async () => {
  const res = await backendApiInstance.get("check/getAllWorker");
  return res.data;
};

export const getCheckById = async (data: { id: number }) => {
  const res = await backendApiInstance.get(`check/get/${data.id}`);
  return res.data;
};

export const getCheckForPrint = async (data: { id: number }) => {
  const res = await backendApiInstance.get(`check/getForPrinter/${data.id}`);
  return res.data;
};

export const createCheck = async (
  data: CheckPostData,
  idempotencyKey: string
) => {
  try {
    const res = await backendApiInstance.post("check/create", data, {
      headers: {
        "Idempotency-Key": idempotencyKey,
      },
    });
    return res.data;
  } catch (e) {
    throw e;
  }
};

export const deactivateCheck = async (id: number) => {
  const res = await backendApiInstance.get(`check/deactivate/${id}`);
  return res.data;
};

export const printCheck = async (data: string) => {
  const res = await axios.post("/terminal/order/print-receipt", {
    data,
  });
  return res.data;
};
