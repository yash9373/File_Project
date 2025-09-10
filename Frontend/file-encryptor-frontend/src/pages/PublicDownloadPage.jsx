import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { downloadSharedFile } from "../redux/slice/linkslice"; // 👈 add this

export default function PublicDownloadPage() {
  const { token } = useParams();
  const dispatch = useDispatch();
  const [password, setPassword] = useState("");
  const { loading, error } = useSelector((state) => state.link); // 👈 use "link" slice

  const handleDownload = async (e) => {
    e.preventDefault();
    const result = await dispatch(downloadSharedFile({ token, password }));

    if (downloadSharedFile.fulfilled.match(result)) {
      const { blob, filename } = result.payload;

      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    }
  };

  return (
    <div>
      {error && <div>{error}</div>}
      <form onSubmit={handleDownload}>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter password"
        />
        <button type="submit" disabled={loading}>
          {loading ? "Downloading..." : "Download"}
        </button>
      </form>
    </div>
  );
}
