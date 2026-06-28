/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
  {
    path: '/',
    page: 'pages/Layout',
    loader: async () => {
      await guard.setupApp();
      return null;
    },
    guard: () => {
      // const gr = guard.shouldLoginRequired();
      // if (!gr.ok) {
      //   return gr;
      // }
      return {
        ok: true,
      };
    },
    children: [
      // {
      //   path: 'auth/login',
      //   page: 'pages/Auth/Login',
      //   guard: () => {
      //     const notLogged = guard.notLogged();
      //     if (notLogged.ok) {
      //       return notLogged;
      //     }

      //     return guard.notActivated();
      //   },
      // },
      {
        path: '*',
        page: 'pages/404',
      },
      {
        path: '50x',
        page: 'pages/50X',
      },
    ],
  },
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
          const notLogged = guard.notLogged();
          if (notLogged.ok) {
            return notLogged;
          }

          return guard.notActivated();
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
