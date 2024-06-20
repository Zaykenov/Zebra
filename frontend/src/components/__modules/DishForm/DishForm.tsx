import React, { FC, useCallback, useEffect, useState } from "react";
import SelectWithSearch from "@shared/ui/SelectWithSearch";
import ModifierFormModal from "../ModifierFormModal";
import LabeledSelect from "@shared/ui/Select/LabeledSelect";
import Checkbox from "@shared/ui/Checkbox/Checkbox";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import ImageUpload from "@common/ImageUpload/ImageUpload";
import AlertMessage from "@common/AlertMessage";

import { Input, LabeledInput } from "@shared/ui/Input";
import { MeasureOption, measureOptions } from "@shared/types/types.";
import { Controller, useForm } from "react-hook-form";
import { PlusIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { getAllProductCategories } from "@api/product-categories";
import { createAllMasterModifiers, createAllModifiers } from "@api/modifiers";
import {
  createDish,
  createMasterDish,
  updateDish,
  updateMasterDish,
} from "@api/dishes";
import { useRouter } from "next/router";
import { ModifierData } from "../ModifierForm/types";
import { Select } from "antd";
import { getAllShops } from "@api/shops";
import { getMasterIngredients } from "@api/ingredient";
import useMasterRole from "@hooks/useMasterRole";

export interface DishFormProps {
  data?: any;
}

const DishForm: FC<DishFormProps> = ({ data }) => {
  const router = useRouter();

  const isMaster = useMasterRole();

  const [categoryOptions, setCategoryOptions] = useState([]);
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
        cost: number;
        measure: string;
      };
    }[]
  >([]);
  const [ingredients, setIngredients] = useState<
    {
      value: number;
      ingredient_id: number;
      cost: number;
      brutto: number;
      measure: MeasureOption;
    }[]
  >([]);

  const [uploadedImage, setUploadedImage] = useState<string | null>(null);

  const [selectedShops, setSelectedShops] = useState<number[]>([]);

  useEffect(() => {
    if (!ingredientOptions.length || !data || !shopOptions.length) return;
    setIngredients(
      data.ingredient_tech_cart?.map((item: any) => ({
        value: item.id,
        cost: item.cost,
        ingredient_id: item.id,
        brutto: item.brutto,
        measure: item.measure,
      })) || []
    );
    setUploadedImage(data.image);
    setModifierSets(
      data.nabors.map((nabor: any) => ({
        ...nabor,
        ingredient_nabor: nabor.nabor_ingredient.map((nabor_item: any) => ({
          value: nabor_item.id,
          ingredient_id: nabor_item.id,
          brutto: nabor_item.brutto,
          price: nabor_item.price,
          measure: nabor_item.measure,
        })),
      }))
    );
    setSelectedShops(
      typeof data.shop_id === "number" ? [data.shop_id] : data.shop_id || []
    );
  }, [data, ingredientOptions, shopOptions]);

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
    getMasterIngredients().then((res) => {
      const ingredients = res.data.map(
        ({
          id,
          name,
          measure,
          cost,
        }: {
          id: string;
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
  }, [data]);

  const { handleSubmit, register, reset, control } = useForm({
    defaultValues: {
      name: "",
      category: "",
      image: "",
      tax: "Фискальный налог",
      measure: MeasureOption.DEFAULT,
      discount: false,
      price: 0,
    },
  });

  useEffect(() => {
    ingredientOptions && data && reset({ ...data, category: data.category_id });
  }, [data, ingredientOptions, reset]);

  const [isModifierModalOpen, setIsModifierModalOpen] = useState(false);
  const [modifierSets, setModifierSets] = useState<ModifierData[]>([]);
  const [selectedModifier, setSelectedModifier] = useState<ModifierData | null>(
    null
  );
  const [totalCost, setTotalCost] = useState<number>(0);
  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  useEffect(() => {
    !isModifierModalOpen && setSelectedModifier(null);
  }, [isModifierModalOpen]);

  const addModifierSet = useCallback((modifierData: ModifierData) => {
    setModifierSets((prevState) => [...prevState, modifierData]);
  }, []);
  const updateModifierSet = useCallback((modifierData: ModifierData) => {
    setModifierSets((prevState) =>
      prevState.map((modifierSet) => {
        if (modifierSet.name !== modifierData.name) return modifierSet;
        return modifierData;
      })
    );
  }, []);
  const removeModifierSet = useCallback((name: string) => {
    setModifierSets((prevState) =>
      prevState.filter((modifierSet) => modifierSet.name !== name)
    );
  }, []);

  const checkIfSubmitionDataIsValid = (submitData: any) => {
    if (submitData.name && submitData.price > 0) {
      // if (ingredients.length > 0 && submitData.name && submitData.price > 0) {
      return true;
    } else {
      showAlertMessage("Некоторые поля не заполнены", AlertMessageType.WARNING);
      return false;
    }
  };

  const [role, setRole] = useState<string>("");

  useEffect(() => {
    const role = localStorage.getItem("zebra.role");
    role && setRole(role);
  }, []);

  const onSubmit = useCallback(
    async (submitData: any) => {
      if (!checkIfSubmitionDataIsValid(submitData)) return;
      const res = isMaster
        ? await createAllMasterModifiers(modifierSets)
        : await createAllModifiers(modifierSets);
      if (res.error) return;
      const nabor =
        res.data.length > 0
          ? res.data.map((item: ModifierData) => ({
              nabor_id: item.id,
            }))
          : [];
      const postData = {
        ...submitData,

        image: uploadedImage,
        category: parseInt(submitData.category),
        ingredient_tech_cart:
          ingredients?.map(({ ingredient_id, brutto }) => ({
            ingredient_id,
            brutto,
          })) || [],
        nabor,
      };
      if (!isMaster) {
        postData.shop_id =
          selectedShops.length > 0
            ? selectedShops
            : shopOptions.map((shop) => shop.value);
      }
      data
        ? isMaster
          ? await updateMasterDish({ ...postData, id: data.id })
          : await updateDish({ ...postData, id: data.id })
        : isMaster
        ? await createMasterDish(postData)
        : await createDish(postData);

      await router.push("/dishes");
      return;
    },
    [modifierSets, uploadedImage, ingredients, router, selectedShops, isMaster]
  );

  const getTotalCost = () => {
    let total = 0;
    ingredients.map((ingredient) => {
      //@ts-ignore
      total +=
        ingredient.cost && ingredient.brutto
          ? Math.round(
              (ingredient.cost * (ingredient.brutto as number) +
                Number.EPSILON) *
                100
            ) / 100
          : 0;
    });
    return total;
  };

  useEffect(() => {
    setTotalCost(getTotalCost());
  }, [ingredients]);

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
          <div className="w-full flex items-center pt-2">
            <label className="w-40 mr-4">Категория</label>
            <Controller
              control={control}
              {...register("category")}
              render={({ field: { onChange, value } }) => (
                <Select
                  allowClear
                  disabled={!isMaster}
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
            disabled={!isMaster}
            label="Ед. измерения"
            options={measureOptions}
          />
          <div className="w-full flex items-center pt-2">
            <label className="w-40 mr-4">Участвует в скидках</label>
            <Checkbox {...register("discount")} disabled={!isMaster} />
          </div>
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
          <div className="flex">
            <span className="w-40 mr-4">Цена</span>
            <div className="flex items-center space-x-4">
              <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
                <label className="text-xs text-gray-600">Итого</label>
                <div className="flex items-center space-x-2">
                  <Input
                    {...register("price", { valueAsNumber: true })}
                    type="text"
                    className="text-right"
                  />{" "}
                  <span className="text-lg">₸</span>
                </div>
              </div>
              <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
                <label className="text-xs text-gray-600">
                  Наценка до налога
                </label>
                <div className="flex items-center space-x-2">
                  <label>
                    {data &&
                      (totalCost !== 0
                        ? Math.round((data.profit / totalCost) * 100)
                        : 0)}
                  </label>
                  <span className="text-lg">%</span>
                </div>
              </div>
              <div className="flex flex-col space-y-1 w-28 whitespace-nowrap">
                <label className="text-xs text-gray-600">
                  Себестоимость без НДС
                </label>
                <div className="flex items-center space-x-2">
                  <label>{totalCost.toFixed(2)}</label>
                  <span className="text-lg">₸</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        {ingredientOptions.length > 0 && (
          <>
            <div className="flex flex-col space-y-3 border-t border-gray-300 pt-3">
              <label className="w-40 font-medium text-lg">Состав</label>
              <div className="w-2/3 flex flex-col space-y-3">
                {ingredients && ingredients.length > 0 && (
                  <div className="w-full flex flex-col space-y-3">
                    <div className="w-full flex items-center border-b border-gray-200 pb-2">
                      <div className="w-1/2 text-sm text-gray-500 font-medium">
                        Ингредиент
                      </div>
                      <div className="w-1/4 text-sm text-gray-500 font-medium">
                        Брутто
                      </div>
                      <div className="w-1/4 text-sm text-gray-500 font-medium">
                        Себестоимость без НДС
                      </div>
                      <div className="w-7" />
                    </div>
                    {ingredients.map((ingredient, idx) => (
                      <div
                        className="w-full flex justify-between items-center"
                        key={`ingredient_${ingredient.ingredient_id}_${idx}`}
                      >
                        <div className="w-1/2 flex pr-4">
                          <SelectWithSearch
                            className="w-full"
                            disabled={!isMaster}
                            options={ingredientOptions}
                            value={ingredient.value}
                            onChange={(value) => {
                              const selectedOption = ingredientOptions.find(
                                // @ts-ignore
                                (option) => option.value === value
                              );
                              const ingredientData =
                                selectedOption?.data || null;
                              //@ts-ignore
                              setIngredients((prevState) =>
                                prevState.map((option, i) => {
                                  if (idx !== i) return option;
                                  return {
                                    value: value,
                                    ingredient_id: value,
                                    brutto: option.brutto,
                                    measure:
                                      ingredientData && ingredientData.measure,
                                    cost: ingredientData && ingredientData.cost,
                                  };
                                })
                              );
                            }}
                          />
                        </div>
                        <div className="w-1/4 flex items-center pr-4">
                          <Input
                            type="text"
                            name="brutto"
                            className="grow"
                            disabled={!isMaster}
                            defaultValue={ingredients[idx].brutto}
                            onChange={(e) => {
                              //@ts-ignore
                              setIngredients((prevState) =>
                                prevState.map((option, i) => {
                                  if (idx !== i) return option;
                                  return {
                                    value: option.value,
                                    ingredient_id: option.ingredient_id,
                                    brutto:
                                      parseFloat(
                                        (e.target as HTMLInputElement).value
                                          .replace(",", ".")
                                          .replace(" ", "")
                                      ) || "",
                                    measure: option.measure,
                                    cost: option.cost,
                                  };
                                })
                              );
                            }}
                          />
                          <div className="ml-2">{ingredients[idx].measure}</div>
                        </div>
                        <div className="w-1/4 flex items-center pr-2">
                          <Input
                            type="text"
                            name="brutto"
                            className="grow"
                            disabled
                            value={(
                              ingredients[idx].cost * ingredients[idx].brutto
                            ).toFixed(2)}
                          />
                        </div>
                        <button
                          onClick={() => {
                            setIngredients((prevState) =>
                              prevState.filter((_, i) => idx !== i)
                            );
                          }}
                          type="button"
                          className="p-2 rounded-md hover:bg-gray-200 transition duration-100"
                        >
                          <XMarkIcon className="w-5 h-5" />
                        </button>
                      </div>
                    ))}
                  </div>
                )}
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
            <div className="w-full flex flex-col border-t border-gray-300 pt-3 space-y-3">
              <div className="font-medium text-lg">Модификаторы</div>
              {modifierSets &&
                modifierSets?.map((modifierSet, idx) => (
                  <div
                    key={modifierSet.name + idx}
                    className="w-full flex flex-col border-b border-gray-200 pb-2"
                  >
                    <div className="flex items-center space-x-3 mb-3">
                      <span className="text-lg font-medium">
                        {modifierSet.name}
                      </span>
                      <button
                        onClick={(e) => {
                          e.preventDefault();
                          setSelectedModifier(modifierSet);
                          setIsModifierModalOpen(true);
                        }}
                        className="text-indigo-400 hover:underline hover:text-indigo-500 text-sm"
                      >
                        Редактировать
                      </button>
                    </div>
                    <div className="w-2/3 flex flex-col">
                      <div className="w-full flex items-center border-b border-gray-200 pb-2 mb-3">
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
                      {modifierSet?.ingredient_nabor?.map(
                        (modifier, modifierIdx) => (
                          <div
                            className="w-full flex justify-between items-center mb-3"
                            key={`modifier_${modifierSet.name}_${modifier.ingredient_id}_${idx}`}
                          >
                            <div className="w-1/2 pr-4">
                              <SelectWithSearch
                                disabled={!isMaster}
                                className="w-full"
                                options={ingredientOptions}
                                value={modifier.value}
                                onChange={(value) => {
                                  const selectedOption = ingredientOptions.find(
                                    // @ts-ignore
                                    (option) => option.value === value
                                  );
                                  const ingredientData =
                                    selectedOption?.data || null;
                                  // @ts-ignore
                                  setModifierSets((prevState) =>
                                    prevState.map((prevModifierSet, pmsIdx) => {
                                      if (pmsIdx !== idx)
                                        return prevModifierSet;
                                      return {
                                        ...prevModifierSet,
                                        ingredient_nabor:
                                          prevModifierSet.ingredient_nabor.map(
                                            (ingr, ingrIdx) => {
                                              if (ingrIdx !== modifierIdx)
                                                return ingr;
                                              return {
                                                ...ingr,
                                                value: value,
                                                ingredient_id: value,
                                                measure:
                                                  ingredientData?.measure || "",
                                              };
                                            }
                                          ),
                                      };
                                    })
                                  );
                                }}
                              />
                            </div>
                            <div className="w-1/4 pr-4 flex items-center">
                              <Input
                                type="text"
                                name="brutto"
                                className="grow"
                                defaultValue={modifier.brutto}
                                disabled={!isMaster}
                                onChange={(e) => {
                                  setModifierSets((prevState) =>
                                    prevState.map((prevModifierSet, pmsIdx) => {
                                      if (pmsIdx !== idx)
                                        return prevModifierSet;
                                      return {
                                        ...prevModifierSet,
                                        ingredient_nabor:
                                          prevModifierSet.ingredient_nabor.map(
                                            (ingr, ingrIdx) => {
                                              if (ingrIdx !== modifierIdx)
                                                return ingr;
                                              return {
                                                ...ingr,
                                                brutto: parseFloat(
                                                  (
                                                    e.target as HTMLInputElement
                                                  ).value
                                                    .replace(",", ".")
                                                    .replace(" ", "")
                                                ),
                                              };
                                            }
                                          ),
                                      };
                                    })
                                  );
                                }}
                              />
                              <div className="ml-2">{modifier.measure}</div>
                            </div>
                            <div className="w-1/4 pr-2">
                              <Input
                                type="text"
                                name="price"
                                className="w-full"
                                disabled={!isMaster}
                                defaultValue={modifier.price}
                                onChange={(e) => {
                                  setModifierSets((prevState) =>
                                    prevState.map((prevModifierSet, pmsIdx) => {
                                      if (pmsIdx !== idx)
                                        return prevModifierSet;
                                      return {
                                        ...prevModifierSet,
                                        ingredient_nabor:
                                          prevModifierSet.ingredient_nabor.map(
                                            (ingr, ingrIdx) => {
                                              if (ingrIdx !== modifierIdx)
                                                return ingr;
                                              return {
                                                ...ingr,
                                                price: parseFloat(
                                                  (
                                                    e.target as HTMLInputElement
                                                  ).value
                                                    .replace(",", ".")
                                                    .replace(" ", "")
                                                ),
                                              };
                                            }
                                          ),
                                      };
                                    })
                                  );
                                }}
                              />
                            </div>
                            <div className="w-5 flex items-center justify-center">
                              <button
                                disabled={!isMaster}
                                onClick={() => {
                                  setModifierSets((prevState) =>
                                    prevState.map((prevModifierSet, pmsIdx) => {
                                      if (pmsIdx !== idx)
                                        return prevModifierSet;
                                      return {
                                        ...prevModifierSet,
                                        ingredient_nabor:
                                          prevModifierSet.ingredient_nabor.filter(
                                            (_, ingrIdx) =>
                                              ingrIdx !== modifierIdx
                                          ),
                                      };
                                    })
                                  );
                                }}
                                type="button"
                                className="p-2 rounded-md hover:bg-gray-200 transition duration-100"
                              >
                                <XMarkIcon className="w-5 h-5" />
                              </button>
                            </div>
                          </div>
                        )
                      )}
                    </div>
                    {isMaster && (
                      <button
                        onClick={() => {
                          setModifierSets((prevState) =>
                            prevState.map((prevModifierSet, pmsIdx) => {
                              if (pmsIdx !== idx) return prevModifierSet;
                              return {
                                ...prevModifierSet,
                                ingredient_nabor: [
                                  ...prevModifierSet.ingredient_nabor,
                                  {
                                    value: 1,
                                    ingredient_id: 1,
                                    price: 0,
                                    brutto: 0,
                                    measure: MeasureOption.DEFAULT,
                                  },
                                ],
                              };
                            })
                          );
                        }}
                        type="button"
                        className="flex items-center text-sm space-x-1 text-indigo-500 hover:text-indigo-700 mt-3 mb-2"
                      >
                        <PlusIcon className="w-4 h-4" />{" "}
                        <span>
                          Добавить модификатор в набор &quot;{modifierSet.name}
                          &quot;
                        </span>
                      </button>
                    )}
                  </div>
                ))}
              {isMaster && (
                <div className="w-full flex flex-col p-2 rounded-md bg-gray-100">
                  <button
                    onClick={() => {
                      setIsModifierModalOpen(true);
                    }}
                    type="button"
                    className="flex items-center text-sm space-x-1 text-indigo-500 hover:text-indigo-700"
                  >
                    <PlusIcon className="w-4 h-4" />{" "}
                    <span>Добавить набор модификаторов</span>
                  </button>
                </div>
              )}
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
      <ModifierFormModal
        isOpen={isModifierModalOpen}
        setIsOpen={setIsModifierModalOpen}
        data={selectedModifier}
        addModifierSet={addModifierSet}
        updateModifierSet={updateModifierSet}
        removeModifierSet={removeModifierSet}
      />
    </>
  );
};

export default DishForm;
