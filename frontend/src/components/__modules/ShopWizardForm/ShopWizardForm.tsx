import React, {
  ChangeEvent,
  FC,
  useCallback,
  useEffect,
  useState,
} from "react";
import { useForm } from "react-hook-form";
import { createShop, ShopFullData, WorkerData } from "@api/shops";
import { Select as CustomSelect } from "@shared/ui/Select";
import { getAllWorkers, WorkerRole } from "@api/workers";
import { getAllAccounts } from "@api/accounts";
import ShopSchetModal from "../ShopSchetModal";
import { useRouter } from "next/router";
import { getMasterMenuItems } from "@api/menu-items";
import {
  TovarData,
  ExistingWorkerData,
  WorkerOption,
} from "@modules/ShopWizardForm/types";
import Checkbox from "@shared/ui/Checkbox/Checkbox";
import clsx from "clsx";
import SelectWithSearch from "@shared/ui/SelectWithSearch";
import { Select } from "antd";
import { XMarkIcon } from "@heroicons/react/24/outline";

const emptyWorker: WorkerData = {
  name: "",
  password: "",
  phone: "",
  username: "",
  role: WorkerRole.WORKER,
  new: true,
};

const roleOptions = [
  { name: "Кассир", value: WorkerRole.WORKER },
  { name: "Менеджер", value: WorkerRole.MANAGER },
];

