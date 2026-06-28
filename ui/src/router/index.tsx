import { Suspense, lazy, useEffect, useState } from 'react';
import { RouteObject } from 'react-router-dom';

import Layout from '@/pages/Layout';
import { mergeRoutePlugins } from '@/utils/pluginKit';

import baseRoutes, { RouteNode } from './routes';
import RouteGuard from './RouteGuard';
import RouteErrorBoundary from './RouteErrorBoundary';

const routeWrapper = (routeNodes: RouteNode[], root: RouteNode[]) => {
  routeNodes.forEach((rn) => {
    if (rn.page === 'pages/Layout') {
      rn.element = rn.guard ? (
        <RouteGuard onEnter={rn.guard} path={rn.path} page={rn.page}>
          <Layout />
        </RouteGuard>
      ) : (
        <Layout />
      );
      rn.errorElement = <RouteErrorBoundary />;
    } else {
      /**
       * cannot use a fully dynamic import statement
       * ref: https://webpack.js.org/api/module-methods/#import-1
       */

      let Ctrl;

      if (typeof rn.page === 'string') {
        const pagePath = rn.page.replace('pages/', '');
        Ctrl = lazy(() => import(`@/pages/${pagePath}`));
        console.log('route page:', rn.page);
        console.log('pagePath:', pagePath);
      } else {
        Ctrl = rn.page;
      }

      rn.element = (
        <Suspense>
          {rn.guard ? (
            <RouteGuard onEnter={rn.guard} path={rn.path} page={rn.page}>
              <Ctrl />
            </RouteGuard>
          ) : (
            <Ctrl />
          )}
        </Suspense>
      );
      rn.errorElement = <RouteErrorBoundary />;
    }
    root.push(rn);
    const children = Array.isArray(rn.children) ? rn.children : null;
    if (children) {
      rn.children = [];
      routeWrapper(children, rn.children);
    }
  });
};

function useMergeRoutes() {
  const [routesState, setRoutes] = useState<RouteObject[]>([]);

  const init = async () => {
    const routes = [];
    let mergedRoutes = baseRoutes;

    try {
      mergedRoutes = await mergeRoutePlugins(baseRoutes);
    } catch (err) {
      console.error('mergeRoutePlugins error:', err);
    }
    console.log('baseRoutes:', baseRoutes);
    console.log('mergedRoutes:', mergedRoutes);
    routeWrapper(mergedRoutes, routes);
    setRoutes(routes);
  };

  useEffect(() => {
    init();
  }, []);

  return routesState;
}

export { useMergeRoutes };
