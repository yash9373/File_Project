import { createBrowserRouter } from "react-router-dom";

import RegisterPage from "./pages/register";
import LoginPage from "./pages/login";
import PrivateRoute from "./privateroute";
import Layout from "./layout";
import FileUploadForm from "./pages/file_add_form";
import PublicDownloadPage from "./pages/PublicDownloadPage";
// 1. Import your new component

export const router = createBrowserRouter([
  { path: "/register", Component: RegisterPage },
  { path: "/login", Component: LoginPage },
  

  { path: "/share/:token", Component: PublicDownloadPage },

  {
    Component: PrivateRoute, // protect all below
    children: [
      { path: "/", Component: Layout },
      { path: "/upload", Component: FileUploadForm },
    ],
  },
]);