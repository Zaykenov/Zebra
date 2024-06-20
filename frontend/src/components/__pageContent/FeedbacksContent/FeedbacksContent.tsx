import React, { FC, useCallback, useEffect, useState } from "react";
import { getExcelFile } from "@api/excel";
import { getDateAndTime, getDateString } from "@utils/dateFormatter";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import extractHeaderAndAccessor from "@utils/extractHeaderAndAccessor";
import { getAllFeedbacks } from "@api/feedbacks";
import filterProperties from "@utils/filterProperties";
import { useFilter } from "@context/index";
import { FeedbackData } from "./types";
import { columns } from "./constants";
import SubRow from "./SubRow";

const FeedbacksContent: FC = () => {
  const { queryOptions, changeTotalResults, changeTotalPages } = useFilter();

  const [tableData, setTableData] = useState<FeedbackData[]>([]);

  const { headers, accessors } = extractHeaderAndAccessor(columns);
  const [excelData, setExcelData] = useState<Partial<FeedbackData>[]>([]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  const handleGetAll = useCallback((res: any) => {
    const transformItem = (feedback: any) => ({
      id: feedback.id,
      userId: feedback.user_id,
      username: feedback.username,
      checkId: feedback.check_id,
      shopId: feedback.shop_id,
      shopName: feedback.shop_name,
      workerId: feedback.worker_id,
      workerName: feedback.worker_name,
      scoreQuality: feedback.score_quality,
      scoreService: feedback.score_service,
      feedbackText: feedback.feedback_text,
      feedbackDate: getDateAndTime(feedback.feedback_date),
      check: JSON.parse(feedback.check_json),
    });

    let data;
    if (res.data.data) {
      data = res.data.data.map(transformItem);
    } else {
      data = res.data.map(transformItem);
    }

    return data;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length) return;
    getAllFeedbacks(queryOptions).then((res) => {
      const data = handleGetAll(res);
      setTableData(data);
      changeTotalPages(res.data.totalPages);
      changeTotalResults(data.length);
    });
    const currentQueryOptions = { ...queryOptions };
    delete currentQueryOptions.page;
    getAllFeedbacks(currentQueryOptions).then((res) => {
      const data = handleGetAll(res);
      const filteredProperties = filterProperties<FeedbackData>(
        data,
        accessors
      );
      setExcelData(filteredProperties);
    });
  }, [queryOptions]);

  return (
    <MainLayout
      title="Отзывы клиентов"
      excelDownloadButton={() =>
        getExcelFile(`Отзывы клиентов ${getDateString()}`, headers, excelData)
      }
      pagination
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          editable={false}
          details
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default FeedbacksContent;