const ShopWizardForm: FC = () => {
  const router = useRouter();

  const { register, handleSubmit, setValue } = useForm<ShopFullData>({
    defaultValues: {
      shop: {
        name: "",
        address: "",
        tis_token: "",
        cash_schet: 0,
        card_schet: 0,
        limit: 0,
      },
      sklad: {
        name: "",
        address: "",
      },
    },
  });

  const [workers, setWorkers] = useState<WorkerData[]>([emptyWorker]);

  const [schetOptions, setSchetOptions] = useState<
    { name: string; value: number }[]
  >([]);

  const [modalOpen, setModalOpen] = useState<boolean>(false);
  const [schetType, setSchetType] = useState<"cash" | "card">("cash");

  const [tovars, setTovars] = useState<TovarData[]>([]);
  const [selectedTovars, setSelectedTovars] = useState<number[]>([]);

  const [existingWorkers, setExistingWorkers] = useState<WorkerOption[]>([]);
  const [selectedExistingWorkers, setSelectedExistingWorkers] = useState<
    WorkerOption[]
  >([]);

  useEffect(() => {
    getMasterMenuItems().then((res) => {
      setTovars(res.data);
      setSelectedTovars(res.data?.map((tovar: TovarData) => tovar.id) || []);
    });
    getAllWorkers().then((res) => {
      setExistingWorkers(
        res.data.map((worker: ExistingWorkerData) => ({
          label: worker.name,
          value: worker.id,
          data: worker,
        }))
      );
    });
  }, []);

  const handleTovarCheck = useCallback(
    (value: number) => () => {
      setSelectedTovars((prevState) => {
        if (prevState.includes(value)) {
          return prevState.filter((selectedTovar) => selectedTovar !== value);
        }
        return [...prevState, value];
      });
    },
    []
  );

  const handleGetAllAccounts = useCallback(() => {
    getAllAccounts().then((res) => {
      setSchetOptions([
        {
          name: "Выберите счет",
          value: 0,
        },
        ...res.data.map((schet: any) => ({
          name: schet.name,
          value: schet.id !== 0 ? schet.id : -1,
        })),
      ]);
    });
  }, []);

  useEffect(() => {
    handleGetAllAccounts();
  }, [handleGetAllAccounts]);

  const handleWorkerChange = useCallback(
    (fieldName: string, idx: number) =>
      (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setWorkers((prevState) =>
          prevState.map((item, itemIdx) => {
            if (itemIdx !== idx) return item;
            return {
              ...item,
              [fieldName]: e.target.value,
            };
          })
        );
      },
    []
  );

  const handleOpenModal = useCallback(
    (schetType: "cash" | "card") => () => {
      setModalOpen(true);
      setSchetType(schetType);
    },
    []
  );

  const onAccountCreate = useCallback(
    (id: number) => {
      setModalOpen(false);
      handleGetAllAccounts();
      schetType === "cash"
        ? setValue("shop.cash_schet", id)
        : setValue("shop.card_schet", id);
    },
    [handleGetAllAccounts, schetType, setValue]
  );

  const onSubmit = useCallback(
    (submitData: ShopFullData) => {
      console.log();
      createShop({
        ...submitData,
        // @ts-ignore
        workers: [
          ...workers,
          // @ts-ignore
          ...selectedExistingWorkers.map((workerId: number) => {
            const worker = existingWorkers.find(
              (worker) => worker.value === workerId
            );
            const workerData = worker?.data;
            return {
              name: workerData?.name,
              password: workerData?.password,
              phone: workerData?.phone,
              username: workerData?.username,
              role: workerData?.role,
              new: false,
            };
          }),
        ],
        products_shop: {
          tovars: selectedTovars,
          tech_carts: [],
        },
      }).then(() => router.replace("/shops"));
    },
    [workers, selectedExistingWorkers, selectedTovars]
  );

  return (
    <>
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="flex flex-col space-y-6"
      >
        <div className="bg-white px-4 py-5 shadow sm:rounded-lg sm:p-6">
          <div className="md:grid md:grid-cols-3 md:gap-6">
            <div className="md:col-span-1">
              <h3 className="text-base font-semibold leading-6 text-gray-900">
                Заведение
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Базовая информация о заведении
              </p>
            </div>
            <div className="mt-5 space-y-6 md:col-span-2 md:mt-0">
              <div className="col-span-3 sm:col-span-2">
                <label
                  htmlFor="company-website"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Название
                </label>
                <input
                  type="text"
                  {...register("shop.name")}
                  className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>

              <div className="col-span-3 sm:col-span-2">
                <label
                  htmlFor="company-website"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Адрес
                </label>
                <input
                  type="text"
                  {...register("shop.address")}
                  className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>

              <div>
                <label
                  htmlFor="about"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Wipon токен
                </label>
                <input
                  type="text"
                  {...register("shop.tis_token")}
                  className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
                <p className="mt-2 text-sm text-gray-500">
                  Уникальный токен для интеграции с ProSklad
                </p>
              </div>

              <div className="col-span-3 sm:col-span-2">
                <label
                  htmlFor="company-website"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Лимит
                </label>
                <input
                  type="text"
                  {...register("shop.limit", { valueAsNumber: true })}
                  className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white px-4 py-5 shadow sm:rounded-lg sm:p-6">
          <div className="md:grid md:grid-cols-3 md:gap-6">
            <div className="md:col-span-1">
              <h3 className="text-base font-semibold leading-6 text-gray-900">
                Счета
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Наличный и безналичный счета для этого заведения
              </p>
            </div>
            <div className="mt-5 md:col-span-2 md:mt-0">
              <div className="grid grid-cols-6 gap-6">
                <div className="col-span-6">
                  <label
                    htmlFor="country"
                    className="block text-sm font-medium leading-6 text-gray-900"
                  >
                    Наличный счет
                  </label>
                  <div className="flex items-end space-x-3">
                    <CustomSelect
                      id="country"
                      options={schetOptions}
                      {...register("shop.cash_schet", { valueAsNumber: true })}
                      className="mt-2 block w-1/3 rounded-md border-0 bg-white px-2 px-2 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                    />
                    <div className="flex items-end">или</div>
                    <button
                      type="button"
                      onClick={handleOpenModal("cash")}
                      className="bg-primary/80 text-white py-1.5 px-2 rounded text-sm font-medium"
                    >
                      Откройте новый
                    </button>
                  </div>
                </div>

                <div className="col-span-6">
                  <label
                    htmlFor="country"
                    className="block text-sm font-medium leading-6 text-gray-900"
                  >
                    Безналичный счет
                  </label>
                  <div className="flex items-end space-x-3">
                    <CustomSelect
                      id="country"
                      options={schetOptions}
                      {...register("shop.card_schet", { valueAsNumber: true })}
                      className="mt-2 block w-1/3 rounded-md border-0 bg-white px-2 px-2 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                    />
                    <div className="flex items-end">или</div>
                    <button
                      type="button"
                      onClick={handleOpenModal("card")}
                      className="bg-primary/80 text-white py-1.5 px-2 rounded text-sm font-medium"
                    >
                      Откройте новый
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white px-4 py-5 shadow sm:rounded-lg sm:p-6">
          <div className="md:grid md:grid-cols-3 md:gap-6">
            <div className="md:col-span-1">
              <h3 className="text-base font-semibold leading-6 text-gray-900">
                Сотрудники
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Менеджера и кассиры ответственные за это заведение
              </p>
            </div>
            <div className="mt-5 md:col-span-2 md:mt-0">
              <div className="space-y-6 divide-y divide-gray-300">
                {workers.map((worker, idx) => (
                  <div
                    key={idx}
                    className={clsx([
                      "grid grid-cols-6 gap-6",
                      idx !== 0 && "pt-4",
                    ])}
                  >
                    <div className="col-span-6 flex items-center justify-between">
                      <span>Сотрудник {idx + 1}</span>
                      <button
                        type="button"
                        className="p-2 rounded text-red-500 hover:bg-red-500 hover:text-white transition-colors duration-100"
                      >
                        <XMarkIcon className="w-6 h-6" />
                      </button>
                    </div>
                    <div className="col-span-6">
                      <label
                        htmlFor="company-website"
                        className="block text-sm font-medium leading-6 text-gray-900"
                      >
                        ФИО
                      </label>
                      <input
                        type="text"
                        name="worker_name"
                        defaultValue={worker.name}
                        onChange={handleWorkerChange("name", idx)}
                        className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                    <div className="col-span-6 sm:col-span-3">
                      <label
                        htmlFor="company-website"
                        className="block text-sm font-medium leading-6 text-gray-900"
                      >
                        Логин (имя пользователя)
                      </label>
                      <input
                        type="text"
                        name="worker_username"
                        defaultValue={worker.username}
                        onChange={handleWorkerChange("username", idx)}
                        className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                    <div className="col-span-6 sm:col-span-3">
                      <label
                        htmlFor="company-website"
                        className="block text-sm font-medium leading-6 text-gray-900"
                      >
                        Пароль
                      </label>
                      <input
                        type="text"
                        name="password"
                        defaultValue={worker.password}
                        onChange={handleWorkerChange("password", idx)}
                        className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                    <div className="col-span-6">
                      <label
                        htmlFor="company-website"
                        className="block text-sm font-medium leading-6 text-gray-900"
                      >
                        Телефон
                      </label>
                      <input
                        type="text"
                        name="phone"
                        defaultValue={worker.phone}
                        onChange={handleWorkerChange("phone", idx)}
                        className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                    <div className="col-span-6">
                      <label
                        htmlFor="country"
                        className="block text-sm font-medium leading-6 text-gray-900"
                      >
                        Доступ
                      </label>
                      <CustomSelect
                        id="country"
                        name="role"
                        options={roleOptions}
                        defaultValue={worker.role}
                        // @ts-ignore
                        onChange={handleWorkerChange("role", idx)}
                        className="mt-2 block w-full rounded-md border-0 bg-white px-2 px-2 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>
                ))}
              </div>
              <div className="flex flex-col">
                <button
                  type="button"
                  className="self-start px-2 py-1.5 mt-6 font-medium text-primary text-sm border border-primary rounded mt-2 hover:bg-primary/80 hover:text-white"
                  onClick={() => {
                    setWorkers((prevState) => [...prevState, emptyWorker]);
                  }}
                >
                  Добавить сотрудника
                </button>
                <span className="font-light mt-4 mb-1 text-sm">
                  Можно выбрать из существующих:
                </span>
                <Select
                  mode="multiple"
                  allowClear
                  style={{ width: "100%", flex: 1 }}
                  options={existingWorkers}
                  value={selectedExistingWorkers}
                  onChange={(value) => setSelectedExistingWorkers(value)}
                  placeholder="Выберите сотрудников"
                  filterOption={(input, option) =>
                    (option?.label ?? "")
                      .trim()
                      .toLowerCase()
                      .includes(input.trim().toLowerCase())
                  }
                />
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white px-4 py-5 shadow sm:rounded-lg sm:p-6">
          <div className="md:grid md:grid-cols-3 md:gap-6">
            <div className="md:col-span-1">
              <h3 className="text-base font-semibold leading-6 text-gray-900">
                Склад
              </h3>
              <p className="mt-1 text-sm text-gray-500">Склад заведения</p>
            </div>
            <div className="mt-5 space-y-6 md:col-span-2 md:mt-0">
              <div className="grid grid-cols-6 gap-6">
                <div className="col-span-6">
                  <label
                    htmlFor="company-website"
                    className="block text-sm font-medium leading-6 text-gray-900"
                  >
                    Название
                  </label>
                  <input
                    type="text"
                    {...register("sklad.name")}
                    className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                  />
                </div>
                <div className="col-span-6">
                  <label
                    htmlFor="company-website"
                    className="block text-sm font-medium leading-6 text-gray-900"
                  >
                    Адрес
                  </label>
                  <input
                    type="text"
                    {...register("sklad.address")}
                    className="mt-2 block w-full rounded-md border-0 px-2 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white px-4 py-5 shadow sm:rounded-lg sm:p-6">
          <div className="md:grid md:grid-cols-3 md:gap-6">
            <div className="md:col-span-1">
              <h3 className="text-base font-semibold leading-6 text-gray-900">
                Меню
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Выберите <strong>товары</strong> из основного меню для вашего
                заведения
              </p>
            </div>
            <div className="mt-5 space-y-6 md:col-span-2 md:mt-0">
              <ul className="grid grid-cols-6 gap-x-2 gap-y-4">
                {tovars?.map(({ id, name }) => (
                  <div key={id} className="col-span-2 flex items-center">
                    <Checkbox
                      name={name}
                      label={name}
                      onChange={handleTovarCheck(id)}
                      defaultChecked
                    />
                  </div>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <div className="flex justify-end px-4 sm:px-0">
          <button
            type="submit"
            className="ml-3 inline-flex justify-center rounded-md bg-primary/80 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-primary focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
          >
            Добавить
          </button>
        </div>
      </form>
      <ShopSchetModal
        isOpen={modalOpen}
        setIsOpen={setModalOpen}
        type={schetType}
        onCreate={onAccountCreate}
      />
    </>
  );
};

export default ShopWizardForm;
