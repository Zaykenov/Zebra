import { MeasureOption } from "@shared/types/types.";

export type DishData = {
  id?: number;
  name: string;
  category: number;
  image: string;
  tax: string;
  measure: MeasureOption;
  discount: boolean;
  cost: number;
  price: number;
  margin: number;
  shop_id: number[] | number;
  ingredient_tech_cart: {
    cost: number;
    ingredient_id: number;
    brutto: string;
  }[];
  nabor: {
    nabor_id: number;
  }[];
};
