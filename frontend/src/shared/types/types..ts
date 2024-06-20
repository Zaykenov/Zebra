export enum MeasureOption {
  DEFAULT = "шт.",
  KILOGRAMS = "кг.",
  LITERS = "л.",
}

export const measureOptions = [
  { name: "шт.", value: MeasureOption.DEFAULT },
  { name: "кг.", value: MeasureOption.KILOGRAMS },
  { name: "л.", value: MeasureOption.LITERS },
];

export enum SortOrder {
  ASC = "asc",
  DESC = "desc",
}
