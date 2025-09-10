import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import api from "../../axios"; // Using the centralized axios instance

/**
 * Creates a share link for a given file.
 * @param {object} shareData - The data for creating the link.
 * @param {string} shareData.file_id - The ID of the file to share.
 * @param {number|null} [shareData.expires_in_minutes] - Optional expiry time.
 * @param {number|null} [shareData.max_downloads] - Optional download limit.
 */
export const createShareLink = createAsyncThunk(
  "share/createLink",
  async (shareData, { rejectWithValue }) => {
    try {
      // The auth token is now handled by the axios interceptor
      const response = await api.post("/share", shareData);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);
import axios from "axios";

export const downloadSharedFile = createAsyncThunk(
  "share/downloadFile",
  async ({ token, password }, { rejectWithValue }) => {
    try {
      const response = await axios.post(
        `http://localhost:8080/share/${token}/download`,
        { password },
        { responseType: "blob" } // blob = file data
      );

      return {
        blob: response.data,
        filename:
          response.headers["content-disposition"]?.split("filename=")[1] ||
          "file",
      };
    } catch (error) {
      return rejectWithValue(
        error.response?.data || { error: "Failed to download file." }
      );
    }
  }
);

/**
 * Deletes a share link.
 * @param {string} token - The token of the share link to delete.
 */
export const deleteShareLink = createAsyncThunk(
  "share/deleteLink",
  async (token, { rejectWithValue }) => {
    try {
      // The auth token is now handled by the axios interceptor
      await api.delete(`/share/${token}`);
      return token; // Return the token for potential state updates
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);


const shareSlice = createSlice({
  name: "share",
  initialState: {
    currentLink: null, // To store the details of a newly created link
    loading: false,
    error: null,
  },
  reducers: {
    clearCurrentLink: (state) => {
      state.currentLink = null;
    },
    clearError: (state) => {
        state.error = null;
    }
  },
  extraReducers: (builder) => {
    builder
      // Create Share Link Reducers
      .addCase(createShareLink.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.currentLink = null;
      })
      .addCase(createShareLink.fulfilled, (state, action) => {
        state.loading = false;
        state.currentLink = action.payload;
      })
      .addCase(createShareLink.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "Failed to create share link.";
      })
      // Delete Share Link Reducers
      .addCase(deleteShareLink.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(deleteShareLink.fulfilled, (state, action) => {
        state.loading = false;
        // If the deleted link is the one we are currently showing, clear it
        if (state.currentLink?.token === action.payload) {
            state.currentLink = null;
        }
      })
      .addCase(deleteShareLink.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "Failed to delete share link.";
      });
  },
});

export const { clearCurrentLink, clearError } = shareSlice.actions;

export default shareSlice.reducer;

