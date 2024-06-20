export const cleanCache = () => {
  if (typeof window !== "undefined") {
    for (const storageItemKey in localStorage) {
      if (storageItemKey.startsWith("zebra.cache")) {
        localStorage.removeItem(storageItemKey);
      }
    }
  }
};
