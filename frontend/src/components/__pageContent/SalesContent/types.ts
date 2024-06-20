export enum Period {
  DAY,
  WEEK,
  MONTH,
}

export enum Stats {
  REVENUE = "revenue",
  PROFIT = "profit",
  CHECKS = "checks",
  VISITORS = "visitors",
  AVG_CHECK = "avg_check",
}

export interface BarChartBoxProps {
  title: string;
  vertical?: boolean;
  switchable?: boolean;
  data: {
    name: string;
    value: number;
    value2?: number;
  }[];
  className?: string;
  isCurrency?: boolean;
  tooltipPayload?: {
    name: string;
    isCurrency: boolean;
  } | null;
}
