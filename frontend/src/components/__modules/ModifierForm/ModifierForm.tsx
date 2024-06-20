import React, { FC, useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { ModifierData } from "./types";
import { Input } from "@shared/ui/Input";
import { Select } from "@shared/ui/Select";
import {
  getMasterModifierNabor,
  getMasterModifierNabors,
} from "@api/modifiers";

export interface ModifierFormProps {
  onClose: () => void;
  data?: ModifierData | null;
  onAdd: (data: ModifierData) => void;
  onUpdate: (data: ModifierData) => void;
  onDelete: (name: string) => void;
}

const ModifierForm: FC<ModifierFormProps> = ({
  onClose,
  data,
  onAdd,
  onUpdate,
  onDelete,
}) => {
  const [isNew, setIsNew] = useState(false);
  const [naborOptions, setNaborOptions] = useState([]);
  const [selectedNabor, setSelectedNabor] = useState<number>(1);

  useEffect(() => {
    getMasterModifierNabors().then((res) => {
      setNaborOptions(
        res.data.map(({ name, id }: { name: string; id: number }) => ({
          name,
          value: id,
        }))
      );
      setSelectedNabor(res.data[0]?.id);
    });
  }, []);

  const {
    register,
    handleSubmit,
    watch,
    formState: { isDirty },
  } = useForm<ModifierData>({
    defaultValues: {
      name: data?.name || "",
      min: 0,
      max: 0,
      ingredient_nabor: [
        {
          ingredient_id: 1,
          brutto: 0,
          price: 0,
        },
      ],
    },
  });

  const onSubmit = useCallback(
    (submitData: ModifierData) => {
      if (isNew) data ? onUpdate(submitData) : onAdd(submitData);
      else {
        selectedNabor &&
          getMasterModifierNabor(selectedNabor).then((res) => {
            console.log(res.data);
            onAdd({
              ...res.data,
              ingredient_nabor: res.data.ingredient_nabor.map((ingr: any) => ({
                ...ingr,
                value: ingr.ingredient_id,
                measure: ingr.measure,
              })),
            });
          });
      }
      onClose();
    },
    [onAdd, onClose, onUpdate, data, isNew, selectedNabor]
  );

  return (
    <div className="flex flex-col items-start w-full">
      <button
        type="button"
        onClick={() => setIsNew((prevState) => !prevState)}
        className="py-2 px-3 rounded-md text-sm border border-primary text-primary hover:bg-primary hover:text-white font-medium mb-4"
      >
        {isNew ? "Выбрать существующий" : "Создать новый"}
      </button>
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="w-full flex flex-col divide-y divide-gray-100"
      >
        {isNew ? (
          <>
            <Input
              {...register("name")}
              type="text"
              placeholder="Название набора"
              className="mb-4"
            />
            <div className="w-full flex flex-col py-4">
              <div className="text-sm mb-4">
                Сколько модификаторов можно выбрать одновременно:
              </div>
              <div className="space-y-2.5">
                <div key="single" className="flex items-center cursor-pointer">
                  <input
                    id="single"
                    name="method"
                    type="radio"
                    defaultChecked={true}
                    className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <label
                    htmlFor="single"
                    className="ml-3 block text-sm font-medium text-gray-700 cursor-pointer"
                  >
                    Только один
                  </label>
                </div>
                <div
                  key="multiple"
                  className="flex items-center cursor-pointer"
                >
                  <input
                    id="multiple"
                    name="method"
                    type="radio"
                    defaultChecked={true}
                    className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <label
                    htmlFor="multiple"
                    className="ml-3 block text-sm font-medium text-gray-700 cursor-pointer"
                  >
                    Несколько
                  </label>
                </div>
              </div>
            </div>
          </>
        ) : (
          <>
            <Select
              options={naborOptions}
              onChange={(e) => {
                setSelectedNabor(
                  parseInt((e.target as HTMLSelectElement).value)
                );
              }}
            />
          </>
        )}
        <div className="w-full flex items-center justify-between pt-4">
          <div className="flex items-center space-x-3">
            <button
              type="submit"
              // onClick={handleSubmit(onSubmit)}
              disabled={isNew && !isDirty}
              className="px-3 pb-1 pt-0.5 bg-primary text-white disabled:opacity-80 disabled:cursor-not-allowed hover:opacity-80 rounded"
            >
              Добавить
            </button>
            <button
              onClick={(e) => {
                e.preventDefault();
                onClose();
              }}
              className="px-3 pb-1 pt-0.5 border border-gray-400 rounded"
            >
              Отменить
            </button>
          </div>
          {data && (
            <button
              onClick={(e) => {
                e.preventDefault();
                onDelete(watch("name"));
                onClose();
              }}
              className="px-3 pb-1 pt-0.5 border border-red-500 text-red-500 hover:bg-red-500 hover:text-white rounded"
            >
              Удалить набор из тех. карты
            </button>
          )}
        </div>
      </form>
    </div>
  );
};

export default ModifierForm;
