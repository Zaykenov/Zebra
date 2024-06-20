const filterProperties = <T extends Record<string, any>>(
  arr: T[],
  properties?: (keyof T)[]
): Partial<T>[] => {
  return arr.map((obj) => {
    const filteredObj: Partial<T> = {};
    if (properties === undefined) {
      properties = Object.keys(obj) as (keyof T)[];
    }
    properties.forEach((prop) => {
      filteredObj[prop] = obj[prop];
    });
    return filteredObj;
  });
};

export default filterProperties