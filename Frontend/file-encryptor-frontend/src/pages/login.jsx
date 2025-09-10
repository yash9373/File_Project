import { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { login } from "../redux/slice/authslice";

export default function LoginPage() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { loading, error } = useSelector((state) => state.auth);
  const [formData, setFormData] = useState({ email: "", password: "" });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!formData.email || !formData.password) {
      alert("Email and password are required!");
      return;
    }
    dispatch(login(formData))
      .unwrap()
      .then(() => {
        navigate("/");
      })
      .catch((err) => alert(err?.message || "Login failed"));
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <form onSubmit={handleSubmit} className="w-full max-w-md bg-white p-8 rounded-2xl shadow-lg border border-gray-200">
        <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">Login</h2>

        <div className="space-y-4">
          <div>
            <label htmlFor="email" className="block text-sm font-semibold text-gray-600 mb-1">Email</label>
            <input
              id="email"
              type="email"
              name="email"
              placeholder="Enter your email"
              value={formData.email}
              onChange={handleChange}
              className="w-full px-4 py-3 bg-slate-100 border-transparent rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-semibold text-gray-600 mb-1">Password</label>
            <input
              id="password"
              type="password"
              name="password"
              placeholder="Enter your password"
              value={formData.password}
              onChange={handleChange}
              className="w-full px-4 py-3 bg-slate-100 border-transparent rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
            />
          </div>
        </div>

        <button type="submit" className="w-full mt-8 bg-blue-600 text-white font-bold py-3 rounded-lg hover:bg-blue-700 transition-colors shadow">
          {loading ? "Logging in..." : "Login"}
        </button>
        {error && <p className="text-red-600 mt-2">{error}</p>}

        <div className="text-black mt-4">
          <p>Don’t have an account?</p>
          <p onClick={() => navigate("/register")} className="text-blue-600 cursor-pointer hover:underline">Register</p>
        </div>
      </form>
    </div>
  );
}
