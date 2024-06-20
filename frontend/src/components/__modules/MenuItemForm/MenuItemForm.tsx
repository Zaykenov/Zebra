import React, { FC, useCallback, useEffect, useState } from "react";
import { Input, LabeledInput } from "@shared/ui/Input";
import LabeledSelect from "@shared/ui/Select/LabeledSelect";
import { Bars2Icon, PlusIcon } from "@heroicons/react/24/outline";
import { useForm, Controller } from "react-hook-form";
import {
  createMenuItem,
  updateMasterMenuItem,
  updateMenuItem,
} from "@api/menu-items";
import { useRouter } from "next/router";
import { getAllProductCategories } from "@api/product-categories";
import { MenuItemData } from "./types";
import ImageUpload from "@common/ImageUpload/ImageUpload";
import { MeasureOption, measureOptions } from "@shared/types/types.";
import { formatInputValue } from "@utils/formatInputValue";
import Checkbox from "@shared/ui/Checkbox/Checkbox";
import AlertMessage from "@common/AlertMessage";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import { Select } from "antd";
import { getAllShops } from "@api/shops";

export interface MenuItemFormProps {
  data?: any;
}

const MenuItemForm: FC<MenuItemFormProps> = ({ data }) => {
  const router = useRouter();
  const [categoryOptions, setCategoryOptions] = useState([]);
  const [shopOptions, setShopOptions] = useState<
    {
      value: number;
      label: string;
    }[]
  >([]);

  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();
  useEffect(() => {
    getAllProductCategories().then((res) => {
      const categories = res.data.map(
        ({ id, name }: { id: string; name: string }) => ({
          label: name,
          value: parseInt(id),
        })
      );
      setCategoryOptions(categories);
    });
    getAllShops().then((res) => {
      const shops = res.data.map(
        ({ id, name }: { id: number; name: string }) => ({
          label: name,
          value: id,
        })
      );
      setShopOptions(shops);
    });
  }, []);

  const [uploadedImage, setUploadedImage] = useState<string | null>(null);
  useEffect(() => {
    data && setUploadedImage(data.image);
  }, [data]);

  const [selectedShops, setSelectedShops] = useState<number[]>([]);

  const { handleSubmit, register, watch, setValue, reset, control } =
    useForm<MenuItemData>({
      defaultValues: {
        name: data?.name || "",
        category: data?.category || "",
        image: data?.image || "",
        price: data?.price || "",
        measure: data?.measure || MeasureOption.DEFAULT,
        discount: false,
        cost: data?.cost || "",
        margin: data?.margin || 0,
      },
    });

  useEffect(() => {
    if (!data || shopOptions.length === 0) return;
    reset({ ...data, category: data.category_id });
    setSelectedShops(
      typeof data.shop_id === "number" ? [data.shop_id] : data.shop_id || []
    );
  }, [data, reset, shopOptions]);

  const [role, setRole] = useState<string>("");

  useEffect(() => {
    const role = localStorage.getItem("zebra.role");
    role && setRole(role);
  }, []);

  const onSubmit = useCallback(
    (submitData: MenuItemData) => {
      //@ts-ignore
      submitData.price = parseFloat(submitData.price);
      //@ts-ignore
      submitData.cost = parseFloat(submitData.cost);
      //@ts-ignore
      submitData.margin = parseFloat(submitData.margin);

      if (role !== "master") {
        submitData.shop_id =
          selectedShops.length > 0
            ? selectedShops
            : shopOptions.map((shop) => shop.value);
      }

      if (
        submitData.category &&
        (submitData.cost as unknown as number) > 0 &&
        submitData.name &&
        submitData.price
      ) {
        if (!data) {
          createMenuItem({ ...submitData, image: uploadedImage || "" }).then(
            () => router.replace("/menu")
          );
        } else {
          role === "master"
            ? updateMasterMenuItem({
                id: data?.id,
                ...submitData,
              }).then(() => router.replace("/menu"))
            : updateMenuItem({
                id: data.id,
                ...submitData,
                image: uploadedImage || "",
              }).then(() => router.replace("/menu"));
        }
      } else {
        showAlertMessage(
          "Некоторые поля не заполнены",
          AlertMessageType.WARNING
        );
      }
    },
    [data, uploadedImage, router, selectedShops, role]
  );

  const [cost] = useState<number>(0);
  const [price] = useState<number>(0);
  const [margin, setMargin] = useState<number>(0);

  useEffect(() => {
    setValue("cost", cost.toFixed(1));
  }, [cost, setValue]);

  useEffect(() => {
    setValue("price", price.toFixed(1));
  }, [price, setValue]);

  useEffect(() => {
    //@ts-ignore
    setValue("margin", margin.toFixed(2));
  }, [margin, setValue]);

  useEffect(() => {
    const subscriptionPrice = watch((value, { name }) => {
      let costValue =
        typeof value.cost === "string"
          ? formatInputValue(value.cost || "").numberValue
          : value.cost || 0;
      let priceValue =
        typeof value.price === "string"
          ? formatInputValue(value.price || "").numberValue
          : value.price || 0;
      if (name == "cost") {
        const marginValue = (100 * (priceValue - costValue)) / costValue;
        if (isNaN(marginValue) || !isFinite(marginValue)) {
          setMargin(0);
          return;
        }
        setMargin(Math.round((marginValue + Number.EPSILON) * 100) / 100);
      } else {
        const marginValue = (100 * (priceValue - costValue)) / costValue;
        if (isNaN(marginValue) || !isFinite(marginValue)) {
          setMargin(0);
          return;
        }
        setMargin(Math.round((marginValue + Number.EPSILON) * 100) / 100);
      }
    });
    return () => {
      subscriptionPrice.unsubscribe();
    };
  }, [watch]);

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-1/2 space-y-5"
    >
      {alertMessage && (
        <AlertMessage
          message={alertMessage.message}
          type={alertMessage.type}
          onClose={hideAlertMessage}
        />
      )}
      <LabeledInput {...register("name")} label="Название" />
      <div className="w-full flex items-center pt-2">
        <label className="w-40 mr-4">Категория</label>
        <Controller
          control={control}
          {...register("category")}
          render={({ field: { onChange, value } }) => (
            <Select
              allowClear
              style={{ width: "100%", flex: 1 }}
              placeholder="Все заведения"
              options={categoryOptions}
              value={value}
              onChange={onChange}
              filterOption={(input, option) =>
                (option?.label ?? "")
                  .trim()
                  .toLowerCase()
                  .includes(input.trim().toLowerCase())
              }
            />
          )}
        />
      </div>
      <div className="w-full flex pt-2">
        <label className="w-40 mr-4">Обложка</label>
        <ImageUpload
          uploadedImage={uploadedImage}
          setUploadedImage={setUploadedImage}
        />
      </div>
      <LabeledSelect
        {...register("measure")}
        fieldClass="w-1/2"
        label="Ед. измерения"
        options={measureOptions}
      />
      <div className="w-full flex items-center pt-2">
        <label className="w-40 mr-4">Участвует в скидках</label>
        <Checkbox {...register("discount")} />
      </div>
      {role !== "master" && (
        <div className="w-full flex items-center pt-2">
          <label className="w-40 mr-4">Заведения</label>
          <Select
            mode="multiple"
            allowClear
            style={{ width: "100%", flex: 1 }}
            placeholder="Все заведения"
            value={selectedShops}
            onChange={(value) => {
              setSelectedShops(value);
            }}
            options={shopOptions}
            filterOption={(input, option) =>
              (option?.label ?? "")
                .trim()
                .toLowerCase()
                .includes(input.trim().toLowerCase())
            }
          />
        </div>
      )}
      <div className="w-full flex">
        <span className="w-40 mr-4">Цена</span>
        <div className="flex items-center space-x-4">
          <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
            <label className="text-xs text-gray-600">Себестоимость</label>
            <div className="flex items-center space-x-2">
              <Input {...register("cost")} type="text" className="text-right" />{" "}
              <span className="text-lg">₸</span>
            </div>
          </div>
          <PlusIcon className="w-6 h-6 text-gray-400 mt-4" />
          <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
            <label className="text-xs text-gray-600">Наценка</label>
            <div className="flex items-center space-x-2">
              <Input
                value={margin}
                type="text"
                className="text-right"
                disabled
              />{" "}
              <span className="text-lg">%</span>
            </div>
          </div>
          <Bars2Icon className="w-6 h-6 text-gray-400 mt-4" />
          <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
            <label className="text-xs text-gray-600">Итого</label>
            <div className="flex items-center space-x-2">
              <Input
                {...register("price")}
                type="text"
                className="text-right"
              />{" "}
              <span className="text-lg">₸</span>
            </div>
          </div>
        </div>
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

export default MenuItemForm;
