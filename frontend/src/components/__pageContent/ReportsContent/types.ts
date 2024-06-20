import { Row } from "react-table";

export interface ReportItem {
  item_name: string;
  type: string;
  initial_ostatki: string;
  initial_sum: string;
  income: string;
  consumption: number;
  final_ostatki: number;
  final_netCost: number;
  final_sum: number;
}

export interface DetailsPopoverProps {
  details: {
    inventarization?: number;
    postavka?: number;
    postavka_cost?: number;
    sales?: number;
    measure: string;
    consumption?: number;
    transfer?: number;
  };
  value: number;
  row: Row;
}
