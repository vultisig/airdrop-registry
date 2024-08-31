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
import ChainsPage from "pages/chains";
import ImportPage from "pages/import";

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
      redirect: vaults.length ? constantPaths.chains : constantPaths.import,
    },
    {
      path: constantPaths.import,
      element: <ImportPage />,
    },
    ...(vaults.length
      ? [
          {
            path: constantPaths.root,
            element: <Layout />,
            children: [
              {
                path: constantPaths.root,
                redirect: constantPaths.chains,
              },
              {
                path: constantPaths.chains,
                element: <ChainsPage />,
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
