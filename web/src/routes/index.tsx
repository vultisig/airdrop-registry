import {
  RouteObject,
  RouterProvider,
  Navigate,
  createBrowserRouter,
} from "react-router-dom";
import React from "react";
import paths from "./constant-paths";
import LandingPage from "../pages/landing";

interface RouteConfig {
  path: string;
  element?: React.ReactNode;
  children?: RouteConfig[];
  redirect?: string;
}

const Component = () => {
  const processRoutes = (routes: RouteConfig[]): RouteObject[] => {
    return routes.reduce<RouteObject[]>(
      (acc, { children, element, path, redirect }) => {
        if (redirect) {
          const processedRoute: RouteObject = {
            path,
            element: <Navigate to={redirect} replace />,
          };
          acc.push(processedRoute);
        } else if (element) {
          const processedRoute: RouteObject = {
            path,
            element,
            children: children ? processRoutes(children) : undefined,
          };
          acc.push(processedRoute);
        }

        return acc;
      },
      []
    );
  };

  const routes: RouteConfig[] = [
    {
      path: paths.root,
      redirect: paths.landing,
    },
    {
      path: paths.landing,
      element: <LandingPage />,
    },
  ];

  const router = createBrowserRouter(processRoutes(routes));

  return <RouterProvider router={router} />;
};

export type { RouteConfig };

export default Component;
