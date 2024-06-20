import { backendApiInstance } from ".";

export interface IdempotencyKey {
  idempotency_key: string;
  time: string;
}

export const sendIdempotencyKeys = async (data: IdempotencyKey[]) => {
  try {
    const res = await backendApiInstance.post("check/idempotencyCheck", {keys: data});
    return res.data;
  } catch (e) {
    throw e;
  }
}