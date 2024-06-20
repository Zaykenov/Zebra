import { Stats } from "./types";

export const statsMapping = {
  [Stats.REVENUE]: "Выручка",
  [Stats.PROFIT]: "Прибыль",
  [Stats.CHECKS]: "Чеки",
  [Stats.VISITORS]: "Посетители",
  [Stats.AVG_CHECK]: "Средний чек",
};

export const weekDaysLabels = {
  monday: "Пн",
  tuesday: "Вт",
  wednesday: "Ср",
  thursday: "Чт",
  friday: "Пт",
  saturday: "Сб",
  sunday: "Вс",
};

export const btns = [
  {
    text: "выручка",
    type: "currency",
    value: Stats.REVENUE,
    isCurrency: true,
  },
  {
    text: "прибыль",
    type: "currency",
    value: Stats.PROFIT,
    isCurrency: true,
  },
  {
    text: "чеки",
    type: "number",
    value: Stats.CHECKS,
    isCurrency: false,
  },
  {
    text: "посетители",
    type: "number",
    value: Stats.VISITORS,
    isCurrency: false,
  },
  {
    text: "средний чек",
    type: "currency",
    value: Stats.AVG_CHECK,
    isCurrency: true,
  },
];
