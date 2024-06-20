import { Popover, Transition } from "@headlessui/react";
import React, {
  Dispatch,
  FC,
  Fragment,
  SetStateAction,
  useCallback,
  useState,
} from "react";
import clsx from "clsx";
import { uploadResource } from "@api/upload";

export interface ImageUploadProps {
  uploadedImage: string | null;
  setUploadedImage: Dispatch<SetStateAction<string | null>>;
}

const ImageUpload: FC<ImageUploadProps> = ({
  uploadedImage,
  setUploadedImage,
}) => {
  const [image, setImage] = useState<File | null>(null);

  const onUpload = useCallback(
    (close: any) => (e: any) => {
      e.preventDefault();
      e.stopPropagation();
      if (!image) return;
      const formData = new FormData();
      formData.append("image", image);
      uploadResource(formData).then((res) => {
        setUploadedImage(res.data);
        close();
      });
    },
    [image, setUploadedImage]
  );

  return (
    <div className="">
      <Popover className="relative">
        {({ open }) => (
          <>
            <Popover.Button
              className={clsx([
                "h-20 w-28 rounded-lg hover:opacity-80 outline-none",
                uploadedImage
                  ? `bg-[url('https://zebra-crm.kz:8029/itemImage/${uploadedImage}')]`
                  : "bg-zinc-400",
              ])}
              style={
                uploadedImage
                  ? {
                      background: `url('https://zebra-crm.kz:8029/itemImage/${uploadedImage}') no-repeat center`,
                    }
                  : {}
              }
            ></Popover.Button>
            <Transition
              as={Fragment}
              enter="transition ease-out duration-200"
              enterFrom="opacity-0 translate-y-1"
              enterTo="opacity-100 translate-y-0"
              leave="transition ease-in duration-150"
              leaveFrom="opacity-100 translate-y-0"
              leaveTo="opacity-0 translate-y-1"
            >
              <Popover.Panel className="absolute left-0 z-10 mt-0.5 w-68 transform shadow-2xl">
                {({ close }) => (
                  <div className="w-full flex flex-col border border-slate-400 rounded bg-white">
                    <div className="w-full text-sm border-b border-b-slate-400 flex items-center justify-between px-2 py-1">
                      <div className="font-semibold">Загрузка</div>
                      <button
                        onClick={() => setUploadedImage(null)}
                        disabled={!uploadedImage}
                        className="text-indigo-500 hover:text-indigo-700 disabled:text-gray-600 disabled:cursor-not-allowed cursor-pointer"
                      >
                        Удалить
                      </button>
                    </div>
                    <div className="w-full px-2 py-4 border-b border-b-slate-400">
                      <input
                        onChange={(e) => {
                          e.target.files && setImage(e.target.files[0]);
                        }}
                        type="file"
                        name="file"
                        id="image"
                      />
                    </div>
                    <div className="w-full p-2 flex items-center">
                      <button
                        onClick={onUpload(close)}
                        className="bg-primary py-1 px-2 text-white text-sm rounded"
                      >
                        Сохранить
                      </button>
                    </div>
                  </div>
                )}
              </Popover.Panel>
            </Transition>
          </>
        )}
      </Popover>
    </div>
  );
};

export default ImageUpload;
