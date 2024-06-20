import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";
import { MenuItemData } from "@modules/MenuItemForm/types";

const processFormData = (data: MenuItemData) => {
  const postData = {
    name: data.name,
    category: data.category,
    image: data.image,
    tax: "Фискальный налог",
    measure: data.measure,
    discount: data.discount,
    cost: data.cost,
    margin: data.margin,
    price: data.price,
    shop_id: data.shop_id,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllMenuItems = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "item/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMasterMenuItems = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "master/item/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMenuItem = async (id: string) => {
  const getRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/item/get/"
      : "item/get?id=";
  const res = await backendApiInstance.get(`${getRoute}${id}`);
  return res.data;
};

export const createMenuItem = async (data: MenuItemData) => {
  const res = await backendApiInstance.post("item/create", processFormData(data));
  return res.data;
};

export const updateMenuItem = async (data: MenuItemData) => {
  const res = await backendApiInstance.post("item/update", processFormData(data));
  return res.data;
};

export const updateMasterMenuItem = async (data: MenuItemData) => {
  const res = await backendApiInstance.post(
    "master/item/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteMenuItem = async (data: { id: number }) => {
  const deleteRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/item/delete"
      : "item/delete";
  const res = await backendApiInstance.post(deleteRoute, data);
  return res.data;
};

export const getProductsByCategory = async (data: { categoryId: number }) => {
  const res = await backendApiInstance.get(
    `item/getWithParams?category=${data.categoryId}`
  );
  return res.data;
};
