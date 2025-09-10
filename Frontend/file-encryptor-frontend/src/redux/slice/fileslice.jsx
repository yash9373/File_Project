import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import api from "../../axios"; // Using the centralized axios instance

// ## Async Thunks for API Calls ##

// 1. Fetch all files for the logged-in user
export const fetchFiles = createAsyncThunk(
  "file/fetchFiles",
  async (_, { rejectWithValue }) => {
    try {
      // The auth token is now handled by the axios interceptor
      const response = await api.get("/files");
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

// 2. Upload a new file
export const uploadFile = createAsyncThunk(
  "file/uploadFile",
  async ({ file, password }, { rejectWithValue }) => {
    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("password", password);

      // The auth token is handled by the interceptor, but we still need the Content-Type header
      const response = await api.post("/files/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

// 3. Delete a file
export const deleteFile = createAsyncThunk(
  "file/deleteFile",
  async (fileId, { rejectWithValue }) => {
    try {
      // The auth token is now handled by the axios interceptor
      await api.delete(`/files/${fileId}`);
      return fileId; // Return the ID for removal from state
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

// Note: This is a helper, not a thunk, as it doesn't dispatch actions.
export const downloadFile = async (fileId, password) => {
  try {
    // The auth token is now handled by the axios interceptor
    const response = await api.get(`/files/${fileId}/download`, {
        params: { password },
        responseType: 'blob',
    });
    return response.data;
  } catch (error) {
    throw error.response.data;
  }
};


const fileSlice = createSlice({
  name: "file",
  initialState: {
    files: [],
    loading: false,
    error: null,
  },
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    // This action can be dispatched on logout to clear the user's file data
    clearFiles: (state) => {
        state.files = [];
        state.error = null;
        state.loading = false;
    }
  },
  extraReducers: (builder) => {
    builder
      // Fetch Files Reducers
      .addCase(fetchFiles.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchFiles.fulfilled, (state, action) => {
        state.loading = false;
        state.files = action.payload;
      })
      .addCase(fetchFiles.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "Failed to fetch files.";
      })
      // Upload File Reducers
      .addCase(uploadFile.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(uploadFile.fulfilled, (state, action) => {
        state.loading = false;
        state.files.push(action.payload);
      })
      .addCase(uploadFile.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "File upload failed.";
      })
      // Delete File Reducers
      .addCase(deleteFile.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(deleteFile.fulfilled, (state, action) => {
        state.loading = false;
        state.files = state.files.filter((file) => file.id !== action.payload);
      })
      .addCase(deleteFile.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload.error || "Failed to delete file.";
      });
  },
});

export const { clearError, clearFiles } = fileSlice.actions;

export default fileSlice.reducer;

