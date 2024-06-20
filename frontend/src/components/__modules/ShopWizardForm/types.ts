import { WorkerRole } from "@api/workers";

export type TovarData = {
  id: number;
  name: string;
  price: number;
  category: string;
};

export type ExistingWorkerData = {
  id: number;
  name: string;
  username: string;
  password: string;
  phone: string;
  role: WorkerRole;
};

export type WorkerOption = {
  label: string;
  value: number;
  data: ExistingWorkerData;
};
