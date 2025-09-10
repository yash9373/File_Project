import React, { useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios"; // Import the original axios library

const PublicDownloadPage = () => {
  const { token } = useParams();
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleDownload = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      // **IMPORTANT FIX**: Manually construct the full URL without the `/api` prefix
      const downloadUrl = `http://localhost:8080/share/${token}/download?password=${password}`;

      const response = await axios.get(downloadUrl, {
        responseType: "blob", // Important for file downloads
      });

      // Get filename from content-disposition header
      const contentDisposition = response.headers["content-disposition"];
      let filename = "downloaded-file"; // fallback filename
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="(.+)"/);
        if (filenameMatch.length > 1) {
          filename = filenameMatch[1];
        }
      }

      // Create a link and trigger the download
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", filename);
      document.body.appendChild(link);
      link.click();
      link.remove();

    } catch (err) {
      setError(err.response?.data?.error || "Download failed. Check the password or link.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto mt-10 max-w-md text-center">
      <h1 className="text-2xl font-bold mb-4">Download Shared File</h1>
      <p className="mb-4">
        This file is password protected. Please enter the password to download.
      </p>
      <form onSubmit={handleDownload}>
        <div className="mb-4">
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Enter file password"
            className="w-full px-3 py-2 border rounded-md"
            required
          />
        </div>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 disabled:bg-gray-400"
        >
          {loading ? "Downloading..." : "Download"}
        </button>
      </form>
    </div>
  );
};

export default PublicDownloadPage;