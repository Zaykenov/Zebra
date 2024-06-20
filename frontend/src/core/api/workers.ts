import { backendApiInstance } from "./index";

export enum WorkerRole {
  MANAGER = "manager",
  WORKER = "worker",
  MASTER = "master",
}

export type WorkerData = {
  name: string;
  phone: string;
  username: string;
  password: string;
  role: WorkerRole;
  id?: number;
  shops: string[];
};

const processFormData = (data: WorkerData) => {
  const postData = {
    name: data.name,
    phone: data.phone,
    username: data.username,
    password: data.password,
    role: data.role,
    shops: data.shops.map((shop) => parseInt(shop)),
  };
  return data.id ? { id: data.id, ...postData } : postData;
};

export const getAllWorkers = async () => {
  const res = await backendApiInstance.get("workers/getAll");
  return res.data;
};

export const getWorker = async (id: string) => {
  const res = await backendApiInstance.get(`workers/get/${id}`);
  return res.data;
};

export const createWorker = async (data: any) => {
  const res = await backendApiInstance.post("authorize", processFormData(data));
  return res.data;
};

export const updateWorker = async (data: any) => {
  const res = await backendApiInstance.post(
    "workers/update",
    processFormData(data)
  );
  return res.data;
};

export const deleteWorker = async (data: { id: number }) => {
  const res = await backendApiInstance.post(`workers/delete/${data.id}`);
  return res.data;
};

export const signIn = async (data: { username: string; password: string }) => {
  const res = await backendApiInstance.post("signin", data);
  return res.data;
};
