import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import api from "../../axios";

// Login thunk
export const login = createAsyncThunk(
  "auth/login",
  async (credentials, thunkAPI) => {
    try {
      const res = await api.post("/auth/login", credentials);
      // Persist token once on successful login
      localStorage.setItem("token", res.data.token);
      return res.data; // { token, user }
    } catch (err) {
      const message = err?.response?.data?.message || err?.message || "Login failed";
      return thunkAPI.rejectWithValue({ message });
    }
  }
);

// Register thunk
export const register = createAsyncThunk(
  "auth/register",
  async (payload, thunkAPI) => {
    try {
      const res = await api.post("/auth/register", payload);
      // Register returns token and user per spec
      localStorage.setItem("token", res.data.token);  
      return res.data; // { token, user }
    } catch (err) {
      const message = err?.response?.data?.message || err?.message || "Register failed";
      return thunkAPI.rejectWithValue({ message });
    }
  }
);

// Get current user
export const getMe = createAsyncThunk("auth/me", async (_, thunkAPI) => {
  try {
    const res = await api.get("/auth/me");
    return res.data; // { user_id, email }
  } catch (err) {
    const message = err?.response?.data?.message || err?.message || "Fetch user failed";
    return thunkAPI.rejectWithValue({ message });
  }
});

const authSlice = createSlice({
  name: "auth",
  initialState: {
    user: null,
    token: localStorage.getItem("token") || null,
    loading: false,
    error: null,
  },
  reducers: {
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.error = null;
      localStorage.removeItem("token");
    },
  },
  extraReducers: (builder) => {
    builder
      // login
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.token = action.payload.token;
        state.user = action.payload.user ?? state.user;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.message || "Login failed";
      })
      // register
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.token = action.payload.token;
        state.user = action.payload.user;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.message || "Register failed";
      })
      // getMe
      .addCase(getMe.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getMe.fulfilled, (state, action) => {
        state.loading = false;
        // backend returns { user_id, email }
        state.user = action.payload;
      })
      .addCase(getMe.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.message || "Fetch user failed";
      });
  },
});

export const { logout } = authSlice.actions;
export default authSlice.reducer;