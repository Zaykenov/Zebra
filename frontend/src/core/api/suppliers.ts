import { backendApiInstance } from "./index";
import { SupplierData } from "../../components/__modules/SupplierForm/types";

const processFormData = (data: SupplierData) => {
  const postData = {
    name: data.name,
    phone: data.phone,
    address: data.address,
    comment: data.comment,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllSuppliers = async () => {
  const res = await backendApiInstance.get("dealer/getAll");
  return res.data;
};

export const getSupplier = async (id: string) => {
  const res = await backendApiInstance.get(`dealer/get?id=${id}`);
  return res.data;
};

export const createSupplier = async (data: any) => {
  const res = await backendApiInstance.post(
    "dealer/create",
    processFormData(data)
  );
  return res.data;
};

export const updateSupplier = async (data: any) => {
  const res = await backendApiInstance.post(
    "dealer/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteSupplier = async (data: { id: number }) => {
  const res = await backendApiInstance.post("dealer/delete", data);
  return res.data;
};
