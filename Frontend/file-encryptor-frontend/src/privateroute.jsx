import { Navigate, Outlet } from "react-router-dom";

export default function PrivateRoute(){
    // Use correct token key
    const isauth = !!localStorage.getItem("token");

    return isauth ? <Outlet/> : <Navigate to="/login" replace/>;
}