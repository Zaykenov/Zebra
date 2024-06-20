import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import { getExcelFile } from "@api/excel";
import { getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { formatNumber } from "@utils/formatNumber";
import extractHeadersAndAccessor from "@utils/extractHeaderAndAccessor";
import { dateToString } from "@api/check";
import { getAllPayments } from "@api/payments";
import filterProperties from "@utils/filterProperties";
import { TotalPaymentData } from "./types";
import { columns } from "./constants";
import { useFilter } from "@context/index";

const PaymentsContent: FC = () => {
  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const [totalData, setTotalData] = useState<TotalPaymentData>({
    total_card: 0,
    total_cash: 0,
    total_check_count: 0,
    total_total: 0,
  });

  const [excelData, setExcelData] = useState<Partial<PaymentItem>[]>([]);

  const { headers, accessors } = extractHeadersAndAccessor(columns(totalData));

  const handleGetAll = useCallback((res: any) => {
    const transformPayment = (payment: any) => ({
      check_count: payment.check_count,
      cash: formatNumber(payment.cash, true, true),
      intCash: payment.cash,
      card: formatNumber(payment.card, true, true),
      intCard: payment.card,
      time: dateToString(payment.time, false, true),
      total: formatNumber(payment.total, true, true),
      intTotal: payment.total,
    });

    const extractTotals = (data: any) => ({
      total_check_count: data.total_check_count,
      total_cash: data.total_cash,
      total_card: data.total_card,
      total_total: data.total_total,
    });

    let tableData;
    let totals;
    if (res.data.data) {
      tableData = res.data.data.payments.map(transformPayment);
      totals = extractTotals(res.data.data);
    } else {
      tableData = res.data.payments.map(transformPayment);
      totals = extractTotals(res.data);
    }
    return { tableData, totals };
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;

    getAllPayments(queryOptions).then((res) => {
      const { tableData, totals } = handleGetAll(res);
      setTableData(tableData);
      setTotalData(totals);
      changeTotalResults(tableData.length);
      changeTotalPages(res.data.totalPages);
    });

    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllPayments(currentQueryOptions).then((res) => {
      const { tableData } = handleGetAll(res);
      const filteredProperties = filterProperties<PaymentItem>(
        tableData,
        accessors,
      );
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Оплаты"
      dateFilter
      searchFilter
      pagination
      filterOptions={[FilterOption.SHOP]}
      excelDownloadButton={() =>
        getExcelFile(`Оплаты ${getDateString()}`, headers, excelData)
      }
    >
      {tableData && (
        <Table
          columns={columns(totalData)}
          data={tableData}
          editable={false}
          hasFooter
        />
      )}
    </MainLayout>
  );
};

export default PaymentsContent;
