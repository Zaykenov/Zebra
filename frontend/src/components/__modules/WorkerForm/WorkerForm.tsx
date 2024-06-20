import React, { FC, useCallback, useEffect, useState } from "react";
import { LabeledSelect } from "@shared/ui/Select";
import { LabeledInput } from "@shared/ui/Input";
import { useForm } from "react-hook-form";
import {
  createWorker,
  updateWorker,
  WorkerData,
  WorkerRole,
} from "@api/workers";
import { useRouter } from "next/router";
import { Select } from "antd";
import { getAllShops } from "@api/shops";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage";

export interface WorkerFormProps {
  data?: any;
}

const accessOptions = [
  { name: "Кассир", value: WorkerRole.WORKER },
  { name: "Менеджер", value: WorkerRole.MANAGER },
];

const WorkerForm: FC<WorkerFormProps> = ({ data }) => {
  const router = useRouter();

  const { handleSubmit, register, reset } = useForm<WorkerData>({
    defaultValues: {
      name: "",
      username: "",
      password: "",
      phone: "",
      role: WorkerRole.WORKER,
    },
  });
  const [confirmPassword, setConfirmPassword] = useState<string>("");
  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const [shopOptions, setShopOptions] = useState<
    { label: string; value: string }[]
  >([]);
  const [selectedShops, setSelectedShops] = useState<string[]>([]);

  const [showPass, setShowPass] = useState<boolean>(false);

  useEffect(() => {
    getAllShops().then((res) => {
      setShopOptions(
        res.data.map((shop: any) => ({
          label: shop.name,
          value: shop.id,
        }))
      );
    });
  }, []);

  useEffect(() => {
    if (!data || !shopOptions) return;
    reset(data);
    setSelectedShops(data.shops);
  }, [data]);

  const onSubmit = useCallback(
    (submitData: WorkerData) => {
      if (submitData.password !== confirmPassword) {
        showAlertMessage("Пароли не совпадают", AlertMessageType.ERROR);
        return;
      }
      if (!data) {
        createWorker({ ...submitData, shops: selectedShops }).then(() =>
          router.replace("/workers")
        );
      } else {
        updateWorker({
          id: data.id,
          ...submitData,
          shops: selectedShops,
        }).then(() => router.replace("/workers"));
      }
    },
    [data, router, selectedShops, confirmPassword]
  );

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
      <LabeledInput {...register("name", { required: true })} label="Имя" />
      <LabeledInput
        {...register("username", { required: true })}
        label="Логин"
      />
      <LabeledInput
        type={showPass ? "text" : "password"}
        {...register("password", { required: true })}
        label="Пароль"
        showPassword={showPass}
        handlePassword={() => {
          setShowPass((prevState) => !prevState);
        }}
      />
      <LabeledInput
        value={confirmPassword}
        onChange={(e) =>
          setConfirmPassword((e.target as HTMLInputElement).value)
        }
        type="password"
        label="Повторите пароль"
      />
      <LabeledInput
        {...register("phone", { required: true })}
        label="Телефон"
      />
      <LabeledSelect
        {...register("role")}
        label="Доступ"
        options={accessOptions}
      />
      <div className="w-full flex items-center">
        <span className="w-40 mr-4">Заведения</span>
        <div className="flex flex-col flex-1 items-start space-y-4">
          <Select
            mode="multiple"
            className="w-full"
            allowClear
            value={selectedShops}
            options={shopOptions}
            onChange={(value) => {
              setSelectedShops(value);
            }}
            optionFilterProp="children"
            filterOption={(input, option) =>
              (option?.label ?? "").toLowerCase().includes(input.toLowerCase())
            }
            filterSort={(optionA, optionB) =>
              (optionA?.label ?? "")
                .toLowerCase()
                .localeCompare((optionB?.label ?? "").toLowerCase())
            }
          />
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

export default WorkerForm;
