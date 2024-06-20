import React, { FC, useCallback, useEffect } from "react";
import { LabeledInput } from "@shared/ui/Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { SupplierData } from "./types";
import { createSupplier, updateSupplier } from "@api/suppliers";

export interface SupplierFormProps {
  data?: any;
}

const SupplierForm: FC<SupplierFormProps> = ({ data }) => {
  const router = useRouter();

  const { handleSubmit, register, reset } = useForm<SupplierData>({
    defaultValues: {
      name: data?.name || "",
      phone: data?.phone || "",
      address: data?.address || "",
      comment: data?.comment || "",
    },
  });

  useEffect(() => {
    reset(data);
  }, [data, reset]);

  const onSubmit = useCallback(
    (submitData: SupplierData) => {
      if (!data)
        createSupplier(submitData).then(() => router.replace("/suppliers"));
      else {
        updateSupplier({
          id: data.id,
          ...submitData,
        }).then(() => router.replace("/suppliers"));
      }
    },
    [data, router]
  );

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-1/2 space-y-5"
    >
      <LabeledInput {...register("name")} label="Название" />
      <LabeledInput {...register("phone")} label="Телефон" />
      <LabeledInput {...register("address")} label="Адрес" />
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

export default SupplierForm;
