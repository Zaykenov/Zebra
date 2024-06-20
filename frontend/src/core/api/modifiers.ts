import {
  backendApiInstance,
  QueryOptionsData,
  queryOptionsToString,
} from "./index";
import { ModifierData } from "../../components/__modules/ModifierForm/types";

const processFormData = (data: ModifierData) => {
  const postData = {
    name: data.name,
    min: data.min || 0,
    max: data.max || 0,
    ingredient_nabor: data.ingredient_nabor,
  };
  return data.id
    ? {
        id: data.id,
        ...postData,
      }
    : postData;
};

export const getAllModifiers = async (queryOptions?: QueryOptionsData) => {
  const res = await backendApiInstance.get(
    "ingredient/nabor/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMasterModifierNabors = async (
  queryOptions?: QueryOptionsData
) => {
  const res = await backendApiInstance.get(
    "master/ingredient/nabor/getAll" + queryOptionsToString(queryOptions, false)
  );
  return res.data;
};

export const getMasterModifierNabor = async (id: number) => {
  const res = await backendApiInstance.get(`master/ingredient/nabor/get/${id}`);
  return res.data;
};

export const getModifier = async (id: number) => {
  const res = await backendApiInstance.get(`ingredient/nabor/get?id=${id}`);
  return res.data;
};

export const getTechCartModifiers = async (id: string) => {
  const res = await backendApiInstance.get(`item/techCart/getNabor?id=${id}`);
  return res.data;
};

export const createModifier = async (data: ModifierData) => {
  const res = await backendApiInstance.post("ingredient/nabor/create", data);
  return res.data;
};

export const createAllModifiers = async (data: ModifierData[]) => {
  const res = await backendApiInstance.post(
    "ingredient/nabor/createAll",
    data.map((modifierData) => processFormData(modifierData))
  );
  return res.data;
};

export const createAllMasterModifiers = async (data: ModifierData[]) => {
  const res = await backendApiInstance.post(
    "master/ingredient/nabor/createAll",
    data.map((modifierData) => processFormData(modifierData))
  );
  return res.data;
};

export const updateModifier = async (data: ModifierData) => {
  const res = await backendApiInstance.post("ingredient/nabor/update", data);
  return res.data;
};

export const createMasterModifier = async (data: ModifierData) => {
  const res = await backendApiInstance.post(
    "master/ingredient/nabor/add",
    data
  );
  return res.data;
};

export const updateMasterModifier = async (data: ModifierData) => {
  const res = await backendApiInstance.post(
    "master/ingredient/nabor/update",
    data
  );
  return res.data;
};

export const deleteModifier = async (data: { id: number }) => {
  const res = await backendApiInstance.post("ingredient/nabor/delete", data);
  return res.data;
};

export const deleteMasterModifier = async (data: { id: number }) => {
  const res = await backendApiInstance.post(
    "master/ingredient/nabor/delete",
    data
  );
  return res.data;
};
