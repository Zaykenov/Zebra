import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type PaymentData = {
  time: string;
  check_count: number;
  cash: number;
  card: number;
  total: number;
};

const processFormData = (data: PaymentData) => {
  const postData = {
    time: data.time,
    check_count: data.check_count,
    cash: data.cash,
    card: data.card,
    total: data.total,
  };
  return postData;
};

export const getAllPayments = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "statistics/payments" + queryOptionsToString(queryOptions)
  );
  return res.data;
};
