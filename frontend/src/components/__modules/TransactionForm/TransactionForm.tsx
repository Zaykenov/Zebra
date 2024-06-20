import React, { FC, useCallback, useEffect, useMemo, useState } from "react";
import { LabeledInput } from "@shared/ui/Input";
import { Controller, useForm } from "react-hook-form";
import { useRouter } from "next/router";
import {
  createTransaction,
  mapTransactionCategoryToString,
  TransactionCategory,
  TransactionData,
  updateTransaction,
} from "@api/transactions";
import { LabeledSelect } from "@shared/ui/Select";
import { getAllAccounts } from "@api/accounts";
import DatePicker from "react-datepicker";

export interface TransactionFormProps {
  data?: any;
  routerName?: string;
}

const excludedCategories: TransactionCategory[] = [
  TransactionCategory.POSTAVKA,
  TransactionCategory.OPEN_SHIFT,
  TransactionCategory.CLOSE_SHIFT,
];

const categoryOptions = Object.values(TransactionCategory)
  .filter((category) => !excludedCategories.includes(category))
  .map((category) => {
    return {
      name: mapTransactionCategoryToString(category as TransactionCategory),
      value: category,
    };
  });

const TransactionForm: FC<TransactionFormProps> = ({ data }) => {
  const router = useRouter();
  const [schetOptions, setSchetOptions] = useState<any[]>([]);

  const shiftId = useMemo(() => router.query.shift, [router]);

  const { handleSubmit, register, reset, control, getValues } =
    useForm<TransactionData>({
      defaultValues: {
        sum: data?.sum || 0,
        schet_id: data?.schet_id || 1,
        category: data?.category || "",
        time: data?.time || "",
        date: data ? new Date(data.time) : new Date(),
        comment: data?.comment || "",
      },
    });

  useEffect(() => {
    if (schetOptions.length === 0 || !data) return;
    reset({ ...data, date: data ? new Date(data.time) : new Date() });
  }, [data, schetOptions, reset]);

  useEffect(() => {
    getAllAccounts().then((res) => {
      setSchetOptions(
        res.data.map((item: any) => ({
          name: item.name,
          value: parseInt(item.id),
        }))
      );
      reset({
        ...getValues(),
        schet_id: res.data[0].id,
      });
    });
  }, []);

  const onSubmit = useCallback(
    (submitData: TransactionData) => {
      if (!data)
        createTransaction({
          sum: submitData.sum,
          schet_id: submitData.schet_id,
          category: submitData.category,
          time: submitData.date?.toISOString(),
          comment: submitData.comment,
          status: "negative",
          ...(shiftId !== undefined
            ? { shift_id: parseInt(shiftId as string) }
            : {}),
        }).then(() =>
          shiftId !== undefined
            ? router.replace("/shifts")
            : router.replace("/transactions")
        );
      else {
        updateTransaction({
          id: data.id,
          sum: submitData.sum,
          schet_id: submitData.schet_id,
          category: data.category,
          time: submitData.date?.toISOString(),
          comment: submitData.comment,
          status: data.status,
        }).then(() =>
          shiftId !== undefined
            ? router.replace("/shifts")
            : router.replace("/transactions")
        );
      }
    },
    [data, router]
  );

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-1/2 space-y-5"
    >
      <LabeledInput
        {...register("sum", { valueAsNumber: true })}
        label="Сумма"
      />
      {(!shiftId ||
        (data && data.category == TransactionCategory.POSTAVKA)) && (
        <LabeledSelect
          options={schetOptions}
          {...register("schet_id", { valueAsNumber: true })}
          label="Со счета"
        />
      )}

      {!data && (
        <LabeledSelect
          options={
            shiftId && !data
              ? [
                  ...categoryOptions,
                  {
                    name: "Закрытие смены",
                    value: TransactionCategory.CLOSE_SHIFT,
                  },
                ]
              : categoryOptions
          }
          {...register("category")}
          label="Категория"
        />
      )}
      <div className="w-full flex items-center">
        <label htmlFor="date" className="w-40 mr-4">
          Дата
        </label>
        <div>
          <Controller
            name="date"
            control={control}
            render={({ field }) => (
              <DatePicker
                locale="ru"
                renderCustomHeader={({ monthDate }) => (
                  <span className="font-medium font-inter text-sm capitalize">
                    {monthDate.toLocaleDateString("default", { month: "long" })}{" "}
                    {monthDate.getFullYear()}
                  </span>
                )}
                selected={field.value}
                onChange={field.onChange}
                timeInputLabel="Время:"
                dateFormat="dd.MM.yyyy HH:mm"
                showTimeInput
                className="rounded text-gray-800 py-2 px-3 border border-gray-300 focus:outline-none focus:border-indigo-500"
              />
            )}
          />
        </div>
      </div>
      <LabeledInput {...register("comment")} label="Комментарий" />

      <div className="pt-5 border-t border-gray-200">
        <button
          type="submit"
          className="py-2 px-3 bg-primary hover:bg-teal-600 transition duration-300 text-white rounded-md"
        >
          Сохранить
        </button>
      </div>
    </form>
  );
};

export default TransactionForm;
