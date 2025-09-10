import React, { useState } from "react";
// Removed Redux and React Router hooks to make the component self-contained.
// In a real application, you would use these hooks as you had them.

export default function FileUploadForm() {
  const [file, setFile] = useState(null);
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Helper function to simulate navigation
  const navigate = (path) => {
    if (path === -1) {
      window.history.back();
    } else {
       // In a real app, react-router-dom's navigate would handle this.
       alert(`Navigating to ${path}`);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!file) {
      alert("Please upload a file!");
      return;
    }
    if (password.length < 6) {
      alert("Password must be at least 6 characters long!");
      return;
    }

    setLoading(true);
    setError(null);

    // This token would typically come from your Redux store or auth context
    const token = localStorage.getItem("token"); 
    if (!token) {
        alert("Authentication error: No token found. Please log in again.");
        setLoading(false);
        return;
    }

    const formData = new FormData();
    formData.append("file", file);
    formData.append("password", password);

    try {
      const response = await fetch("http://localhost:8080/api/files/upload", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Something went wrong");
      }

      alert("File uploaded successfully!");
      navigate("/"); // Navigate back to the dashboard on success
    } catch (err) {
      setError(err.message);
      console.error("Upload failed:", err);
      alert(`Upload failed: ${err.message || "Please try again."}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <form
        onSubmit={handleSubmit}
        className="bg-white p-6 rounded-2xl shadow-md space-y-4 w-80"
      >
        <h2 className="text-xl font-semibold text-center">Upload a Secure File</h2>

        {/* File Upload */}
        <input
          type="file"
          onChange={(e) => setFile(e.target.files[0])}
          className="block w-full text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
        />

        {/* Password */}
        <input
          type="password"
          placeholder="Enter encryption password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full px-3 py-2 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 text-black"
        />
        
        {error && <p className="text-sm text-center text-red-500">{error}</p>}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 rounded-xl hover:bg-blue-700 transition disabled:bg-blue-300"
        >
          {loading ? "Uploading..." : "Upload and Encrypt"}
        </button>

        {/* Back Button */}
        <button
          type="button"
          onClick={() => navigate(-1)} // Navigates to the previous page
          className="w-full bg-gray-200 text-gray-700 py-2 rounded-xl hover:bg-gray-300 transition"
        >
          Back
        </button>
      </form>
    </div>
  );
}

