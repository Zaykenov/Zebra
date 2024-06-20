export type IngredientData = {
  id?: number;
  name: string;
  category: number;
  measure: string;
  count: number;
  cost: number;
  sklad: string;
  shop_id: number[] | number;
};
