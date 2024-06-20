import React, {
  ChangeEvent,
  FC,
  useCallback,
  useEffect,
  useState,
} from "react";
import SelectWithSearch from "@shared/ui/SelectWithSearch";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage";

import { Input, LabeledInput } from "@shared/ui/Input";
import { MeasureOption } from "@shared/types/types.";
import { useForm } from "react-hook-form";
import { PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import {
  createMasterModifier,
  createModifier,
  updateMasterModifier,
  updateModifier,
} from "@api/modifiers";
import { useRouter } from "next/router";
import { Select } from "antd";
import { getAllShops } from "@api/shops";
import { getMasterIngredients } from "@api/ingredient";
import useMasterRole from "@hooks/useMasterRole";

export interface NaborFormProps {
  isEdit?: boolean;
  data?: any;
}

const NaborForm: FC<NaborFormProps> = ({ data }) => {
  const router = useRouter();

  const isMaster = useMasterRole();

  const { handleSubmit, register, reset } = useForm({
    defaultValues: {
      name: "",
      min: 0,
      max: 0,
      shops: [],
      ingredient_nabor: [],
      replaces: [],
    },
  });

  const [shopOptions, setShopOptions] = useState<
    {
      value: number;
      label: string;
    }[]
  >([]);
  const [ingredientOptions, setIngredientOptions] = useState<
    {
      name: string;
      value: number;
      data: {
        price: number;
        measure: string;
      };
    }[]
  >([]);
  const [ingredients, setIngredients] = useState<
    {
      ingredient_id: number;
      price: number;
      brutto: number;
      measure: MeasureOption;
    }[]
  >([]);

  const [selectedShops, setSelectedShops] = useState<number[]>([]);
  const [selectedReplaces, setSelectedReplaces] = useState<number[]>([]);

  useEffect(() => {
    if (
      !ingredientOptions.length ||
      !data ||
      !shopOptions.length ||
      isMaster === null
    )
      return;
    const ingredients = isMaster
      ? data.ingredient_nabor
      : data.nabor_ingredient;
    setIngredients(
      ingredients?.map((item: any) => ({
        ingredient_id: item.ingredient_id,
        price: item.price,
        brutto: item.brutto,
        measure: item.measure,
      })) || []
    );
    setSelectedShops(
      typeof data.shop_id === "number" ? [data.shop_id] : data.shop_id || []
    );
    setSelectedReplaces(data.replaces || []);
  }, [data, ingredientOptions, shopOptions, isMaster]);

  const handleGetIngredients = useCallback((res: any) => {
    const ingredients = res.data.map(
      ({
        id,
        name,
        measure,
        cost,
      }: {
        id: number;
        name: string;
        measure: MeasureOption;
        cost: number;
      }) => ({
        name,
        value: id,
        data: {
          measure,
          cost,
        },
      })
    );
    setIngredientOptions(ingredients);
  }, []);

  useEffect(() => {
    getMasterIngredients().then(handleGetIngredients);

    getAllShops().then((res) => {
      const shops = res.data.map(
        ({ id, name }: { id: number; name: string }) => ({
          label: name,
          value: id,
        })
      );
      setShopOptions(shops);
    });
  }, [handleGetIngredients]);

  useEffect(() => {
    ingredientOptions && data && reset({ ...data, category: data.category_id });
  }, [data, ingredientOptions, reset]);

  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const checkIfSubmitionDataIsValid = (submitData: any) => {
    if (submitData.name && submitData.max > 0) {
      return true;
    } else {
      showAlertMessage("Некоторые поля не заполнены", AlertMessageType.WARNING);
      return false;
    }
  };

  const handleSelectIngredient = useCallback(
    (idx: number) => (value: string | number) => {
      const selectedOption = ingredientOptions.find(
        (option) => option.value === value
      );
      const ingredientData = selectedOption?.data || null;
      setIngredients((prevState) =>
        prevState.map((option, i) => {
          if (idx !== i) return option;
          return {
            ingredient_id: value as number,
            brutto: option.brutto,
            measure:
              (ingredientData && (ingredientData.measure as MeasureOption)) ||
              MeasureOption.DEFAULT,
            price: (ingredientData && ingredientData.price) || 0,
          };
        })
      );
    },
    [ingredientOptions]
  );

  const handleBruttoChange = useCallback(
    (idx: number) => (e: ChangeEvent<HTMLInputElement>) => {
      setIngredients((prevState) =>
        prevState.map((option, i) => {
          if (idx !== i) return option;
          return {
            ingredient_id: option.ingredient_id,
            brutto:
              parseFloat(
                (e.target as HTMLInputElement).value
                  .replace(",", ".")
                  .replace(" ", "")
              ) || 0,
            measure: option.measure,
            price: option.price,
          };
        })
      );
    },
    []
  );

  const handlePriceChange = useCallback(
    (idx: number) => (e: ChangeEvent<HTMLInputElement>) => {
      setIngredients((prevState) =>
        prevState.map((option, i) => {
          if (idx !== i) return option;
          return {
            ingredient_id: option.ingredient_id,
            brutto: option.brutto,
            measure: option.measure,
            price:
              parseFloat(
                (e.target as HTMLInputElement).value
                  .replace(",", ".")
                  .replace(" ", "")
              ) || 0,
          };
        })
      );
    },
    []
  );

  const onSubmit = useCallback(
    async (submitData: any) => {
      if (!checkIfSubmitionDataIsValid(submitData)) return;
      const processedData = {
        ...submitData,
        min: parseInt(submitData.min),
        max: parseInt(submitData.max),
        shops: selectedShops,
        ingredient_nabor: ingredients,
        replaces: selectedReplaces,
      };
      const processedDataForUpdate = {
        ...processedData,
        id: parseInt(router.query.id as string),
      };

      if (isMaster) {
        data
          ? await updateMasterModifier(processedDataForUpdate)
          : await createMasterModifier(processedData);
      } else {
        data
          ? await updateModifier(processedDataForUpdate)
          : await createModifier(processedData);
      }

      await router.push("/nabors");
    },
    [data, ingredients, router, selectedReplaces, selectedShops, isMaster]
  );

  return (
    <>
      {alertMessage && (
        <AlertMessage
          message={alertMessage.message}
          type={alertMessage.type}
          onClose={hideAlertMessage}
        />
      )}
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="flex flex-col space-y-5"
      >
        <div className="w-1/2 flex flex-col space-y-5">
          <LabeledInput {...register("name")} label="Название" />
          {!isMaster && (
            <div className="w-full flex items-center pt-2">
              <label className="w-40 mr-4">Заведения</label>
              <Select
                mode="multiple"
                allowClear
                disabled={!isMaster}
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
          <LabeledInput
            {...register("min", { valueAsNumber: true })}
            label="Мин. кол-во"
            disabled={!isMaster}
          />
          <LabeledInput
            {...register("max", { valueAsNumber: true })}
            label="Макс. кол-во"
            disabled={!isMaster}
          />
          <div className="w-full flex items-center pt-2">
            <label className="w-40 mr-4">Заменяет след. ингредиенты</label>
            <Select
              mode="multiple"
              disabled={!isMaster}
              placeholder="..."
              value={selectedReplaces}
              style={{ width: "100%", flex: 1 }}
              onChange={(value) => {
                setSelectedReplaces(value);
              }}
              options={ingredientOptions.map((option) => ({
                label: option.name,
                value: option.value,
              }))}
              filterOption={(input, option) =>
                (option?.label ?? "")
                  .trim()
                  .toLowerCase()
                  .includes(input.trim().toLowerCase())
              }
            />
          </div>
        </div>
        {ingredientOptions.length > 0 && (
          <>
            <div className="w-full flex flex-col border-t border-gray-300 pt-3 space-y-3">
              <div className="font-medium text-lg">Модификаторы</div>
              <div className="w-2/3 flex flex-col space-y-3">
                {
                  <div className="w-full flex flex-col space-y-3">
                    <div className="w-full flex items-center border-b border-gray-200 pb-2">
                      <div className="w-1/2 text-sm text-gray-500 font-medium">
                        Модификатор
                      </div>
                      <div className="w-1/4 text-sm text-gray-500 font-medium">
                        Брутто
                      </div>
                      <div className="w-1/4 text-sm text-gray-500 font-medium">
                        Цена
                      </div>
                      <div className="w-7" />
                    </div>
                    {ingredients.map((ingredient, idx) => (
                      <div
                        className="w-full flex justify-between items-center"
                        key={`${ingredient.ingredient_id}_${idx}`}
                      >
                        <div className="w-1/2 flex pr-4">
                          <SelectWithSearch
                            disabled={!isMaster}
                            className="w-full"
                            options={ingredientOptions}
                            value={ingredient.ingredient_id}
                            onChange={handleSelectIngredient(idx)}
                          />
                        </div>
                        <div className="w-1/4 flex items-center pr-4">
                          <Input
                            type="text"
                            name="brutto"
                            disabled={!isMaster}
                            className="grow"
                            defaultValue={ingredients[idx].brutto}
                            onChange={handleBruttoChange(idx)}
                          />
                          <div className="ml-2">{ingredients[idx].measure}</div>
                        </div>
                        <div className="w-1/4 flex items-center pr-2">
                          <Input
                            type="text"
                            name="brutto"
                            disabled={!isMaster}
                            className="grow"
                            value={ingredients[idx].price}
                            onChange={handlePriceChange(idx)}
                          />
                        </div>
                        <button
                          onClick={() => {
                            setIngredients((prevState) =>
                              prevState.filter((_, i) => idx !== i)
                            );
                          }}
                          type="button"
                          disabled={!isMaster}
                          className="p-2 rounded-md hover:bg-gray-200 transition duration-100"
                        >
                          <XMarkIcon className="w-5 h-5" />
                        </button>
                      </div>
                    ))}
                  </div>
                }
                {isMaster && (
                  <button
                    onClick={() => {
                      //@ts-ignore
                      setIngredients((prevState) => [
                        ...prevState,
                        {
                          value: "",
                          ingredient_id: 1,
                          brutto: 0,
                          cost: 0,
                          measure: MeasureOption.DEFAULT,
                        },
                      ]);
                    }}
                    type="button"
                    className="flex items-center space-x-1 text-indigo-500 hover:text-indigo-700"
                  >
                    <PlusIcon className="w-4 h-4" />{" "}
                    <span>Добавить ингредиент</span>
                  </button>
                )}
              </div>
            </div>
          </>
        )}
        <div className="pt-5 border-t border-gray-200">
          <button
            type="submit"
            className="py-2 px-3 bg-primary hover:bg-teal-600 transition duration-300 text-white rounded-md"
          >
            Сохранить
          </button>
        </div>
      </form>
    </>
  );
};

export default NaborForm;
