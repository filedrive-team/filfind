import { lazy } from 'react';
import { RouteModel } from '@/models/RouteModel';

enum RouterPath {
  home = '/',
  detail = '/detail',
  error = '/error',
  auth = '/auth',
  signUp = '/signUp',
  storageProvide = '/storageProvide/:id',
  client = '/client/:id',
  clientList = '/clientList',
  changePassword = '/changePassword',
  resetPassword = '/resetPassword',
}

const Routes: RouteModel[] = [
  {
    name: 'Home',
    path: RouterPath.home,
    auth: false,
    component: lazy(() => import('@/pages/Home')),
  },
  {
    name: 'StorageProvide',
    path: RouterPath.storageProvide,
    auth: false,
    component: lazy(() => import('@/pages/StorageProvide')),
  },
  {
    name: 'Client',
    path: RouterPath.client,
    auth: false,
    component: lazy(() => import('@/pages/Client')),
  },
  {
    name: 'ClientList',
    path: RouterPath.clientList,
    auth: false,
    component: lazy(() => import('@/pages/ClientList')),
  },
  {
    name: 'detail',
    path: RouterPath.detail,
    auth: true,
    component: lazy(() => import('@/pages/Detail')),
  },
  {
    name: 'auth',
    path: RouterPath.auth,
    auth: false,
    component: lazy(() => import('@/pages/Auth')),
  },
  {
    name: 'resetPassword',
    path: RouterPath.resetPassword,
    auth: false,
    component: lazy(() => import('@/pages/ResetPassword')),
  },
  {
    name: 'signUp',
    path: RouterPath.signUp,
    auth: false,
    component: lazy(() => import('@/pages/SignUp')),
  },
  {
    name: 'changePassword',
    path: RouterPath.changePassword,
    auth: false,
    component: lazy(() => import('@/pages/ChangePassword')),
  },
  {
    name: '404',
    path: RouterPath.error,
    auth: false,
    component: lazy(() => import('@/404')),
  },
];

export { RouterPath, Routes };
