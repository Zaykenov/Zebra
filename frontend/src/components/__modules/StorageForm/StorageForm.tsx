import React, { FC, useCallback, useEffect } from "react";
import { LabeledInput } from "@shared/ui/Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { StorageData } from "./types";
import { createSklad, updateSklad } from "@api/sklad";

export interface MenuItemFormProps {
  data?: any;
}

const StorageForm: FC<MenuItemFormProps> = ({ data }) => {
  const router = useRouter();

  const { handleSubmit, register, reset } = useForm<StorageData>({
    defaultValues: {
      name: data?.name || "",
      address: data?.address || "",
    },
  });

  useEffect(() => {
    reset(data);
  }, [data, reset]);

  const onSubmit = useCallback(
    (submitData: StorageData) => {
      if (!data)
        createSklad(submitData).then(() => router.replace("/storages"));
      else {
        updateSklad({
          id: data.id,
          ...submitData,
        }).then(() => router.replace("/storages"));
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
      <LabeledInput {...register("address")} label="Адрес" />
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

export default StorageForm;
