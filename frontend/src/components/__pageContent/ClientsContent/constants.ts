import { Column } from "react-table";
import { ClientData } from "./types";

export const columns: Column[] = [
  {
    Header: "Имя",
    accessor: "firstName",
  },
  {
    Header: "День рождения",
    accessor: "birthDate",
  },
  {
    Header: "Электронная почта",
    accessor: "email",
  },
  {
    Header: "Дата регистрации",
    accessor: "registrationDate",
  },
  {
    Header: "Скидка",
    accessor: "discount",
  },
  {
    Header: "Баланс ZebraCoin",
    accessor: "zebraCoinBalance",
  },
];

export const filterProperties = (
  arr: any[],
  properties?: (keyof ClientData)[]
): Partial<ClientData>[] => {
  return arr.map((obj) => {
    const filteredObj: Partial<ClientData> = {};
    if (properties === undefined) {
      properties = Object.keys(obj) as (keyof ClientData)[];
    }
    properties.forEach((prop) => {
      filteredObj[prop] = obj[prop];
    });
    return filteredObj;
  });
};
