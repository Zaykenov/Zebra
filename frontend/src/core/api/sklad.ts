import { backendApiInstance } from "./index";
import { StorageData } from "../../components/__modules/StorageForm/types";

const processFormData = (data: StorageData) => {
  const postData = {
    name: data.name,
    address: data.address,
  };

  return data.id ? { id: data.id, ...postData } : postData;
};

export const getAllSklads = async () => {
  const res = await backendApiInstance.get("sklad/getAll");
  return res.data;
};

export const getSklad = async (id: string) => {
  const res = await backendApiInstance.get(`sklad/get?id=${id}`);
  return res.data;
};

export const createSklad = async (data: any) => {
  const res = await backendApiInstance.post(
    "sklad/create",
    processFormData(data)
  );
  return res.data;
};

export const updateSklad = async (data: any) => {
  const res = await backendApiInstance.post(
    "sklad/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteSklad = async (data: { id: number }) => {
  const res = await backendApiInstance.post("sklad/delete", data);
  return res.data;
};
