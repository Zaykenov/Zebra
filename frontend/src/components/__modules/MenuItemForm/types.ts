import { MeasureOption } from "@shared/types/types.";

export type MenuItemData = {
  id?: number;
  name: string;
  category: string;
  category_id: number;
  image: string;
  measure: MeasureOption;
  discount: boolean;
  cost: string;
  margin: number;
  price: string;
  shop_id: number[] | number;
};
