import type { IndexRouteObject, NonIndexRouteObject } from 'react-router-dom';

import { guard } from '@/utils';
import type { TGuardFunc } from '@/utils/guard';

type IndexRouteNode = Omit<IndexRouteObject, 'children'>;
type NonIndexRouteNode = Omit<NonIndexRouteObject, 'children'>;
type UnionRouteNode = IndexRouteNode | NonIndexRouteNode;

export type RouteNode = UnionRouteNode & {
  page: string;
  children?: RouteNode[];
  /**
   * a method to auto guard route before route enter
   * if the `ok` field in guard returned `TGuardResult` is true,
   * it means the guard passed then enter the route.
   * if guard returned the `TGuardResult` has `redirect` field,
   * then auto redirect route to the `redirect` target.
   */
  guard?: TGuardFunc;
};

const routes: RouteNode[] = [
  // {
  //   path: '/',
  //   page: 'pages/Layout',
  //   loader: async () => {
  //     await guard.setupApp();
  //     return null;
  //   },
  //   guard: () => {
  //     // const gr = guard.shouldLoginRequired();
  //     // if (!gr.ok) {
  //     //   return gr;
  //     // }
  //     return {
  //       ok: true,
  //     };
  //   },
  //   children: [
  //     // {
  //     //   path: 'auth/login',
  //     //   page: 'pages/Auth/Login',
  //     //   guard: () => {
  //     //     const notLogged = guard.notLogged();
  //     //     if (notLogged.ok) {
  //     //       return notLogged;
  //     //     }

  //     //     return guard.notActivated();
  //     //   },
  //     // },
  //     {
  //       path: '*',
  //       page: 'pages/404',
  //     },
  //     {
  //       path: '50x',
  //       page: 'pages/50X',
  //     },
  //   ],
  // },
  {
    path: '/',
    page: 'pages/Layout/AuthLayout',
    loader: async () => {
      await guard.setupApp();
      return null;
    },
    guard: () => {
      const gr = guard.shouldLoginRequired();
      if (!gr.ok) {
        return gr;
      }

      return {
        ok: true,
      };
    },
    children: [
      {
        path: 'auth/login',
        page: 'pages/Auth/Login',
        guard: () => {
          return guard.notLogged();
          // if (notLogged.ok) {
          //   return notLogged;
          // }

          // return guard.notActivated();
        },
      },
      {
        path: 'auth/account-recovery',
        page: 'pages/Auth/RecoveryAccount',
        guard: () => {
          return guard.notLogged();
        },
      },
      {
        path: 'auth/register',
        page: 'pages/Auth/Register',
        guard: () => {
          return guard.notLogged();
        },
      },
    ],
  },
  {
    path: '/',
    page: 'pages/Layout',
    loader: async () => {
      await guard.setupApp();
      return null;
    },
    children: [
      {
        path: '403',
        page: 'pages/403',
      },
    ],
  },
  {
    path: '/install',
    page: 'pages/Install',
  },
  {
    path: '/maintenance',
    page: 'pages/Maintenance',
  },
];
export default routes;
