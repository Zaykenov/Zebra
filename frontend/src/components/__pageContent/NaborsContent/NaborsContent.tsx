import React, { FC, useCallback, useEffect, useState } from "react";
import { FilterOption } from "@layouts/MainLayout/types";
import Table from "@common/Table";
import MainLayout from "@layouts/MainLayout";
import { getAllModifiers, getMasterModifierNabors } from "@api/modifiers";
import SubRow from "./SubRow";
import { columns } from "./constants";
import { useFilter } from "@context/index";
import useMasterRole from "@hooks/useMasterRole";

const NaborsContent: FC = () => {
  const isMaster = useMasterRole();

  const { queryOptions, changeTotalPages, changeTotalResults } = useFilter();

  const [tableData, setTableData] = useState([]);

  const renderSubComponent = useCallback((row: any) => {
    return <SubRow rowData={row.row} />;
  }, []);

  useEffect(() => {
    if (!queryOptions || !Object.keys(queryOptions).length || isMaster === null)
      return;
    (async function () {
      const res = isMaster
        ? await getMasterModifierNabors(queryOptions)
        : await getAllModifiers(queryOptions);

      changeTotalPages(res.data.totalPages);
      setTableData(res.data.data);
      changeTotalResults(res.data.data.length);
    })();
  }, [queryOptions, isMaster]);

  return (
    <MainLayout
      title="Наборы модификаторов"
      addHref={isMaster ? "/nabors/nabor_form" : ""}
      searchFilter
      pagination
      filterOptions={[FilterOption.SHOP]}
    >
      {tableData && (
        <Table
          columns={columns}
          data={tableData}
          details={true}
          renderRowSubComponent={renderSubComponent}
        />
      )}
    </MainLayout>
  );
};

export default NaborsContent;
