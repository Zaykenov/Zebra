export const objectToQueryParamsString = (data: any) => {
  return (
    "?" +
    Object.keys(data)
      .map((key) => {
        const value = data[key];
        return Array.isArray(value)
          ? value.map((elem) => `${key}=${elem}`).join("&")
          : `${key}=${value}`;
      })
      .join("&")
  );
};
