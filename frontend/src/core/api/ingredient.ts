import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";
import { IngredientData } from "@modules/IngredientForm/types";

const processFormData = (data: IngredientData) => {
  const postData = {
    name: data.name,
    category: data.category,
    measure: data.measure,
    cost: data.cost,
    shop_id: data.shop_id,
  };
  return data.id ? { id: data.id, ...postData } : postData;
};

export const getAllIngredients = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "ingredient/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMasterIngredients = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "master/ingredient/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getIngredient = async (id: string) => {
  const getRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/ingredient/get/"
      : "ingredient/get?id=";
  const res = await backendApiInstance.get(`${getRoute}${id}`);
  return res.data;
};

export const createIngredient = async (data: any) => {
  const res = await backendApiInstance.post(
    "ingredient/create",
    processFormData(data)
  );
  return res.data;
};

export const updateIngredient = async (data: any) => {
  const res = await backendApiInstance.post(
    "ingredient/update",
    processFormData(data)
  );
  return res.data;
};

export const updateMasterIngredient = async (data: any) => {
  const res = await backendApiInstance.post(
    "master/ingredient/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteIngredient = async (data: { id: number }) => {
  const deleteRoute =
    typeof window !== "undefined" &&
    localStorage.getItem("zebra.role") === "master"
      ? "master/ingredient/delete"
      : "ingredient/delete";
  const res = await backendApiInstance.post(deleteRoute, data);
  return res.data;
};
