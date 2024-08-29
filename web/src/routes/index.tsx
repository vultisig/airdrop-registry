import {
  RouteObject,
  RouterProvider,
  Navigate,
  createBrowserRouter,
} from "react-router-dom";
import React from "react";

import { useVaultContext } from "context";
import constantPaths from "routes/constant-paths";

import Layout from "layout";

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
  const { vaults } = useVaultContext();

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
      path: constantPaths.root,
      redirect: vaults.length ? constantPaths.balance : constantPaths.landing,
    },
    {
      path: constantPaths.landing,
      element: <LandingPage />,
    },
    ...(vaults.length
      ? [
          {
            path: constantPaths.root,
            element: <Layout />,
            children: [
              {
                path: constantPaths.root,
                redirect: constantPaths.balance,
              },
              {
                path: constantPaths.balance,
                element: <BalancePage />,
              },
              {
                path: constantPaths.asset,
                element: <AssetPage />,
              },
              {
                path: "*",
                redirect: constantPaths.root,
              },
            ],
          },
        ]
      : []),

    {
      path: "*",
      redirect: constantPaths.root,
    },
  ];

  const router = createBrowserRouter(processRoutes(routes), {
    basename: constantPaths.basePath,
  });

  return <RouterProvider router={router} />;
};

export type { RouteConfig };

export default Component;
