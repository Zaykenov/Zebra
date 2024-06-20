import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";

export type WasteItemData = {
  item_id: number | string;
  type: string;
  quantity: number | string;
  details: string;
  measure?: string;
};

export type WasteData = {
  sklad_id?: number;
  reason: string;
  comment: string;
  items: WasteItemData[];
  time: string;
  date?: Date;
  id?: number;
};

const processFormData = (data: WasteData) => {
  const postData = {
    sklad_id: data.sklad_id,
    reason: data.reason,
    comment: data.comment,
    items: data.items,
    time: data.time,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllWastes = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "sklad/getRemoved" + queryOptionsToString(queryOptions)
  );
  return res.data;
};

export const getWasteById = async (id: string) => {
  const res = await backendApiInstance.get(`sklad/getRemoved/${id}`);
  return res.data;
};

export const createWaste = async (data: WasteData) => {
  const res = await backendApiInstance.post("sklad/remove", data);
  return res.data;
};

export const updateWaste = async (data: WasteData) => {
  const res = await backendApiInstance.post("sklad/spisanie/update", data);
  return res.data;
};

export const deleteWaste = async (data: { id: number }) => {
  const res = await backendApiInstance.post("sklad/spisanie/delete", data);
  return res.data;
};

export const confirmRemoveWaste = async (id: number) => {
  const res = await backendApiInstance.post(`sklad/confirm/${id}`);
  return res.data;
};

export const rejectRemoveWaste = async (id: number) => {
  const res = await backendApiInstance.post(`sklad/reject/${id}`);
  return res.data;
};
