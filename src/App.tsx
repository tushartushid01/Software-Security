import { Outlet, RouteObject, createBrowserRouter } from "react-router-dom";
import { Login } from "./Pages/login";
import { ForgotPassword } from "./Pages/forgetPassword";
import { Dashboard } from "./Pages/dashboard";

const GlobalProvider = () => {
  return <Outlet />;
};

const AuthenticatedRoutes: RouteObject = {
  path: "/",
  children: [
    { path: "/", Component: Login },
    {
      path: "/forget-password",
      Component: ForgotPassword,
    },
    {
      path: "/dashboard",
      Component: Dashboard,
    },
  ],
};

export const App = createBrowserRouter([
  {
    path: "/",
    Component: GlobalProvider,
    children: [AuthenticatedRoutes],
  },
]);
