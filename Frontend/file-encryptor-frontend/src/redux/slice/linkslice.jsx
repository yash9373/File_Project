import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import api from "../../axios";

export const createShareLink = createAsyncThunk(
  "share/createLink",
  // This accepts a payload object (e.g., { file_id: "..." })
  async (payload, { rejectWithValue }) => {
    try {
      // It sends that payload directly to the backend
      const response = await api.post("/share", payload);

      const token = response.data.token;
      const shareableLink = `${window.location.origin}/share/${token}`;
      
      // It returns the final, complete URL
      return { link: shareableLink };
    } catch (error) {
      console.error("Error creating share link:", error.response.data);
      return rejectWithValue(error.response.data);
    }
  }
);

const shareSlice = createSlice({
  name: "share",
  initialState: {
    shareableLink: null,
    loading: false,
    error: null,
  },
  reducers: {
    resetShareLink: (state) => {
      state.shareableLink = null;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(createShareLink.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.shareableLink = null;
      })
      .addCase(createShareLink.fulfilled, (state, action) => {
        state.loading = false;
        state.shareableLink = action.payload.link;
      })
      .addCase(createShareLink.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "Failed to create share link.";
      });
  },
});

export const { resetShareLink } = shareSlice.actions;
export default shareSlice.reducer;