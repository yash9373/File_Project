import { configureStore } from '@reduxjs/toolkit'
import authReducer from './slice/authslice.jsx'
import filereducer from './slice/fileslice.jsx'
import linkreducer from './slice/linkslice.jsx'
export const store = configureStore({
  reducer: {
    auth: authReducer,
    file:filereducer,
    link:linkreducer
  },
})