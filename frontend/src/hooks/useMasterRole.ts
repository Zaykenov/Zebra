import { useEffect, useState } from "react";

const useMasterRole = () => {
  const [isMaster, setIsMaster] = useState<boolean | null>(null);

  useEffect(() => {
    if (typeof window === "undefined") return;
    const role = localStorage.getItem("zebra.role");
    if (!role) return;
    setIsMaster(role === "master");
  }, []);

  return isMaster;
};

export default useMasterRole;
