import axios from "axios";

// const BASE_URL = "http://192.168.8.117:8080/";
const BASE_URL = "http://192.168.43.75:8080/";

export const axiosKaspiInstance = axios.create({
  baseURL: BASE_URL,
});

export const getProcessId = async (amount: number) => {
  const res = await axiosKaspiInstance.get(`payment?amount=${amount}`);
  return res.data;
};

export const getProcessStatus = async (processId: string) => {
  const res = await axiosKaspiInstance.get(`status?processId=${processId}`);
  return res.data;
};
