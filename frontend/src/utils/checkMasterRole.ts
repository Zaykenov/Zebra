export const isMaster = () => {
  if (typeof window === "undefined") return false;
  const role = localStorage.getItem("zebra.role");
  if (!role || role !== "master") return false;
  return true;
};
