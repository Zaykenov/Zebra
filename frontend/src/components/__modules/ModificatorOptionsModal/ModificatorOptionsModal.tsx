import React, {
  Dispatch,
  FC,
  SetStateAction,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import ItemCard from "@common/ItemCard/ItemCard";
import clsx from "clsx";

type IngredientData = {
  id: number;
  nabor_id: number;
  name: string;
  category: number;
  measure: string;
  cost: number;
  brutto: number;
  price: number;
  image: string;
};

export interface ModificatorFormModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  title: string;
  product: any;
  onSubmit: (product: any) => void;
  data: {
    id: number;
    name: string;
    max: number;
    min: number;
    nabor_ingredient: {
      id: number;
      nabor_id: number;
      name: string;
      category: number;
      measure: string;
      cost: number;
      brutto: number;
      price: number;
      image: string;
    }[];
  }[];
}

export type SelectedModificator = {
  id: number;
  nabor_id: number;
  name: string;
  quantity: number;
  totalBrutto: number;
  brutto: number;
  cost: number;
};

const ModificatorOptionsModal: FC<ModificatorFormModalProps> = ({
  isOpen,
  setIsOpen,
  title,
  product,
  data,
  onSubmit,
}) => {
  const [selectedModificators, setSelectedModificators] = useState<
    SelectedModificator[]
  >([]);

  const [totalPrice, setTotalPrice] = useState(0);

  const [naborsCount, setNaborsCount] = useState<
    {
      id: number;
      min: number;
      max: number;
      count: number;
    }[]
  >([]);

  const [submitDisabled, setSubmitDisabled] = useState(false);

  useEffect(() => {
    setNaborsCount(
      data.map((nabor) => ({
        id: nabor.id,
        min: nabor.min,
        max: nabor.max,
        count: 0,
      })),
    );
  }, [data]);

  const reachedMax = useCallback(
    (naborId: number) => {
      const naborInfo = naborsCount.find((nabor) => nabor.id === naborId);
      if (!naborInfo) return true;
      return naborInfo.count === naborInfo.max;
    },
    [naborsCount],
  );

  useEffect(() => {
    if (!naborsCount.length) return;
    const disabled = naborsCount.some(
      ({ count, min, max }) => count > 0 && (count < min || count > max),
    );
    setSubmitDisabled(disabled);
  }, [naborsCount]);

  useEffect(() => {
    product && setTotalPrice(product.price);
  }, [product]);

  const isSelected = useCallback(
    (id: number) => {
      for (let i = 0; i < selectedModificators.length; i++) {
        if (selectedModificators[i].id === id) return true;
      }
      return false;
    },
    [selectedModificators],
  );

  const onModificatorSelect = useCallback(
    (ingredient: IngredientData) => (e: any) => {
      e.preventDefault();
      e.stopPropagation();
      if (reachedMax(ingredient.nabor_id)) return;
      setNaborsCount((prevState) =>
        prevState.map((naborInfo) => {
          if (naborInfo.id === ingredient.nabor_id) {
            return { ...naborInfo, count: naborInfo.count + 1 };
          }
          return naborInfo;
        }),
      );
      setSelectedModificators((prevState) => {
        setTotalPrice(totalPrice + ingredient.price);

        if (prevState.find((item) => item.id === ingredient.id)) {
          return prevState.map((item) => {
            if (item.id === ingredient.id) {
              const quantity = item.quantity + 1;
              return {
                ...item,
                quantity,
                totalBrutto: item.brutto * quantity,
              };
            }
            return item;
          });
        }
        return [
          ...prevState,
          {
            id: ingredient.id,
            nabor_id: ingredient.nabor_id,
            name: ingredient.name,
            quantity: 1,
            totalBrutto: ingredient.brutto,
            brutto: ingredient.brutto,
            cost: ingredient.cost,
            price: ingredient.price,
          },
        ];
      });
    },
    [totalPrice, reachedMax],
  );

  const onModificatorRemove = useCallback(
    (ingredient: IngredientData) => (e: any) => {
      e.preventDefault();
      e.stopPropagation();
      // if (reachedMin) return;
      setNaborsCount((prevState) =>
        prevState.map((naborInfo) => {
          if (naborInfo.id === ingredient.nabor_id) {
            return { ...naborInfo, count: naborInfo.count - 1 };
          }
          return naborInfo;
        }),
      );
      setSelectedModificators((prevState) => {
        const removedIngredient = prevState.find(
          (item) => item.id === ingredient.id,
        );
        if (!removedIngredient) {
          return prevState;
        }
        setTotalPrice(totalPrice - ingredient.price);
        if (removedIngredient.quantity === 1) {
          return prevState.filter((item) => item.id !== ingredient.id);
        } else {
          return prevState.map((item) => {
            if (item.id !== ingredient.id) return item;
            const quantity = item.quantity - 1;
            return {
              ...item,
              quantity,
              totalBrutto: item.brutto * quantity,
            };
          });
        }
      });
    },
    [totalPrice],
  );

  const getQuantity = (ingredientId: number) => {
    const modificator = selectedModificators.find(
      (item) => item.id === ingredientId,
    );
    if (!modificator) return undefined;
    return modificator.quantity;
  };

  const onClose = useCallback(() => {
    setTotalPrice(product.price);
    setIsOpen(false);
  }, [setIsOpen, product]);

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onClose={() => {
        setSelectedModificators([]);
      }}
      title={title}
      maxWidth="max-w-3xl"
      fullScreen
      footer={
        <ModalFooter
          totalPrice={totalPrice}
          onClose={onClose}
          onSubmit={onSubmit}
          product={product}
          selectedModificators={selectedModificators}
          setSelectedModificators={setSelectedModificators}
          submitDisabled={submitDisabled}
        />
      }
    >
      <div className="w-full flex flex-col">
        <div className="flex flex-col space-y-3">
          {data.map((modificator) => (
            <div
              key={modificator.id}
              className="flex overflow-x-auto flex-col space-y-2.5"
            >
              <div className="text-gray-500">{modificator.name}</div>
              <div className="flex flex-nowrap space-x-3.5 pb-2">
                {modificator.nabor_ingredient.map((ingredient) => (
                  <div
                    key={`${modificator.id}_${ingredient.id}`}
                    className="w-36"
                  >
                    <ItemCard
                      name={ingredient.name}
                      cover={ingredient.image}
                      price={ingredient.price}
                      className={clsx([
                        "m-1 rounded w-36",
                        isSelected(ingredient.id)
                          ? "ring ring-indigo-400"
                          : "ring-0",
                      ])}
                      height="h-[110px]"
                      onSelect={onModificatorSelect({
                        ...ingredient,
                        nabor_id: modificator.id,
                      })}
                      quantity={getQuantity(ingredient.id)}
                      onRemove={onModificatorRemove({
                        ...ingredient,
                        nabor_id: modificator.id,
                      })}
                    />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </ModalLayout>
  );
};

const ModalFooter: FC<{
  totalPrice: number;
  onClose: () => void;
  onSubmit: (product: any) => void;
  product: any;
  selectedModificators: SelectedModificator[];
  setSelectedModificators: Dispatch<SetStateAction<SelectedModificator[]>>;
  submitDisabled: boolean;
}> = ({
  totalPrice,
  onClose,
  onSubmit,
  product,
  selectedModificators,
  setSelectedModificators,
  submitDisabled,
}) => {
  return (
    <div className="flex items-center justify-between px-6 py-4 rounded-md">
      <div>
        Итого: <span className="font-bold">{totalPrice.toFixed(2)} ₸</span>
      </div>
      <div className="flex items-center space-x-2">
        <button
          onClick={onClose}
          className="px-4 py-2 rounded-md border border-gray-300"
        >
          Отменить
        </button>
        <button
          onClick={() => {
            onSubmit({
              ...product,
              selectedModificators,
              price: totalPrice,
              itemPrice: product.price,
            });
            setSelectedModificators([]);
            onClose();
          }}
          disabled={submitDisabled}
          className="px-4 py-2 rounded-md border border-gray-300 bg-primary hover:opacity-80 text-white font-bold disabled:bg-gray-400"
        >
          Добавить чек
        </button>
      </div>
    </div>
  );
};

export default ModificatorOptionsModal;
