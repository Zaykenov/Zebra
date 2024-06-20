export const formatNumber = (
  num: number,
  isCurrency: boolean = false,
  isToFixed: boolean = true
) => {
  return (
    num.toLocaleString("ru-RU", {
      style: "decimal",
      ...(isToFixed && { maximumFractionDigits: 2, minimumFractionDigits: 2 }),
    }) + (isCurrency ? " ₸" : "")
  );
};
