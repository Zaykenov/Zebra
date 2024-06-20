import { backendApiInstance } from "@api/index";

export const getMasterShops = async () => {
  const res = await backendApiInstance.get("master/getShops");
  return res.data;
};
