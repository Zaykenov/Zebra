import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type TransferItemData = {
  item_id: number | string;
  type: string;
  quantity: number | string;
  transfer_id?: number | string;
  measure?: string;
};

export type TransferData = {
  time: string;
  date?: Date;
  from_sklad: number;
  to_sklad: number;
  item_transfers: TransferItemData[];
  id?: number;
};

export const getAllTransfers = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "sklad/transfer/getAll" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getTransferById = async (id: string) => {
  const res = await backendApiInstance.get(`sklad/transfer/get/${id}`);
  return res.data;
};

export const createTransfer = async (data: TransferData) => {
  const res = await backendApiInstance.post("sklad/transfer/create", data);
  return res.data;
};

export const updateTransfer = async (data: TransferData) => {
  const res = await backendApiInstance.post("sklad/transfer/update", data);
  return res.data;
};

export const deleteTransfer = async (data: { id: number }) => {
  const res = await backendApiInstance.post("sklad/transfer/delete", data);
  return res.data;
};
