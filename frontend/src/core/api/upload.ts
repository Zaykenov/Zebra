import { backendApiInstance } from "./index";

export const uploadResource = async (formData: any) => {
  const res = await backendApiInstance.post("upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return res.data;
};
