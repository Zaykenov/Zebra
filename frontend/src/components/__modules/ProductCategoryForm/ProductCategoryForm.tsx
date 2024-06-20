import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledInput } from "@shared/ui/Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import {
  createProductCategory,
  updateProductCategory,
} from "@api/product-categories";
import { ProductCategoryData } from "./types";
import ImageUpload from "@common/ImageUpload/ImageUpload";

export interface ProductCategoryFormProps {
  data?: any;
}

const ProductCategoryForm: FC<ProductCategoryFormProps> = ({ data }) => {
  const router = useRouter();

  const [uploadedImage, setUploadedImage] = useState<string | null>(null);
  useEffect(() => {
    data && setUploadedImage(data.image);
  }, [data]);

  const { handleSubmit, register, reset } = useForm<ProductCategoryData>({
    defaultValues: {
      name: data?.name || "",
      image: data?.image || "",
    },
  });

  useEffect(() => {
    reset(data);
  }, [data, reset]);

  const onSubmit = useCallback(
    (submitData: any) => {
      if (!data)
        createProductCategory({
          ...submitData,
          image: uploadedImage || "",
        }).then(() => router.replace("/categories_products"));
      else {
        updateProductCategory({
          id: data.id,
          ...submitData,
          image: uploadedImage || "",
        }).then(() => router.replace("/categories_products"));
      }
    },
    [data, uploadedImage, router]
  );

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-1/2 space-y-5"
    >
      <LabeledInput {...register("name")} label="Название" />
      <div className="w-full flex pt-2">
        <label className="w-40 mr-4">Обложка</label>
        <ImageUpload
          uploadedImage={uploadedImage}
          setUploadedImage={setUploadedImage}
        />
      </div>
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

export default ProductCategoryForm;
