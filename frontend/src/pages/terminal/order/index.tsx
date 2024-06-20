import React, { useCallback, useEffect, useState } from "react";
import { NextPage } from "next";
import {
  Category,
  Modificator,
  NaborIngredient,
  Product,
  getTerminalStartData,
} from "@api/terminal";
import TerminalLayout from "@layouts/TerminalLayout";
import useAlertMessage, { AlertMessageType } from "@hooks/useAlertMessage";
import clsx from "clsx";
import { ChevronLeftIcon } from "@heroicons/react/24/solid";
import AlertMessage from "@common/AlertMessage";
import DiscountModal from "@modules/DiscountModal";
import CodeModal from "@modules/CodeModal";
import CommentCheckModal from "@modules/CommentCheckModal";
import ModificatorOptionsModal from "@modules/ModificatorOptionsModal";
import PaymentModal from "@modules/PaymentModal";
import {
  CheckItemData,
  CheckStatus,
  TechCartCheckData,
  TovarCheckData,
  createCheck,
} from "@api/check";
import { PaymentMethod } from "@api/check";
import useLocalStorage from "@hooks/useLocalStorage";
import { UserMobileData } from "@api/mobile";
import InfoModal from "@modules/InfoModal";
import { checkShift } from "@api/shifts";
import useOnlineStatus from "@hooks/useOnlineStatus";
import TerminalModalButton from "@modules/TerminalButtons/TerminalModalButton";
import ProductList from "@modules/TerminalLists/ProductList";
import CategoryList from "@modules/TerminalLists/CategoryList";
import MainPanelProductList from "@modules/TerminalLists/MainPanelProductList";
import { regenerateIdempotency } from "@reducers/idempotencySlice";
import { useAppDispatch } from "@hooks/useAppDispatch";
import { useAppSelector } from "@hooks/useAppSelector";
import { processPrintObject } from "@utils/processPrintObject";

export function prepareCheck(
  status: CheckStatus,
  payment: PaymentMethod,
  card: number,
  cash: number,
  check: CheckItemData[],
  userData: UserMobileData,
  discount: number,
  comment: string,
  orderId: number | null,
  pager: number,
) {
  let tovarCheck: TovarCheckData[] = [];
  let techCartCheck: TechCartCheckData[] = [];
  check.forEach((item) => {
    if (item.type === "tovar") {
      tovarCheck.push({
        tovar_id: item.id,
        quantity: item.count,
        comments: item.comments,
        modifications: "",
      });
    } else {
      // type === "tech"
      techCartCheck.push({
        tech_cart_id: item.id,
        quantity: item.count,
        comments: item.comments,
        modificators: item.selectedModificators || [], // or empty array [] if not defined
      });
    }
  });
  const data = {
    discount_percent: userData ? userData.discount : discount / 100,
    payment,
    cash,
    card,
    comment: `Пейджер №${pager}: ${comment}`,
    status,
    mobile_user_id: userData ? userData.userId : null,
    tovarCheck,
    techCartCheck,
  };
  return orderId ? { ...data, id: orderId } : data;
}

const compareModificators = (productA: any, productB: any) => {
  if (!productA.selectedModificators || !productB.selectedModificators)
    return false;
  const n = productA.selectedModificators.length;
  if (n !== productB.selectedModificators.length) return false;
  const A = productA.selectedModificators.map((item: any) => item.id).sort();
  const B = productB.selectedModificators.map((item: any) => item.id).sort();

  for (let i = 0; i < n; i++) {
    if (A[i] !== B[i]) return false;
  }
  return true;
};

