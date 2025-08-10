import { observer } from "mobx-react-lite";
import { createRoot } from 'react-dom/client'
import { Toaster } from "react-hot-toast";
import { RouterProvider } from "react-router-dom";
import { initialUserStore } from "@/store/user";
import './index.css'
import router from "./router";

const Main = observer(() => (
  <>
    <RouterProvider router={router} />
    <Toaster position="top-right" />
  </>
));

(async () => {
  await initialUserStore();
  
  const container = document.getElementById("root");
  const root = createRoot(container as HTMLElement);
  root.render(<Main />);
})();
