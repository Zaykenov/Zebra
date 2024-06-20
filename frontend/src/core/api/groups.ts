import { backendApiInstance } from "./index";

export type InventoryGroupData = {
  id?: number;
  name: string;
  sklad_id: number;
  sklad_name: string;
  measure: string;
  type: string;
  items: InventoryGroupItemData[];
};

export type InventoryGroupItemData = {
  item_id: number;
  type?: string;
};

export const getAllInventoryGroups = async (skladId?: number) => {
  const queryParam = !!skladId ? `?sklad_id=${skladId}` : "";
  const res = await backendApiInstance.get(
    "sklad/inventarization/group/getAll" + queryParam
  );
  return res.data;
};

export const getInventoryGroup = async (id: string) => {
  const res = await backendApiInstance.get(
    `sklad/inventarization/group/get/${id}`
  );
  return res.data;
};

export const createInventoryGroup = async (data: InventoryGroupData) => {
  const res = await backendApiInstance.post(
    "sklad/inventarization/group/create",
    data
  );
  return res.data;
};

export const updateInventoryGroup = async (data: InventoryGroupData) => {
  const res = await backendApiInstance.post(
    "sklad/inventarization/group/update",
    data
  );
  return res.data;
};

export const deleteMenuItem = async (data: { id: number }) => {
  const res = await backendApiInstance.post("item/delete", data);
  return res.data;
};

export const getProductsByCategory = async (data: { categoryId: number }) => {
  const res = await backendApiInstance.get(
    `item/getWithParams?category=${data.categoryId}`
  );
  return res.data;
};