const OrderPage: NextPage = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [shiftData, setShiftData] = useState<any>([]);
  const [check, setCheck] = useState<CheckItemData[]>([]);
  //ids
  const [checkItemId, setCheckItemId] = useState<number | null>(null);
  const [orderId, setOrderId] = useState<number | null>(null);

  const [discount, setDiscount] = useState<number>(0);
  const [totalPrice, setTotalPrice] = useState<number>(0);
  const [totalPriceWithDiscount, setTotalPriceWithDiscount] =
    useState<number>(0);

  const { alertMessage, showAlertMessage, hideAlertMessage } =
    useAlertMessage();

  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    null,
  );

  const [selectedProduct, setSelectedProduct] = useState<Product>();
  const [selectedModificators, setSelectedModificators] = useState<any[]>([]);
  const [mainPanelItems, setMainPanelItems] = useState<any[]>([]);

  const [comment, setComment] = useState<string>("");
  const [userData, setUserData] = useState<any>(null);

  const [paymentDone, setPaymentDone] = useState<boolean>(false);
  const [clearPaymentModalsStates, setClearPaymentModalsStates] =
    useState(false);

  const isOnline = useOnlineStatus();

  const [savedCategories, setSavedCategories] = useLocalStorage(
    "zebra.categories",
    "",
  );
  const [savedMainPanelItems, setSavedMainPanelItems] = useLocalStorage(
    "zebra.mainPanel",
    "",
  );

  const dispatch = useAppDispatch();
  const idempotencyState = useAppSelector((state) => state.idempotency);

  // modal states
  const [modificatorsModalOpen, setModificatorsModalOpen] =
    useState<boolean>(false);
  const [commentModalOpen, setCommentModalOpen] = useState<boolean>(false);
  const [codeModalOpen, setCodeModalOpen] = useState<boolean>(false);
  const [paymentModalOpen, setPaymentModalOpen] = useState<boolean>(false);
  const [discountModalOpen, setDiscountModalOpen] = useState<boolean>(false);
  const [infoModalOpen, setInfoModalOpen] = useState<boolean>(false);

  useEffect(() => {
    if (savedCategories.length == 0 && savedMainPanelItems.length == 0) {
      getTerminalStartData().then((res) => {
        const currCategories = res.categories.filter(
          (categoryItem: any) => categoryItem.category !== "Главный экран",
        );
        setSavedCategories(currCategories);
        setSavedMainPanelItems(res.mainDisplay);
        setCategories(currCategories);
        setMainPanelItems(res.mainDisplay);
      });
    } else {
      setCategories(savedCategories);
      setMainPanelItems(savedMainPanelItems);
    }
    checkShift().then((res) => {
      setShiftData(res.data);
    });
  }, []);

  useEffect(() => {
    const isAuthedRightNow = localStorage.getItem("zebra.authed");
    if (!!isAuthedRightNow) setInfoModalOpen(true);
  }, []);

  const getTotalPriceWithDiscountFromOrder = () => {
    let totalPrice = 0;
    check?.map((product) => {
      let sumOfModificatorPrice =
        product.selectedModificators?.reduce(
          (acc, mod) => acc + mod.price * mod.quantity,
          0,
        ) || 0;
      sumOfModificatorPrice *= product.count;
      totalPrice += product.hasDiscount
        ? (product.total - sumOfModificatorPrice) * (1 - discount / 100) +
          sumOfModificatorPrice
        : product.total;
    });
    return totalPrice;
  };

  const updateCheck = useCallback(() => {
    const totalPriceFromOrder = check.reduce(
      (totalPrice, product) => totalPrice + product.total,
      0,
    );
    setTotalPrice(Math.floor(totalPriceFromOrder));
    setTotalPriceWithDiscount(Math.floor(getTotalPriceWithDiscountFromOrder()));
  }, [check, discount, userData, comment, orderId]);

  useEffect(() => {
    if (check.length === 0) {
      setTotalPrice(0);
      setTotalPriceWithDiscount(0);
      return;
    }
    updateCheck();
  }, [check, updateCheck]);

  useEffect(() => {
    if (clearPaymentModalsStates) {
      setClearPaymentModalsStates(false);
    }
  }, [clearPaymentModalsStates]);

  const getModificatorsByProductId = (productId: number): Modificator[] => {
    for (const category of categories) {
      for (const product of category.products) {
        if (product.id === productId && product.nabor) {
          const modifiers: Modificator[] = [];

          for (const modifier of product.nabor) {
            const uniqueNaborIngredients: NaborIngredient[] = [];

            for (const ingredient of modifier.nabor_ingredient) {
              if (!uniqueNaborIngredients.some((i) => i.id === ingredient.id)) {
                uniqueNaborIngredients.push(ingredient);
              }
            }

            modifier.nabor_ingredient = uniqueNaborIngredients;
            modifiers.push(modifier);
          }

          return modifiers;
        }
      }
    }

    return [];
  };

  const getModificatorsByMainPanelProductId = (
    productId: number,
  ): Modificator[] => {
    for (const product of mainPanelItems) {
      if (product.id === productId && product.nabor) {
        const modifiers: Modificator[] = [];

        for (const modifier of product.nabor) {
          const uniqueNaborIngredients: NaborIngredient[] = [];

          for (const ingredient of modifier.nabor_ingredient) {
            if (!uniqueNaborIngredients.some((i) => i.id === ingredient.id)) {
              uniqueNaborIngredients.push(ingredient);
            }
          }

          modifier.nabor_ingredient = uniqueNaborIngredients;
          modifiers.push(modifier);
        }

        return modifiers;
      }
    }
    return [];
  };

  const getNewCheck = (product: Product) => {
    let added = false;
    const newCheck = check?.map((checkItem) => {
      if (product.type === "techCart") {
        if (
          checkItem.type === "techCart" &&
          checkItem.id === product.id &&
          checkItem.name === product.name &&
          compareModificators(checkItem, product)
        ) {
          checkItem.count += 1;
          checkItem.total += product.price;
          added = true;
        }
        return checkItem;
      } else {
        if (checkItem.id === product.id && checkItem.type === product.type) {
          checkItem.count += 1;
          checkItem.total += product.price;
          added = true;
        }
        return checkItem;
      }
    });
    return { newCheck, added };
  };

  // add item to check
  const addCheck = useCallback(
    (product: any) => {
      let { newCheck, added } = getNewCheck(product);
      newCheck = added
        ? newCheck
        : [
            ...newCheck,
            {
              id: product.id,
              name: product.name,
              count: 1,
              hasDiscount: product.discount,
              price: product.price,
              total: product.price,
              type: product.type,
              comments: "",
              selectedModificators:
                product.selectedModificators?.map((item: any) => ({
                  ...item,
                  brutto: item.totalBrutto,
                })) || [],
            },
          ];
      setCheck(newCheck);
    },
    [check],
  );

  // +1 item
  const addItem = useCallback(
    (product: any) => () => {
      const { newCheck } = getNewCheck(product);
      setCheck(newCheck);
    },
    [check],
  );

  // -1 item
  const removeItem = useCallback(
    (product: any) => () => {
      const newCheck = check?.map((checkItem) => {
        if (product.type === "techCart") {
          if (
            checkItem.type === "techCart" &&
            checkItem.id === product.id &&
            checkItem.name === product.name &&
            compareModificators(checkItem, product)
          ) {
            checkItem.count -= 1;
            checkItem.total -= product.price;
          }
          return checkItem;
        } else {
          if (checkItem.id === product.id && checkItem.name === product.name) {
            checkItem.count -= 1;
            checkItem.total -= checkItem.price;
          }
          return checkItem;
        }
      });
      setCheck(newCheck.filter((checkItem) => checkItem.count > 0));
    },
    [check],
  );

  const checkIfHasModal = (product: Product) => {
    if (product.nabor) {
      return product.nabor.length > 0;
    }
    return false;
  };
  // handle item selection
  const onItemSelect = (product: Product, route: string) => () => {
    if (product.type === "techCart") {
      const modificators =
        route === "categories"
          ? getModificatorsByProductId(product.id)
          : getModificatorsByMainPanelProductId(product.id);
      if (modificators.length === 0) {
        addCheck({
          ...product,
          selectedModificators: [],
          itemPrice: product.price,
        });
      } else {
        setSelectedProduct(product);
        setSelectedModificators(modificators);
        setModificatorsModalOpen(true);
      }
    } else {
      addCheck({ ...product, itemPrice: product.price });
    }
  };

  const onCheckSubmit = (
    payment: PaymentMethod,
    card: number,
    cash: number,
    pager: number,
  ) => {
    const checkData = prepareCheck(
      CheckStatus.CLOSE,
      payment,
      card,
      cash,
      check,
      userData,
      discount,
      comment,
      orderId,
      pager,
    );
    const idempotency = idempotencyState.idempotencyKey;
    dispatch(regenerateIdempotency());
    return createCheck(checkData, idempotency);
    // const printData = processPrintObject(checkData)
    // return {registeredCheck, printData}
  };

  useEffect(() => {
    if (!isOnline) {
      showAlertMessage(
        "Проверьте ваше интернет подлкючение",
        AlertMessageType.ERROR,
      );
    } else {
      hideAlertMessage();
    }
  }, [isOnline]);

  return (
    <TerminalLayout>
      <div className="w-full flex">
        <div className="w-[38%] py-2 bg-white flex flex-col">
          <div className="pt-4 pb-3 px-6 flex justify-between font-bold text-sm border-b border-gray-200">
            <div className="">Наименование</div>
            <div className="flex">
              <div className="w-20 text-right">Кол-во</div>
              <div className="w-20 text-right">Цена</div>
              <div className="w-20 text-right">Итого</div>
            </div>
          </div>
          <div className="grow flex flex-col divide-y divide-gray-200 overflow-y-auto">
            {check?.map((checkItem, idx) => (
              <div
                key={checkItem.name + idx}
                className="px-6 py-3 flex justify-between text-sm"
              >
                <button
                  onClick={() => {
                    setCheckItemId(idx);
                    setCommentModalOpen(true);
                  }}
                  className="flex flex-col hover:bg-slate-100"
                >
                  <div className="text-base">{checkItem.name}</div>
                  <div className="w-full text-left max-w-[250px] text-xs text-gray-600 overflow-hidden whitespace-nowrap truncate">
                    {checkItem.selectedModificators &&
                      checkItem.selectedModificators
                        ?.map((item) => item.name)
                        .join(", ")}
                  </div>
                  <div className="w-full text-left max-w-[250px] text-xs text-gray-600 overflow-hidden whitespace-nowrap truncate">
                    {checkItem.comments}
                  </div>
                </button>
                <div className="flex">
                  <div className="w-20 relative flex items-center justify-end space-x-2">
                    <button
                      onClick={removeItem(checkItem)}
                      className="w-5 h-5 pb-0.5 border border-gray-400 text-lg flex items-center justify-center rounded-full font-bold bg-zinc-200"
                    >
                      -
                    </button>
                    <span>{checkItem.count}</span>
                    <button
                      onClick={addItem(checkItem)}
                      className="w-5 h-5 pb-0.5 border border-gray-400 text-lg flex items-center justify-center rounded-full font-bold bg-zinc-200"
                    >
                      +
                    </button>
                  </div>
                  <div className="w-20 text-right flex items-center justify-end">
                    {checkItem.price.toFixed(2)}
                  </div>
                  <div className="w-20 text-right flex items-center justify-end">
                    {checkItem.total.toFixed(2)}
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="pt-4 pb-2 px-6 flex flex-col space-y-4 border-t border-gray-200">
            <div className="flex justify-between font-bold">
              <div className="">К оплате</div>
              <div className="flex flex-col space-y-2">
                <div
                  className={clsx([
                    totalPriceWithDiscount < totalPrice
                      ? "TotalPrice text-sm line-through"
                      : "TotalPrice text-lg",
                  ])}
                >
                  {totalPrice.toFixed(2)} ₸
                </div>
                {totalPriceWithDiscount < totalPrice && (
                  <div className="TotalPriceWithDiscount text-lg">
                    {totalPriceWithDiscount.toFixed(2)} ₸
                  </div>
                )}
              </div>
            </div>
            {userData && (
              <div className="flex justify-between font-bold">
                <div className="flex flex-col space-y-1">
                  Текущий клиент: {userData.name} (скидка:{" "}
                  {userData.discount * 100}%)
                </div>
                <div className="flex flex-col space-y-2"></div>
                <div className="text-lg">
                  <button
                    className="col-span-2 py-2 px-1 text-sm tablet:text-base border border-gray-300 rounded-md bg-neutral-100 hover:bg-neutral-200 shadow"
                    onClick={() => {
                      setUserData(null);
                      setDiscount(0);
                    }}
                  >
                    Удалить
                  </button>
                </div>
              </div>
            )}
            <div className="flex justify-between items-center">
              <div className="grid grid-cols-6 gap-3">
                <TerminalModalButton
                  setIsModalOpen={setCommentModalOpen}
                  buttonContent="Комментарий"
                />
                <TerminalModalButton
                  setIsModalOpen={setCodeModalOpen}
                  buttonContent="Код"
                />
                {!userData && (
                  <TerminalModalButton
                    setIsModalOpen={setDiscountModalOpen}
                    buttonContent="Скидка"
                  />
                )}
                <button
                  className="col-span-3 text-white pt-1.5 pb-2 px-3 bg-primary text-lg tablet:text-sm font-semibold rounded-md hover:bg-teal-600"
                  onClick={(e) => {
                    e.preventDefault();
                    showAlertMessage(
                      "Функционал сохранения временно недоступен",
                      AlertMessageType.INFO,
                    );
                  }}
                >
                  Сохранить
                </button>
                <button
                  className="col-span-3 text-white pt-1.5 pb-2 px-3 bg-primary text-lg tablet:text-sm font-semibold rounded-md hover:bg-teal-600"
                  onClick={() => {
                    setPaymentModalOpen(true);
                  }}
                >
                  Оплатить
                </button>
              </div>
            </div>
          </div>
        </div>
        <div className="w-[62%] py-2 px-5 bg-[#d8dbe6] flex flex-col overflow-auto">
          <div className="flex space-x-3">
            <button
              onClick={() => {
                setSelectedCategory(null);
              }}
              className={clsx([
                "font-bold pt-5 pb-1",
                selectedCategory ? "text-indigo-500" : "text-gray-600",
              ])}
            >
              Все товары
            </button>
            {selectedCategory && (
              <button
                className={clsx([
                  "text-gray-600 font-bold pt-5 pb-1 flex items-center space-x-3",
                ])}
              >
                <ChevronLeftIcon className="w-4 h-4" />{" "}
                <span>{selectedCategory.category}</span>
              </button>
            )}
          </div>
          <div className="grow flex flex-col py-4">
            {!selectedCategory ? (
              <>
                <CategoryList
                  categories={categories}
                  setSelectedCategory={setSelectedCategory}
                  mainPanelCategoryId={13}
                />
                <MainPanelProductList
                  mainPanelProducts={mainPanelItems}
                  checkIfHasModal={checkIfHasModal}
                  onItemSelect={onItemSelect}
                />
              </>
            ) : (
              <ProductList
                selectedCategory={selectedCategory}
                checkIfHasModal={checkIfHasModal}
                onItemSelect={onItemSelect}
              />
            )}
          </div>
        </div>
      </div>
      <PaymentModal
        isOpen={paymentModalOpen}
        shiftData={shiftData}
        setIsOpen={setPaymentModalOpen}
        total={totalPriceWithDiscount}
        clearStates={clearPaymentModalsStates}
        setPaymentDone={setPaymentDone}
        paymentDone={paymentDone}
        onClose={() => {
          if (paymentDone && totalPrice !== 0) {
            setPaymentModalOpen(false);
            setCheck([]);
            setUserData(null);
            setComment("");
            setDiscount(0);
            setSelectedProduct(undefined);
            setTotalPrice(0);
            setOrderId(null);
            setCheckItemId(null);
            setClearPaymentModalsStates(true);
          }
        }}
        onSubmit={onCheckSubmit}
        discount={userData ? userData.discount : discount / 100}
      />
      <ModificatorOptionsModal
        isOpen={modificatorsModalOpen}
        setIsOpen={setModificatorsModalOpen}
        title={`${selectedProduct?.name} - Модификаторы`}
        data={selectedModificators}
        product={selectedProduct}
        onSubmit={addCheck}
      />
      <CommentCheckModal
        isOpen={commentModalOpen}
        setIsOpen={setCommentModalOpen}
        data={
          checkItemId !== null && checkItemId !== undefined
            ? check[checkItemId]?.comments
            : comment
        }
        setData={(data: string) => {
          if (checkItemId !== null) {
            setCheck(
              check?.map((checkItem, idx) =>
                idx === checkItemId
                  ? { ...checkItem, comments: data }
                  : checkItem,
              ),
            );
          } else setComment(data);
        }}
        onClose={() => {
          setCheckItemId(null);
        }}
      />
      <CodeModal
        isOpen={codeModalOpen}
        setIsOpen={setCodeModalOpen}
        setDiscount={setDiscount}
        setData={(data: UserMobileData) => {
          setUserData(data);
        }}
      />
      <DiscountModal
        shiftData={shiftData}
        isOpen={discountModalOpen}
        setIsOpen={setDiscountModalOpen}
        setData={(data: number) => {
          setDiscount(data);
        }}
        data={discount}
      />
      <InfoModal
        isOpen={infoModalOpen}
        shiftData={shiftData}
        setIsOpen={setInfoModalOpen}
        onClose={() => {
          if (typeof window !== "undefined") {
            localStorage.removeItem("zebra.authed");
          }
        }}
      />
      {alertMessage && (
        <AlertMessage
          message={alertMessage.message}
          type={alertMessage.type}
          onClose={hideAlertMessage}
        />
      )}
    </TerminalLayout>
  );
};

export default OrderPage;
