import { StrictMode } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { Layout } from './layout';

const Dashboard = () => <h1>Dashboard</h1>;

const router = createBrowserRouter([
  {
    path: 'dashboard',
    element: <Layout />,
    children: [
      { path: '', element: <Dashboard /> },
    ]
  }
]);

export const App = () => {
  return (
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>
  );
};





