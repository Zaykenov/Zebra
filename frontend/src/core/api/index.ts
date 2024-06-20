import axios from "axios";
import { deleteMenuItem } from "./menu-items";
import { deleteIngredient } from "./ingredient";
import { deleteProductCategory } from "./product-categories";
import {
  deleteIngredientCategory,
  getTechCartByIngredientID,
} from "./ingredient-category";
import { deleteSklad } from "./sklad";
import { deleteDish } from "./dishes";
import { deleteSupplier } from "./suppliers";
import { deleteAccount } from "./accounts";
import { deleteSupply } from "./supplies";
import { deleteWorker } from "./workers";
import { deleteTransaction } from "./transactions";
import { deleteInventory, deleteInventoryItem } from "./inventory";
import { deactivateCheck, deleteCheckById } from "./check";
import Router from "next/router";
import { deleteWaste } from "./wastes";
import { deleteTransfer } from "./transfers";
import { deleteMasterModifier, deleteModifier } from "@api/modifiers";
import { objectToQueryParamsString } from "@utils/objectToQueryParamsString";
import { isMaster } from "@utils/checkMasterRole";

export const BACKEND_URL = process.env.NEXT_PUBLIC_API_URL;
export const MOBILE_URL = process.env.NEXT_PUBLIC_MOBILE_URL;
export const EXCEL_URL = process.env.NEXT_PUBLIC_EXCEL_URL;

let token = "";
if (typeof localStorage !== "undefined")
  token = localStorage.getItem("zebra.authToken") || "";

export const backendApiInstance = axios.create({
  baseURL: BACKEND_URL,
  headers: {
    Authorization: `Bearer ${token}`,
  },
});

export const mobileApiInstance = axios.create({
  baseURL: MOBILE_URL,
});

export const excelGeneratorApiInstance = axios.create({
  baseURL: EXCEL_URL,
});

backendApiInstance.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err?.response?.status === 401 && window !== undefined) {
      localStorage.removeItem("zebra.authToken");
      return Router.replace("/login");
    } else return Promise.reject(err);
  }
);

export const getDeleteResourceDetails = async (id: number, path: string) => {
  switch (path) {
    case "/ingredients":
      return getTechCartByIngredientID(id);
    default:
      break;
  }
};

export const deleteResource = (id: any, path: string) => {
  switch (path) {
    case "/menu":
      return deleteMenuItem({ id });
    case "/dishes":
      return deleteDish({ id });
    case "/ingredients":
      return deleteIngredient({ id });
    case "/categories_products":
      return deleteProductCategory({ id });
    case "/categories_ingredients":
      return deleteIngredientCategory({ id });
    case "/storages":
      return deleteSklad({ id });
    case "/suppliers":
      return deleteSupplier({ id });
    case "/accounts":
      return deleteAccount({ id });
    case "/supplier":
      return deleteSupplier({ id });
    case "/supply":
      return deleteSupply({ id });
    case "/terminal/orders":
      return deleteCheckById({ id });
    case "/receipts":
      return deactivateCheck(id);
    case "/workers":
      return deleteWorker({ id });
    case "/transactions":
      return deleteTransaction({ id });
    case "/inventory":
      return deleteInventory({ id });
    case "/inventory/[id]":
      return deleteInventoryItem({ id });
    case "/terminal/inventory/[id]":
      return deleteInventoryItem({ id });
    case "/waste":
      return deleteWaste({ id });
    case "/transfer":
      return deleteTransfer({ id });
    case "/nabors":
      return isMaster() ? deleteMasterModifier({ id }) : deleteModifier({ id });
    default:
      return;
  }
};

export const getItems = async () => {
  const res = await backendApiInstance.get("item/getItems");
  return res.data;
};

export enum QueryOptions {
  FROM = "from",
  TO = "to",
  SORT = "sort",
  SEARCH = "search",
  PAGE = "page",
  CATEGORY = "category",
  WORKER = "worker_id",
  SKLAD = "sklad_id",
  SCHET = "schet_id",
  PAYMENT = "payment_type",
  DEALER = "dealer_id",
  TYPE = "type",
  STATUS = "status",
  SHOP = "shop",
}

export type QueryOptionsData = {
  [QueryOptions.FROM]?: string;
  [QueryOptions.TO]?: string;
  [QueryOptions.SORT]?: string;
  [QueryOptions.SEARCH]?: string;
  [QueryOptions.PAGE]?: number;
  [QueryOptions.CATEGORY]?: number;
  [QueryOptions.WORKER]?: number;
  [QueryOptions.SKLAD]?: number | number[];
  [QueryOptions.SCHET]?: number;
  [QueryOptions.PAYMENT]?: string;
  [QueryOptions.DEALER]?: number;
  [QueryOptions.TYPE]?: string;
  [QueryOptions.STATUS]?: string;
  [QueryOptions.SHOP]?: number | number[];
};

export const queryOptionsToString = (
  queryOptions?: QueryOptionsData,
  withDate: boolean = true
) => {
  if (!queryOptions) return "";
  const defaultStartDate = new Date();
  defaultStartDate.setMonth(new Date().getMonth() - 1);
  const defaultEndDate = new Date();
  const from =
    queryOptions[QueryOptions.FROM] || (withDate && defaultStartDate);
  const to = queryOptions[QueryOptions.TO] || (withDate && defaultEndDate);
  const search = queryOptions[QueryOptions.SEARCH];
  const sort = queryOptions[QueryOptions.SORT];
  const page = queryOptions[QueryOptions.PAGE];
  const category = queryOptions[QueryOptions.CATEGORY];
  const worker = queryOptions[QueryOptions.WORKER];
  const sklad = queryOptions[QueryOptions.SKLAD];
  const schet = queryOptions[QueryOptions.SCHET];
  const payment = queryOptions[QueryOptions.PAYMENT];
  const dealer = queryOptions[QueryOptions.DEALER];
  const type = queryOptions[QueryOptions.TYPE];
  const status = queryOptions[QueryOptions.STATUS];
  const shop = queryOptions[QueryOptions.SHOP];

  let queryObj = {
    ...(search ? { [QueryOptions.SEARCH]: search } : {}),
    ...(sort ? { [QueryOptions.SORT]: sort } : {}),
    ...(page ? { [QueryOptions.PAGE]: page } : {}),
    ...(from ? { [QueryOptions.FROM]: from } : {}),
    ...(to ? { [QueryOptions.TO]: to } : {}),
    ...(category ? { [QueryOptions.CATEGORY]: category } : {}),
    ...(worker ? { [QueryOptions.WORKER]: worker } : {}),
    ...(sklad ? { [QueryOptions.SKLAD]: sklad } : {}),
    ...(schet ? { [QueryOptions.SCHET]: schet !== -1 ? schet : 0 } : {}),
    ...(payment ? { [QueryOptions.PAYMENT]: payment } : {}),
    ...(dealer ? { [QueryOptions.DEALER]: dealer } : {}),
    ...(type ? { [QueryOptions.TYPE]: type } : {}),
    ...(status ? { [QueryOptions.STATUS]: status } : {}),
    ...(shop ? { [QueryOptions.SHOP]: shop } : {}),
  };
  queryObj = Object.fromEntries(
    Object.entries(queryObj).filter(([_, v]) => v !== null)
  );
  // @ts-ignore
  return objectToQueryParamsString(queryObj);
};
