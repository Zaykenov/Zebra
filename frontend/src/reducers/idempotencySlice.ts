import { createSlice } from "@reduxjs/toolkit";
import uuid from "react-uuid";

interface IdempotencyState {
  idempotencyKey: string;
}

const generateInitialIdempotencyKey = () => `${uuid()}_1`;

const initialState: IdempotencyState = {
  idempotencyKey: generateInitialIdempotencyKey()
};

const idempotencySlice = createSlice({
  name: "idempotency",
  initialState,
  reducers: {
    regenerateIdempotency: (state) => {
      state.idempotencyKey = generateInitialIdempotencyKey();
    },
    updateIdempotencyVersion: (state) => {
      const [orderUUID, versionStr] = state.idempotencyKey.split('_');
      const newVersion = parseInt(versionStr) + 1;
      state.idempotencyKey = `${orderUUID}_${newVersion}`;
    }
  },
});

export const { regenerateIdempotency, updateIdempotencyVersion } = idempotencySlice.actions;

export default idempotencySlice.reducer;