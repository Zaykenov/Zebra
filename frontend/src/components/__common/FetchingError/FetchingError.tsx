import Image from "next/image";
import { FC } from "react";

interface FetchingErrorProps {
  errorMessage: string;
}

const FetchingError: FC<FetchingErrorProps> = ({ errorMessage }) => {
  return (
    <div className="mt-20 flex flex-col items-center">
      <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
        Извините, произошла ошибка при получении данных:
      </p>
      <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
        <strong>"{errorMessage}"</strong>
      </p>
      <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
        Скорее всего, это связано с <strong>перезагрузкой сервера</strong>.
      </p>
      {/* <Image src={"/images/jdun-meme.gif"} height={200} width={200} /> */}
      <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
        Пожалуйста, подождите <strong>3-5 минут</strong>.
      </p>
      <p className="mb-4 text-gray-500 dark:text-gray-300 text-xl">
        Eсли проблема не решится, свяжитесь с разработчиками.
      </p>
    </div>
  );
};

export default FetchingError;
