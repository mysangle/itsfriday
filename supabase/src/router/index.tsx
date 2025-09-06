import { Suspense, lazy } from "react";
import { createBrowserRouter } from "react-router-dom";
import App from "@/App";
import RootLayout from "@/layouts/RootLayout";
import Home from "@/pages/Home";
import Loading from "@/pages/Loading";

const Libro = lazy(() => import("@/pages/Libro"));
const Monero = lazy(() => import("@/pages/Monero"));
const NotFound = lazy(() => import("@/pages/NotFound"));
const SignIn = lazy(() => import("@/pages/SignIn"));
const SignUp = lazy(() => import("@/pages/SignUp"));

export enum Routes {
  ROOT = "/",
  AUTH = "/auth",
  LIBRO = "/libro",
  MONERO = "/monero",
}

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: Routes.AUTH,
        children: [
          {
            path: "",
            element: (
              <Suspense fallback={<Loading />}>
                <SignIn />
              </Suspense>
            ),
          },
          {
            path: "signup",
            element: (
              <Suspense fallback={<Loading />}>
                <SignUp />
              </Suspense>
            ),
          },
        ],
      },
      {
        path: Routes.ROOT,
        element: <RootLayout />,
        children: [
          {
            path: "",
            element: <Home />,
          },
          {
            path: Routes.LIBRO,
            element: (
              <Suspense fallback={<Loading />}>
                <Libro />
              </Suspense>
            ),
          },
          {
            path: Routes.MONERO,
            element: (
              <Suspense fallback={<Loading />}>
                <Monero />
              </Suspense>
            ),
          },
        ],
      },
      {
        path: "404",
        element: (
          <Suspense fallback={<Loading />}>
            <NotFound />
          </Suspense>
        ),
      },
    ],
  },
]);

export default router;

