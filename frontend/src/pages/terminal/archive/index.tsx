import React, { FC, useEffect, useMemo, useState } from "react";
import TerminalLayout from "@layouts/TerminalLayout";
import TerminalCabinetLayout from "@layouts/TerminalCabinetLayout";
import CheckArchive from "@modules/CheckArchive";
import { getAllChecks, getAllWorkerChecks } from "@api/check";

const ArchivePage: FC = () => {
  const [data, setData] = useState();

  useEffect(() => {
    getAllWorkerChecks().then((res) => {
      setData(res.data);
    });
  }, []);

  const navigation = useMemo(
    () => [
      {
        id: 0,
        name: "Все чеки",
        component: data && <CheckArchive data={data} />,
      },
    ],
    [data]
  );

  return (
    <TerminalLayout>
      <TerminalCabinetLayout navigation={navigation} />
    </TerminalLayout>
  );
};

export default ArchivePage;
