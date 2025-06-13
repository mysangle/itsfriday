import { Suspense, lazy } from 'react';
import { createBrowserRouter } from 'react-router-dom';
import { Layout } from '@/components/Layout';
import Loading from "@/components/Loading";

const Main = lazy(() => import("@/components/Main"));
const SignUp = lazy(() => import("@/components/SignUp"));
const Login = lazy(() => import("@/components/Login"));
const Welcome = lazy(() => import("@/components/Welcome"));

export enum Routes {
  ROOT = "dashboard",
  AUTH = "auth",
  WELCOME = "welcome"
}

const router = createBrowserRouter([
  {
    path: Routes.ROOT,
    element: <Layout />,
    children: [
      {
        path: "",
        element: <Main />
      },
      {
        path: Routes.WELCOME,
        element: <Welcome />
      },
      {
        path: Routes.AUTH,
        children: [
          {
            path: "",
            element: <Layout />,
          },
          {
            path: "signup",
            element: (
              <Suspense fallback={<Loading />}>
                <SignUp />
              </Suspense>
            ),
          },
          {
            path: "login",
            element: (
              <Suspense fallback={<Loading />}>
                <Login />
              </Suspense>
            ),
          },
        ],
      },
    ]
  }
]);

export default router;
