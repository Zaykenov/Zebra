import { backendApiInstance } from ".";

export const getUserByCode = async (id: string) => {
    const res = await backendApiInstance.get(`user/getUserByCode/${id}`);
    return res.data;
};