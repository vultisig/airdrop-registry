import {
  RouteObject,
  RouterProvider,
  Navigate,
  createBrowserRouter,
} from "react-router-dom";
import React from "react";
import paths from "routes/constant-paths";

import DefaultLayout from "layouts/default";

import AssetPage from "pages/asset";
import BalancePage from "pages/balance";
import LandingPage from "pages/landing";

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
      //should remove when redux added to project
    },
    {
      path: paths.landing,
      element: <LandingPage />,
    },
    {
      path: paths.root,
      element: <DefaultLayout />,
      children: [
        {
          path: paths.root,
          redirect: paths.balance,
        },
        {
          path: paths.balance,
          element: <BalancePage />,
        },
        {
          path: paths.asset,
          element: <AssetPage />,
        },
        {
          path: "*",
          redirect: paths.root,
        },
      ],
    },
    {
      path: "*",
      redirect: paths.root,
    },
  ];

  const router = createBrowserRouter(processRoutes(routes));

  return <RouterProvider router={router} />;
};

export type { RouteConfig };

export default Component;
