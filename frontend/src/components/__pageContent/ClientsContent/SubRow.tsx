import React, { FC, useEffect, useState } from "react";
import { Row } from "react-table";
import clientsService from "@services/clientsService";

const SubRow: FC<{ rowData: Row<any> }> = ({ rowData }) => {
  const [data, setData] = useState<any>(null);

  const getProductNames = (check: any) => {
    const techCarts = JSON.parse(check).TechCarts.map((item: any) => item.Name);
    const tovars = JSON.parse(check).Tovars.map((item: any) => item.Name);
    return techCarts.concat(tovars).join(" + ");
  };
  useEffect(() => {
    clientsService.getClientById(rowData.original.id).then((res) => {
      setData(res.data);
    });
  }, []);

  return (
    <div className="p-4 px-10 bg-white">
      <div className="w-full flex flex-col pr-4">
        <div className="text-center mb-5 font-bold">Отзывы</div>
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          <div className="w-1/6 text-left">ID Склада</div>
          <div className="w-1/6 text-left">Оценка качества</div>
          <div className="w-1/6 text-left">Оценка сервиса</div>
          <div className="w-1/6 text-left">Отзыв</div>
          <div className="w-1/6 text-left">Товары</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {data &&
            data.feedbacks.map((feedback: any) => (
              <div key={`tech-${feedback.id}`} className="flex text-sm">
                <div className="w-1/6 text-left">{feedback.shop_id}</div>
                <div className="w-1/6 text-left">{feedback.score_quality}</div>
                <div className="w-1/6 text-left">{feedback.score_service}</div>
                <div className="w-1/6 text-left">{feedback.feedback_text}</div>
                <div className="w-1/6 text-left">
                  {getProductNames(feedback.check_json)}
                </div>
              </div>
            ))}
        </div>
        <div className="text-center mb-5 font-bold">Чеки</div>
        <div className="flex text-xs text-gray-500 pb-2 border-b border-gray-300">
          {/* <div className="grow">Отзывы</div> */}
          <div className="w-1/3 text-left">Наличными</div>
          <div className="w-1/3 text-left">Картой</div>
          <div className="w-1/3 text-left">Скидка</div>
          <div className="w-1/3 text-left">Сумма</div>
        </div>
        <div className="flex flex-col space-y-3 mb-4 pt-2 font-medium">
          {data &&
            data.checks.map((check: any) => (
              <div key={`tech-${check.id}`} className="flex text-sm">
                <div className="w-1/3 text-left">{check.cash}</div>
                <div className="w-1/3 text-left">{check.card}</div>
                <div className="w-1/3 text-left">{check.discount}</div>
                <div className="w-1/3 text-left">{check.sum}</div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SubRow;
