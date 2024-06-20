import { backendApiInstance } from "./index";
import { WorkerRole } from "./workers";

export type ShopData = {
  name: string;
  address: string;
  tis_token: string;
  cash_schet: number;
  card_schet: number;
  limit: number;
};

export type SkladData = {
  name: string;
  address: string;
};

export type WorkerData = {
  name: string;
  phone: string;
  username: string;
  password: string;
  role: WorkerRole;
  new: boolean;
};

export type ShopFullData = {
  shop: ShopData;
  sklad: SkladData;
  workers: WorkerData[];
  tovars: number[];
  products_shop: {
    tovars: number[];
    tech_carts: number[];
  };
};

export const getAllShops = async () => {
  const res = await backendApiInstance.get("shop/getAll");
  return res.data;
};

export const getShop = async (id: string) => {
  const res = await backendApiInstance.get(`shop/get?id=${id}`);
  return res.data;
};

export const createShop = async (data: ShopFullData) => {
  const res = await backendApiInstance.post("shop/create", data);
  return res.data;
};
