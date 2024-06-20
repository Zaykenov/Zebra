import { backendApiInstance } from ".";

export type NaborIngredient = {
  id: number;
  name: string;
  price: number;
};

export type Modificator = {
  id: number;
  name: string;
  price: number;
  nabor_ingredient: NaborIngredient[];
};

export type Product = {
  id: number;
  name: string;
  categoryName: string;
  type: "tovar" | "techCart";
  image?: string;
  price: number;
  discount: boolean;
  nabor?: Modificator[];
};

export type Category = {
  id: number;
  category: string;
  image: string;
  products: Product[];
};

interface TerminalStartData {
    categories: Category[],
    mainDisplay: Product[],
    currentShift: any
}

export const getTerminalStartData = async (): Promise<TerminalStartData> => {
  const response = await backendApiInstance.get("/terminal/start");
  const { categories, mainDisplay, current_shift } = response.data.data;
  return {
    categories,
    mainDisplay,
    currentShift: current_shift
  }
};
