import { observer } from "mobx-react-lite";
import { Suspense, useEffect } from "react";
import { Outlet } from "react-router-dom";
import Header from "@/components/Header";
import Navigation from "@/components/Navigation";
import useCurrentUser from "@/hooks/useCurrentUser";
import { cn } from "@/lib/utils";
import Loading from "@/pages/Loading";
import { Routes } from "@/router";

const RootLayout = observer(() => {
  const currentUser = useCurrentUser();

  useEffect(() => {
    if (!currentUser) {
      window.location.href = Routes.AUTH;
    }
  }, []);

  return (
    <div className="w-full min-h-full flex flex-row justify-center items-start pl-16">
      <div
        className={cn(
            "group flex flex-row justify-between items-center fixed top-0 left-0 right-0 select-none h-16 bg-sidebar",
            "h-8 px-2",
            "border-b border-border",
          )}
      >
        <Header className="py-4 md:pl-6" />
      </div>
      <div
        className={cn(
            "group flex flex-col justify-start items-start fixed top-8 left-0 select-none h-full bg-sidebar",
            "w-16 px-2",
            "border-r border-border",
          )}
      >
        <Navigation className="py-4 pt-6" />
      </div>
      <main className="w-full h-auto grow shrink flex flex-col justify-start items-center">
        <Suspense fallback={<Loading />}>
          <Outlet />
        </Suspense>
      </main>
    </div>
  );
});

export default RootLayout;
