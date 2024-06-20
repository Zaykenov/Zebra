import React, {
  Context,
  createContext,
  FC,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter } from "next/router";
import { QueryOptions, QueryOptionsData } from "@api/index";
import { defaultEndDate, defaultStartDate } from "@shared/constants/dates";
import { SortOrder } from "@shared/types/types.";

interface FilterContextProps {
  curPage?: number;
  totalPages?: number;
  dateRange?: {
    startDate: Date | null;
    endDate: Date | null;
  };
  searchValue?: string;
  queryOptions: QueryOptionsData;
  toNextPage: () => void;
  toPreviousPage: () => void;
  toPage: (value: number) => void;
  changeTotalPages: (value: number) => void;
  handleDateChange: (dateRange: {
    startDate: Date | null;
    endDate: Date | null;
  }) => void;
  handleSearch: (value: string) => void;
  handleFilterChange: (data: QueryOptionsData) => void;
  handleSort: (field: string) => void;
  getFilterValue: (
    filterOption: QueryOptions
  ) => string | number | number[] | undefined;
  totalResults: number;
  changeTotalResults: (total: number) => void;
}

// @ts-ignore
export const FilterContext: Context<FilterContextProps> = createContext({});
const { Provider } = FilterContext;

const FilterProvider: FC<{
  children: ReactNode;
  defaultFilters?: {
    hasPagination?: boolean;
    hasDatePicker?: boolean;
  };
}> = ({ children, defaultFilters }) => {
  const router = useRouter();

  const [queryOptions, setQueryOptions] = useState<QueryOptionsData>({});

  const [totalPages, setTotalPages] = useState<number>(1);
  const [totalResults, setTotalResults] = useState<number>(0);

  useEffect(() => {
    if (
      !defaultFilters ||
      !router.isReady ||
      Object.keys(router.query).length > 0
    )
      return;
    const { hasPagination, hasDatePicker } = defaultFilters;
    router.query = {
      ...(hasPagination && { page: "1" }),
      ...(hasDatePicker && {
        from: defaultStartDate.toISOString().substring(0, 10),
        to: defaultEndDate.toISOString().substring(0, 10),
      }),
    };
    router.push(router).then();
  }, [defaultFilters, router]);

  useEffect(() => {
    const optionsData: QueryOptionsData = {};
    router.query[QueryOptions.FROM] &&
      (optionsData[QueryOptions.FROM] = router.query[
        QueryOptions.FROM
      ] as string);
    router.query[QueryOptions.TO] &&
      (optionsData[QueryOptions.TO] = router.query[QueryOptions.TO] as string);
    router.query[QueryOptions.SORT] &&
      (optionsData[QueryOptions.SORT] = router.query[
        QueryOptions.SORT
      ] as string);
    router.query[QueryOptions.SEARCH] &&
      (optionsData[QueryOptions.SEARCH] = router.query[
        QueryOptions.SEARCH
      ] as string);
    router.query[QueryOptions.PAGE] &&
      (optionsData[QueryOptions.PAGE] = parseInt(
        router.query[QueryOptions.PAGE] as string
      ));
    router.query[QueryOptions.CATEGORY] &&
      (optionsData[QueryOptions.CATEGORY] = parseInt(
        router.query[QueryOptions.CATEGORY] as string
      ));
    router.query[QueryOptions.WORKER] &&
      (optionsData[QueryOptions.WORKER] = parseInt(
        router.query[QueryOptions.WORKER] as string
      ));
    router.query[QueryOptions.SKLAD] &&
      (optionsData[QueryOptions.SKLAD] = Array.isArray(
        router.query[QueryOptions.SKLAD]
      )
        ? (router.query[QueryOptions.SKLAD] as string[]).map((elem) =>
            parseInt(elem)
          )
        : parseInt(router.query[QueryOptions.SKLAD] as string));
    router.query[QueryOptions.SCHET] &&
      (optionsData[QueryOptions.SCHET] = parseInt(
        router.query[QueryOptions.SCHET] as string
      ));
    router.query[QueryOptions.DEALER] &&
      (optionsData[QueryOptions.DEALER] = parseInt(
        router.query[QueryOptions.DEALER] as string
      ));
    router.query[QueryOptions.PAYMENT] &&
      (optionsData[QueryOptions.PAYMENT] = router.query[
        QueryOptions.PAYMENT
      ] as string);
    router.query[QueryOptions.TYPE] &&
      (optionsData[QueryOptions.TYPE] = router.query[
        QueryOptions.TYPE
      ] as string);
    router.query[QueryOptions.STATUS] &&
      (optionsData[QueryOptions.STATUS] = router.query[
        QueryOptions.STATUS
      ] as string);
    router.query[QueryOptions.SHOP] &&
      (optionsData[QueryOptions.SHOP] = Array.isArray(
        router.query[QueryOptions.SHOP]
      )
        ? (router.query[QueryOptions.SHOP] as string[]).map((elem) =>
            parseInt(elem)
          )
        : parseInt(router.query[QueryOptions.SHOP] as string));
    setQueryOptions(optionsData);
  }, [router.query]);

  const toNextPage = useCallback(async () => {
    const pageStr = router.query.page;
    if (!pageStr) return;
    const page = parseInt(pageStr as string);
    if (page === totalPages) return;
    router.query = {
      ...router.query,
      page: `${page + 1}`,
    };
    await router.push(router);
  }, [router.query, totalPages]);

  const toPreviousPage = useCallback(async () => {
    const pageStr = router.query.page;
    if (!pageStr) return;
    const page = parseInt(pageStr as string);
    if (page === 1) return;
    router.query = {
      ...router.query,
      page: `${page - 1}`,
    };
    await router.push(router);
  }, [router.query, totalPages]);

  const toPage = useCallback(
    async (page: number) => {
      if (!Object.keys(router.query).length) return;
      router.query = {
        ...router.query,
        page: `${page}`,
      };
      await router.push(router);
    },
    [router.query]
  );

  const changeTotalPages = useCallback(
    (value: number) => setTotalPages(value),
    []
  );

  const handleDateChange = useCallback(
    async ({
      startDate,
      endDate,
    }: {
      startDate: Date | null;
      endDate: Date | null;
    }) => {
      router.query = {
        ...router.query,
        from: startDate?.toISOString().substring(0, 10),
        to: endDate?.toISOString().substring(0, 10),
      };

      await router.push(router);
    },
    [router.query]
  );

  const handleSearch = useCallback(
    async (value: string) => {
      router.query = {
        ...router.query,
        search: value,
      };
      if (value.trim().length === 0) delete router.query?.search;
      await router.push(router);
    },
    [router.query]
  );

  const handleFilterChange = useCallback(
    async (data: QueryOptionsData) => {
      if (!Object.keys(router.query).length) return;
      // @ts-ignore
      router.query = {
        ...router.query,
        ...data,
      };
      await router.push(router);
    },
    [router.query]
  );

  const handleSort = useCallback(
    async (field: string) => {
      if (!router.isReady) return;
      const sortInfo = router.query.sort;
      let order: SortOrder;

      if (!sortInfo) order = SortOrder.ASC;
      else {
        const sortInfoArray = (sortInfo as string).split(" ");
        const oldField = sortInfoArray[0];
        const oldOrder = sortInfoArray[1];
        if (oldField === field) {
          order = oldOrder === SortOrder.ASC ? SortOrder.DESC : SortOrder.ASC;
        } else order = SortOrder.ASC;
      }

      router.query = {
        ...router.query,
        [QueryOptions.SORT]: `${field} ${order}`,
      };
      await router.push(router);
    },
    [router.query]
  );

  const changeTotalResults = useCallback((total: number) => {
    setTotalResults(total);
  }, []);

  const curPage = useMemo(() => queryOptions.page, [queryOptions.page]);
  const dateRange = useMemo(
    () => ({
      startDate: queryOptions.from ? new Date(queryOptions.from) : null,
      endDate: queryOptions.to ? new Date(queryOptions.to) : null,
    }),
    [queryOptions.from, queryOptions.to]
  );
  const searchValue = useMemo(() => queryOptions.search, [queryOptions.search]);

  const getFilterValue = useCallback(
    (filterOption: QueryOptions) => queryOptions[filterOption],
    [queryOptions]
  );

  return (
    <Provider
      value={{
        curPage,
        totalPages,
        dateRange,
        searchValue,
        queryOptions,
        toNextPage,
        toPreviousPage,
        toPage,
        changeTotalPages,
        handleDateChange,
        handleSearch,
        handleFilterChange,
        handleSort,
        getFilterValue,
        totalResults,
        changeTotalResults,
      }}
    >
      {children}
    </Provider>
  );
};

export const useFilter = (): FilterContextProps => {
  return useContext(FilterContext);
};

export default FilterProvider;
