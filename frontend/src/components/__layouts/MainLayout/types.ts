import { ReactNode } from "react";

export enum FilterOption {
  DEALER = "dealer",
  ITEMS_CATEGORY = "itemsCategory",
  PRODUCTS_CATEGORY = "productsCategory",
  INGREDIENTS_CATEGORY = "ingredientsCategory",
  ITEM_TYPE = "itemType",
  SCHET = "schet",
  SHOP = "shop",
  SKLAD = "sklad",
  WORKER = "worker",
}

export interface MainLayoutProps {
  title: string;
  addHref?: string;
  backBtn?: boolean;
  dateFilter?: boolean;
  searchFilter?: boolean;
  pagination?: boolean;
  children?: ReactNode;
  customBtns?: ReactNode[];
  content?: ReactNode;
  withoutOverflow?: boolean;
  excelDownloadButton?: any;
  filterOptions?: FilterOption[];
}
