export const clearStorage = () => {
  if (typeof window !== "undefined") {
    for (const storageItemKey in localStorage) {
      if (storageItemKey.startsWith("zebra.categories") || storageItemKey.startsWith("zebra.mainPanel")) {
        localStorage.removeItem(storageItemKey);
      }
    }
  }
};