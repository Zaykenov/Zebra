import React, {
  ChangeEvent,
  FC,
  HTMLProps,
  useCallback,
  useEffect,
  useMemo,
} from "react";
import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import debounce from "lodash.debounce";
import { useFilter } from "@context/index";

interface SearchProps extends HTMLProps<HTMLInputElement> {}

const Search: FC<SearchProps> = ({ onChange }) => {
  const { handleSearch } = useFilter();

  useEffect(() => {
    return () => {
      debouncedResults.cancel();
    };
  }, []);

  const handleChange = useCallback(
    async (e: ChangeEvent<HTMLInputElement>) => {
      await handleSearch(e.target.value);
    },
    [handleSearch],
  );

  const debouncedResults = useMemo(() => {
    return debounce(handleChange, 300);
  }, [handleChange]);

  return (
    <div className="flex items-center">
      <div className="relative rounded-md shadow-sm border border-gray-300">
        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
          <MagnifyingGlassIcon
            className="h-5 w-5 text-gray-400"
            aria-hidden="true"
          />
        </div>
        <input
          type="text"
          name="text"
          id="email"
          onChange={onChange ? onChange : debouncedResults}
          className="block w-full rounded-md border-gray-300 pl-10 py-2 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          placeholder="Быстрый поиск"
        />
      </div>
    </div>
  );
};

export default Search;
