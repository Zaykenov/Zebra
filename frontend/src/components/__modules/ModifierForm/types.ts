import { MeasureOption } from "@shared/types/types.";

export type ModifierData = {
  name: string;
  min?: number;
  max?: number;
  ingredient_nabor: {
    value: number;
    ingredient_id: number;
    brutto: number;
    price: number;
    measure: string;
  }[];
  id?: number;
};
