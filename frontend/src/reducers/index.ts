import { combineReducers } from '@reduxjs/toolkit';
import idempotencySlice from './idempotencySlice';

const rootReducer = combineReducers({
    idempotency: idempotencySlice,
});

export default rootReducer;