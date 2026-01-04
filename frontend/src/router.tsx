import { createBrowserRouter } from 'react-router-dom';
import { StartPage } from './components/pages/Start';
import { SettingsPage } from '@/components/pages/Settings';
import Layout from '@/components/Layout';
import { RepoViewPage } from '@/components/views/RepoViewPage';
import ErrorPage from '@/components/ErrorPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        path: '/',
        element: <StartPage />,
      },
    ],
  },
  {
    path: '/resource',
    element: <Layout />,
    children: [
      {
        path: '/resource/:user/:repoPath',
        element: <RepoViewPage />,
      },
    ],
    errorElement: <ErrorPage />,
  },
  {
    path: '/settings',
    element: <Layout />,
    children: [
      {
        path: '/settings',
        element: <SettingsPage />,
      },
    ],
  },
]);
