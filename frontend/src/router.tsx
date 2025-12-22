import { createBrowserRouter, useParams } from 'react-router-dom';
import { StartPage } from './components/pages/Start';
import { SettingsPage } from '@/components/pages/Settings';
import Layout from '@/components/Layout';
import ErrorPage from '@/components/ErrorPage';
// import { redirect } from 'react-router-dom';

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
