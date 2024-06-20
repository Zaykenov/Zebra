import { backendApiInstance } from "./index";
import { AccountData } from "../../components/__modules/AccountForm/types";

const processFormData = (data: AccountData) => {
  const postData = {
    name: data.name,
    currency: data.currency,
    type: data.type,
    start_balance: data.start_balance,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllAccounts = async () => {
  const res = await backendApiInstance.get("finance/schet/getAll");
  return res.data;
};

export const getAccount = async (id: string) => {
  const res = await backendApiInstance.get(`finance/schet/get?id=${id}`);
  return res.data;
};

export const createAccount = async (data: AccountData) => {
  const res = await backendApiInstance.post("finance/schet/create", data);
  return res.data;
};

export const updateAccount = async (data: AccountData) => {
  const res = await backendApiInstance.post("finance/schet/update", data);
  return res.data;
};

export const deleteAccount = async (data: { id: number }) => {
  const res = await backendApiInstance.post("finance/schet/delete", data);
  return res.data;
};
