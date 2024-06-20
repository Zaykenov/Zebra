import { backendApiInstance } from "./index";

export type TagData = {
  id?: number;
  text: string;
};

const processFormData = (data: TagData) => {
  const postData = {
    text: data.text || "",
  };
  return postData;
};

export const getAllTags = async () => {
  const res = await backendApiInstance.get("check/tag/getAll");
  return res.data;
};

export const getTag = async (id: string) => {
  const res = await backendApiInstance.get(`check/tag/get/${id}`);
  return res.data;
};

export const createTag = async (data: TagData) => {
  const res = await backendApiInstance.post(
    "check/tag/create",
    processFormData(data)
  );
  return res.data;
};

export const updateTag = async (data: TagData) => {
  const res = await backendApiInstance.post(
    "check/tag/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteTag = async (data: { id: number }) => {
  const res = await backendApiInstance.post("check/tag/delete", data);
  return res.data;
};
