import React, { Dispatch, FC, SetStateAction, useCallback } from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import ModifierForm from "../ModifierForm";
import { ModifierData } from "../ModifierForm/types";

export interface ModifierFormModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  data?: ModifierData | null;
  addModifierSet: (data: ModifierData) => void;
  updateModifierSet: (data: ModifierData) => void;
  removeModifierSet: (name: string) => void;
}

const ModifierFormModal: FC<ModifierFormModalProps> = ({
  isOpen,
  setIsOpen,
  data,
  addModifierSet,
  updateModifierSet,
  removeModifierSet,
}) => {
  const handleClose = useCallback(() => {
    setIsOpen(false);
  }, [setIsOpen]);

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Добавление набора модификаторов"
    >
      <div className="flex flex-col items-start">
        <ModifierForm
          data={data}
          onClose={handleClose}
          onAdd={addModifierSet}
          onUpdate={updateModifierSet}
          onDelete={removeModifierSet}
        />
      </div>
    </ModalLayout>
  );
};

export default ModifierFormModal;
