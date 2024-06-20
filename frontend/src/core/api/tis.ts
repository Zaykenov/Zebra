import axios from "axios";
import { v4 as uuidv4 } from "uuid";
import { backendApiInstance } from "./index";

// dev
const BASE_URL = "https://dev.kassa.wipon.kz/api/";
const TOKEN = "qoGuMj6ylslMLnuCZVTRlJBZEK2zIW";

// prod
// const BASE_URL = "https://app.kassa.wipon.kz/api/";
// const TOKEN = "f0b08oH8PKa31FJinYvD9hLdhU8CZc";

export enum TisType {
  SELL = 2,
  RETURN = 3,
}

export enum TisCompareFieldType {
  VENDOR = "vendor_code",
  BARCODE = "barcode",
}

export enum TisPaymentMethod {
  CASH = 0,
  CARD = 1,
}

export type TisData = {
  token?: string;
  type: TisType;
  items: {
    name: string;
    price: number;
    quantity: number;
    discount: number;
    kgd_code: number;
    compare_field: {
      type: TisCompareFieldType;
      value: string;
    };
  }[];
  payments: {
    payment_method: TisPaymentMethod;
    sum: number;
  }[];
};

export const measureToKgd = (measure: "кг" | "л" | "шт.") => {
  switch (measure) {
    case "кг":
      return 166;
    case "л":
      return 112;
    case "шт.":
      return 796;
    default:
      return 796;
  }
};

export const axiosTisInstance = axios.create({
  baseURL: BASE_URL,
});

export const sendTisData = async (data: TisData) => {
  const res = await axiosTisInstance.post(
    "ticket/send",
    {
      ...data,
      token: TOKEN,
    },
    {
      headers: {
        Accept: "application/json",
        "Idempotency-Key": uuidv4(),
      },
    }
  );
  return res.data;
};

export const saveTisResponse = async (data: any) => {
  const res = await backendApiInstance.post("external/saveCheck", data);
  return res.data;
};
