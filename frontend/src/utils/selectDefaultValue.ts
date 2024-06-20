//функция для того чтобы передвигать нужную опцию селекта в начало для дефолтного значения
export const selectRightDefaultValue = (arr: [], prop: any, val: any) => {
    const index = arr.findIndex((element) => element[prop] === val);
    if (index !== -1) {
      const element = arr.splice(index, 1)[0];
      arr.unshift(element);
    }
    return arr;
}