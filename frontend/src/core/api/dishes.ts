import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";
import { DishData } from "@modules/DishForm/types";

const processFormData = (data: DishData) => {
  const postData = {
    name: data.name,
    category: data.category,
    image: data.image,
    tax: "Фискальный налог",
    measure: data.measure,
    discount: data.discount,
    price: data.price,
    ingredient_tech_cart: data.ingredient_tech_cart,
    nabor: data.nabor,
    shop_id: data.shop_id,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllDishes = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "item/techCart/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMasterDishes = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "master/techCart/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getDish = async (id: string) => {
  const getRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/techCart/get/"
      : "item/techCart/get?id=";
  const res = await backendApiInstance.get(`${getRoute}${id}`);
  return res.data;
};

export const getDishesByCategory = async (data: { categoryId: number }) => {
  const res = await backendApiInstance.get(
    `item/techCart/getWithParams?category=${data.categoryId}`
  );
  return res.data;
};

export const createDish = async (data: DishData) => {
  const res = await backendApiInstance.post(
    "item/techCart/create",
    processFormData(data)
  );
  return res.data;
};

export const createMasterDish = async (data: DishData) => {
  const res = await backendApiInstance.post(
    "master/item/techCart/create",
    processFormData(data)
  );
  return res.data;
};

export const updateDish = async (data: DishData) => {
  const res = await backendApiInstance.post(
    "item/techCart/update",
    processFormData(data)
  );
  return res.data;
};

export const updateMasterDish = async (data: DishData) => {
  const res = await backendApiInstance.post(
    "master/techCart/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteDish = async (data: { id: number }) => {
  const deleteRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/techCart/delete"
      : "item/techCart/delete";
  const res = await backendApiInstance.post(deleteRoute, data);
  return res.data;
};
