import { backendApiInstance } from "./index";

export const getAllIngredientCategories = async () => {
  const res = await backendApiInstance.get("ingredient/category/getAll");
  return res.data;
};

export const getTechCartByIngredientID = async (id: number) => {
  const res = (await backendApiInstance.post(
    "ingredient/getTechCartByIngredientID",
    {id}
  )).data;
  return {
    name: res.data[0].name,
    link: `/dishes/${res.data[0].id}`,
    message: 'Если удалить этот ингредиент, изменится состав полуфабрикатов и тех.карт, в которые он входит:'
  };
};

export const getIngredientCategory = async (id: string) => {
  const res = await backendApiInstance.get(`ingredient/category/get?id=${id}`);
  return res.data;
};

export const createIngredientCategory = async (data: any) => {
  const res = await backendApiInstance.post("ingredient/category/create", data);
  return res.data;
};

export const updateIngredientCategory = async (data: any) => {
  const res = await backendApiInstance.post("ingredient/category/update", data);
  return res.data;
};

export const deleteIngredientCategory = async (data: any) => {
  const res = await backendApiInstance.post("ingredient/category/delete", data);
};
