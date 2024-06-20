import { backendApiInstance } from "./index";

export const getAllProductCategories = async () => {
  const res = await backendApiInstance.get("item/category/getAll");
  return res.data;
};

export const getProductCategory = async (id: string) => {
  const res = await backendApiInstance.get(`item/category/get?id=${id}`);
  return res.data;
};

export const createProductCategory = async (data: any) => {
  const res = await backendApiInstance.post("item/category/create", data);
  return res.data;
};

export const updateProductCategory = async (data: any) => {
  const res = await backendApiInstance.post("item/category/update", data);
  return res.data;
};

export const deleteProductCategory = async (data: any) => {
  const res = await backendApiInstance.post("/item/category/delete", data);
};
