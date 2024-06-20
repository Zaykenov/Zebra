export type TotalPaymentData = {
  total_card: number;
  total_cash: number;
  total_check_count: number;
  total_total: number;
};

export interface PaymentItem {
  time: string;
  check_count: string;
  intCash: string;
  intCard: string;
  intTotal: string;
  comment: string;
  status: string;
}
