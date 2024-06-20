import React, { useEffect, useState } from "react";
import { NextPage } from "next";
import PageLayout from "@layouts/PageLayout";
import MainLayout from "@layouts/MainLayout";
import NaborForm from "@modules/NaborForm";
import { useRouter } from "next/router";
import { getMasterModifierNabor, getModifier } from "@api/modifiers";
import useMasterRole from "@hooks/useMasterRole";

const NaborFormPage: NextPage = () => {
  const router = useRouter();

  const isMaster = useMasterRole();

  const [data, setData] = useState(null);

  useEffect(() => {
    if (isMaster === null) return;
    const id = router.query.id;
    if (!id) return;
    const naborId = parseInt(id as string);
    (async function () {
      const res = isMaster
        ? await getMasterModifierNabor(naborId)
        : await getModifier(naborId);
      setData(res.data);
    })();
  }, [router, isMaster]);

  return (
    <PageLayout>
      <MainLayout title="Добавление набора" backBtn={true}>
        <div className="p-5">
          <NaborForm data={data} isEdit />
        </div>
      </MainLayout>
    </PageLayout>
  );
};

export default NaborFormPage;
