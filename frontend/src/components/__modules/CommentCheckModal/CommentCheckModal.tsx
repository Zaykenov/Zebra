import React, {
  Dispatch,
  FC,
  SetStateAction,
  useCallback,
  useEffect,
  useState,
} from "react";
import ModalLayout from "@layouts/ModalLayout/ModalLayout";
import {
  createTag,
  deleteTag,
  getAllTags,
  TagData,
} from "@api/tags";
import { XMarkIcon } from "@heroicons/react/24/solid";
import clsx from "clsx";
import { Input } from "@shared/ui/Input";

export interface CommentCheckModalProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  setData: (data: string) => void;
  onClose: () => void;
  data: string;
}

const CommentCheckModal: FC<CommentCheckModalProps> = ({
  isOpen,
  setIsOpen,
  setData,
  onClose,
  data,
}) => {
  const [tags, setTags] = useState<TagData[]>([]);
  const [text, setText] = useState<string>("");
  const [comment, setComment] = useState<string | null>("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  useEffect(() => {
    setComment(data);
    data && data.length > 0
      ? setSelectedTags(data.split(", "))
      : setSelectedTags([]);
  }, [data]);

  useEffect(() => {
    setComment(selectedTags.join(", "));
  }, [selectedTags]);

  const retrieveTags = useCallback(() => {
    getAllTags().then((res) => {
      setTags(res.data);
    });
  }, []);

  useEffect(() => {
    retrieveTags();
  }, [retrieveTags]);

  const addTag = useCallback(
    (text: string) => {
      createTag({ text }).then(() => {
        retrieveTags();
      });
    },
    [retrieveTags]
  );

  const removeTag = useCallback(
    ({ id, text }: TagData) => {
      if (!id) return;
      setSelectedTags(selectedTags.filter((tagText) => tagText !== text));
      deleteTag({ id }).then(() => {
        retrieveTags();
      });
    },
    [retrieveTags, selectedTags]
  );

  return (
    <ModalLayout
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onClose={onClose}
      title="Добавить комментарий к заказу"
    >
      <form className="flex flex-col">
        <ul className="w-full flex flex-wrap gap-2 mb-4">
          {tags.map((tag) => (
            <li
              key={tag.id}
              className={clsx([
                "pt-0.5 pb-1 pl-3 pr-2 flex items-center space-x-2 cursor-pointer rounded-3xl text-white",
                selectedTags.includes(tag.text)
                  ? "bg-primary ring ring-teal-600"
                  : "bg-primary/75",
              ])}
              onClick={() => {
                selectedTags.includes(tag.text)
                  ? setSelectedTags(
                      selectedTags.filter((tagText) => tagText !== tag.text)
                    )
                  : setSelectedTags((prevState) => [...prevState, tag.text]);
              }}
            >
              <span className="text-sm font-medium">{tag.text}</span>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  removeTag(tag);
                }}
                className="flex items-center justify-center mt-0.5 p-0.5 rounded-md hover:bg-gray-300/40"
              >
                <XMarkIcon className="w-4 h-4" />
              </button>
            </li>
          ))}
        </ul>
        <div className="w-full flex space-x-3 mb-2">
          <Input
            className=""
            type="text"
            onInput={(e) => {
              setText((e.target as HTMLInputElement).value);
            }}
            name="text"
          />
          <button
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              addTag(text);
            }}
            className="pb-1 pt-0.5 px-3 whitespace-nowrap bg-primary text-white font-medium rounded-md"
          >
            Добавить фактор
          </button>
        </div>
        <div className="w-full flex">
          <label></label>
          <textarea
            placeholder="Комментарий"
            value={comment as string}
            onInput={(e) => {
              setComment((e.target as HTMLTextAreaElement).value);
            }}
            className="w-full p-2 text-sm h-20 border border-gray-300 outline-none"
          />
        </div>
        <button
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            setData(comment || "");
            setSelectedTags([]);
            setComment("");
            onClose();
            setIsOpen(false);
          }}
          className="mt-3 pt-1 pb-1.5 text-sm flex items-center justify-center bg-primary hover:opacity-80 text-white rounded"
        >
          Сохранить
        </button>
      </form>
    </ModalLayout>
  );
};

export default CommentCheckModal;
