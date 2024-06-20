import React, {
  FC,
  FormEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Controller, useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { LabeledSelect } from "@shared/ui/Select";
import { getAllSklads } from "@api/sklad";
import DatePicker from "react-datepicker";
import {
  createInventory,
  InventoryData,
  InventoryItem,
  updateInventory,
  updateInventoryParams,
} from "@api/inventory";
import { LabeledButtonGroup } from "@shared/ui/ButtonGroup";
import NestedCheckbox from "@shared/ui/NestedCheckbox";
import { getAllIngredients } from "@api/ingredient";
import { getAllMenuItems } from "@api/menu-items";
import { Node } from "react-checkbox-tree";
import Search from "@common/Search/Search";
import { getAllIngredientCategories } from "@api/ingredient-category";
import { QueryOptions } from "@api/index";
import { getAllProductCategories } from "@api/product-categories";
import clsx from "clsx";
import Loader from "@common/Loader/Loader";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import AlertMessage from "@common/AlertMessage";
import { getAllInventoryGroups } from "@api/groups";

export interface InventoryFormProps {
  data?: any;
  isEdit?: boolean;
  fromTerminal?: boolean;
}

const dateOptions: { label: string; value: string }[] = [
  {
    label: "Задним числом",
    value: "custom",
  },
  {
    label: "Временем проведения",
    value: "default",
  },
];

const typeOptions: { label: string; value: string }[] = [
  {
    label: "Полная",
    value: "full",
  },
  {
    label: "Полная (без расходников)",
    value: "fullPartial",
  },
  {
    label: "Частичная",
    value: "partial",
  },
];

const getHighlightText = (text: string, keyword: string) => {
  const startIndex = text.toLowerCase().indexOf(keyword.toLowerCase());
  return startIndex !== -1 ? (
    <span>
      {text.substring(0, startIndex)}
      <span style={{ color: "red" }}>
        {text.substring(startIndex, startIndex + keyword.length)}
      </span>
      {text.substring(startIndex + keyword.length)}
    </span>
  ) : (
    <span>{text}</span>
  );
};

const keywordFilter = (nodes: Node[], keyword: string) => {
  let newNodes = [];
  for (let node of nodes) {
    const n = { ...node };
    if (n.children) {
      const nextNodes = keywordFilter(n.children, keyword);
      if (nextNodes.length > 0) {
        n.children = nextNodes;
        // @ts-ignore
      } else if (n.label.toLowerCase().includes(keyword.toLowerCase())) {
        n.children = nextNodes.length > 0 ? nextNodes : [];
      }
      if (
        nextNodes.length > 0 ||
        // @ts-ignore
        n.label.toLowerCase().includes(keyword.toLowerCase())
      ) {
        n.label = getHighlightText(n.label as string, keyword);
        newNodes.push(n);
      }
    } else {
      // @ts-ignore
      if (n.label.toLowerCase().includes(keyword.toLowerCase())) {
        n.label = getHighlightText(n.label as string, keyword);
        newNodes.push(n);
      }
    }
  }
  return newNodes;
};

export const getAllValuesFromNodes = (
  nodes: Node[],
  firstLevel: boolean
): string[] => {
  if (firstLevel) {
    const values = [];
    for (let n of nodes) {
      values.push(n.value);
      if (n.children) {
        values.push(...getAllValuesFromNodes(n.children, false));
      }
    }
    return values;
  } else {
    const values = [];
    for (let n of nodes) {
      values.push(n.value);
      if (n.children) {
        values.push(...getAllValuesFromNodes(n.children, false));
      }
    }
    return values;
  }
};

const InventoryForm: FC<InventoryFormProps> = ({
  data,
  isEdit = false,
  fromTerminal = false,
}) => {
  const router = useRouter();

  const [ingredients, setIngredients] = useState<Node[]>([]);
  const [products, setProducts] = useState<Node[]>([]);
  const [groups, setGroups] = useState<Node[]>([]);
  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const [checkboxNodes, setCheckboxNodes] = useState<Node[]>([
    {
      value: "ingredients",
      label: "Ингредиенты",
      children: [],
    },
    {
      value: "products",
      label: "Товары",
      children: [],
    },
    {
      value: "groups",
      label: "Группы продуктов",
      children: [],
    },
  ]);
  const [checked, setChecked] = useState<string[]>([]);
  const [expanded, setExpanded] = useState<string[]>([]);

  const [skladOptions, setSkladOptions] = useState<any[]>([]);

  const [dateOption, setDateOption] = useState<string>("");
  const [type, setType] = useState<string>("");

  const [submitLoading, setSubmitLoading] = useState(false);

  const [searchValue, setSearchValue] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [inventoryDataLoaded, setInventoryDataLoaded] = useState(false);

  const [itemsLoading, setItemsLoading] = useState<boolean>(false);

  const [defaultValues] = useState<InventoryData>({
    sklad_id: 1,
    time: new Date().toISOString(),
    date: new Date(),
    type: "partial",
    status: "opened",
    items: [],
  });

  const { handleSubmit, register, control, reset, setValue } =
    useForm<InventoryData>({
      defaultValues,
    });

  const getAllItems = useCallback(async (id: number) => {
    setItemsLoading(true);
    await Promise.all([
      getAllIngredientCategories(),
      getAllProductCategories(),
    ]).then(([ingredientCategoriesRes, productCategoriesRes]) => {
      // Set Ingredients
      const ingredientsPromises = ingredientCategoriesRes.data.map(
        (category: any) =>
          getAllIngredients({
            [QueryOptions.CATEGORY]: category.id,
            [QueryOptions.SKLAD]: id,
          })
      );
      Promise.all(ingredientsPromises).then((ingredientsResponses) => {
        setItemsLoading(false);
        ingredientsResponses.forEach((ingredientsRes, index) => {
          const category = ingredientCategoriesRes.data[index];
          const ingredientsOfCategory = ingredientsRes.data;
          ingredientsOfCategory.length > 0 &&
            setIngredients((prevState) => [
              ...prevState,
              {
                value: `ingredients_category_${category.id}`,
                label: category.name,
                children: ingredientsOfCategory.map((ingredient: any) => ({
                  value: `${ingredient.ingredient_id}_ingredient`,
                  label: ingredient.name,
                })),
              },
            ]);
        });
      });

      // Set Products
      const productsPromises = productCategoriesRes.data.map((category: any) =>
        getAllMenuItems({
          [QueryOptions.CATEGORY]: category.id,
          [QueryOptions.SKLAD]: id,
        })
      );
      Promise.all(productsPromises).then((productsResponses) => {
        productsResponses.forEach((productsRes, index) => {
          const category = productCategoriesRes.data[index];
          const productsOfCategory = productsRes.data;
          productsOfCategory.length > 0 &&
            setProducts((prevState) => [
              ...prevState,
              {
                value: `products_category_${category.id}`,
                label: category.name,
                children: productsOfCategory.map((product: any) => ({
                  value: `${product.tovar_id}_tovar`,
                  label: product.name,
                })),
              },
            ]);
        });
      });

      getAllInventoryGroups(id).then((res) => {
        setGroups(
          res.data.map((group: any) => ({
            value: `${group.id}_group`,
            label: group.name,
          }))
        );
      });
      setIsLoading(false);
    });
  }, []);

  useEffect(() => {
    setIsLoading(true);
    getAllSklads().then((skladRes) => {
      // Set Sklad Options
      setSkladOptions(
        skladRes.data.map((item: any) => ({
          name: item.name,
          value: parseInt(item.id),
        }))
      );
      skladRes.data[0] && setValue("sklad_id", skladRes.data[0].id);

      // await getAllItems(skladRes.data[0].id);
    });
  }, []);

  useEffect(() => {
    if (!skladOptions.length) return;
    if (!isEdit) {
      getAllItems(skladOptions[0].value);
      return;
    }
    if (!data) return;
    getAllItems(data.sklad_id);
  }, [isEdit, data, skladOptions]);

  useEffect(() => {
    setCheckboxNodes((prevState) =>
      prevState.map((node) => {
        switch (node.value) {
          case "ingredients":
            return {
              ...node,
              children: ingredients,
            };
          case "products":
            return {
              ...node,
              children: products,
            };
          case "groups":
            return {
              ...node,
              children: groups,
            };
          default:
            return node;
        }
      })
    );
  }, [ingredients, products, groups]);

  // if edit
  useEffect(() => {
    if (!data || !checkboxNodes || inventoryDataLoaded) return;
    const inventoryData = data;
    setInventoryDataLoaded(true);
    setDateOption("custom");
    setType(inventoryData.type);
    setChecked(
      inventoryData.items.map((item: any) => `${item.item_id}_${item.type}`)
    );
    reset({
      sklad_id: inventoryData.sklad_id,
      date: new Date(inventoryData.time),
      type: inventoryData.type,
      time: new Date(inventoryData.time).toISOString(),
      items: inventoryData.items,
      status: "opened",
    });
  }, [data, checkboxNodes, inventoryDataLoaded, reset]);

  const onSubmit = useCallback(
    (submitData: InventoryData) => {
      setSubmitLoading(true);
      const postData = {
        sklad_id: 1,
        ...(!fromTerminal && submitData),
        type,
        items: checked.map((item) => {
          const dataItem = data?.items?.find(
            (dataItem: InventoryItem) =>
              dataItem.item_id === parseInt(item) &&
              dataItem.type === item.split("_")[1]
          );
          return {
            item_id: parseInt(item),
            type: item.split("_")[1],
            time: submitData.date?.toISOString(),
            ...(data
              ? {
                  fact_quantity: dataItem?.fact_quantity || 0,
                  group_id: dataItem?.group_id || 0,
                }
              : {}),
          };
        }),
        time: fromTerminal
          ? new Date().toISOString()
          : submitData.date?.toISOString(),
        status: "opened",
      };
      data
        ? // @ts-ignore
          updateInventoryParams({
            ...postData,
            id: data.id,
            sklad_id: data.sklad_id,
            // @ts-ignore
            items: postData.items.map((item) => ({
              ...item,
              id:
                data.items.find(
                  (dataItem: any) =>
                    dataItem.item_id === item.item_id &&
                    dataItem.type === item.type
                )?.id || 0,
              inventarization_id: data.id,
            })),
          })
            .then((res) => router.push(`/inventory/${res.data.id}`))
            .catch((err) => {
              err.response.status === 445 &&
                showAlertMessage(
                  "Нельзя провести инвентаризацию в промежутке между двумя уже проведенными",
                  AlertMessageType.WARNING
                );
            })
            .finally(() => setSubmitLoading(false))
        : // @ts-ignore
          createInventory(postData)
            .then((res) => {
              return router.push(
                fromTerminal
                  ? `/terminal/inventory/${res.data.id}`
                  : `/inventory/${res.data.id}`
              );
            })
            .catch((err) => {
              err.response.status === 445 &&
                showAlertMessage(
                  "Нельзя провести инвентаризацию в промежутке между двумя уже проведенными",
                  AlertMessageType.WARNING
                );
            })
            .finally(() => {
              setSubmitLoading(false);
            });
    },
    [products, router, checked, type, data]
  );

  const searchedNodes = useMemo(() => {
    return searchValue.trim()
      ? keywordFilter(checkboxNodes, searchValue.trim())
      : checkboxNodes;
  }, [searchValue, checkboxNodes]);

  const onSearchChange = useCallback((e: FormEvent<HTMLInputElement>) => {
    const value = (e.target as HTMLInputElement).value.trim();
    setSearchValue(value);
  }, []);

  useEffect(() => {
    setExpanded(searchValue ? getAllValuesFromNodes(searchedNodes, true) : []);
  }, [searchValue, searchedNodes]);

  return itemsLoading ? (
    <div className="flex flex-col items-center justify-center">
      <Loader />
    </div>
  ) : (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className={clsx(["flex flex-col space-y-5"])}
    >
      {alertMessage && (
        <AlertMessage
          message={alertMessage.message}
          type={alertMessage.type}
          onClose={hideAlertMessage}
        />
      )}
      <div className="flex flex-col space-y-5 w-1/2">
        {!fromTerminal && (
          <>
            <LabeledSelect
              {...register("sklad_id", {
                disabled: !!data,
                valueAsNumber: true,
                onChange: async (e) => {
                  setIngredients([]);
                  setProducts([]);
                  setGroups([]);
                  setChecked([]);
                  await getAllItems(e.target.value);
                },
              })}
              label="Склад"
              options={skladOptions}
            />

            <LabeledButtonGroup
              label="Проверка остатков"
              buttons={dateOptions}
              value={dateOption}
              onBtnClick={(value) => {
                setDateOption(value);
              }}
            />

            {dateOption === "custom" && (
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
                            {monthDate.toLocaleDateString("default", {
                              month: "long",
                            })}{" "}
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
            )}
          </>
        )}

        <LabeledButtonGroup
          label="Тип инвентаризации"
          buttons={typeOptions}
          value={type}
          onBtnClick={(value) => {
            setType(value);
          }}
        />
      </div>

      {type === "partial" && (
        <div className="flex flex-col">
          <div>
            Выберите продукты или категории, чтобы проверить их остатки на
            складе
          </div>
          <Search onChange={onSearchChange} />
          {isLoading ? (
            <Loader />
          ) : (
            <NestedCheckbox
              nodes={checkboxNodes}
              searchedNodes={searchedNodes}
              checked={checked}
              expanded={expanded}
              setChecked={setChecked}
            />
          )}
        </div>
      )}

      <div className="pt-5 border-t border-gray-200">
        <button
          type="submit"
          className="py-2 px-3 bg-primary disabled:bg-gray-400/60 hover:bg-teal-600 transition duration-300 text-white rounded-md"
          disabled={submitLoading}
        >
          Сохранить
        </button>
      </div>
    </form>
  );
};

export default InventoryForm;
