import {
  DealerSelect,
  ItemsCategorySelect,
  ItemTypeSelect,
  SchetSelect,
  ShopSelect,
  SkladSelect,
  WorkerSelect,
} from "@common/FilterSelect";
import { FilterOption } from "@layouts/MainLayout/types";
import { ReactNode } from "react";

export const filterComponents: { [key in FilterOption]: ReactNode } = {
  dealer: <DealerSelect />,
  itemsCategory: <ItemsCategorySelect />,
  productsCategory: <ItemsCategorySelect type="product" />,
  ingredientsCategory: <ItemsCategorySelect type="ingredient" />,
  itemType: <ItemTypeSelect />,
  schet: <SchetSelect />,
  shop: <ShopSelect />,
  sklad: <SkladSelect />,
  worker: <WorkerSelect />,
};
