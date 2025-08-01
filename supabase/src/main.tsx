import { observer } from "mobx-react-lite";
import { createRoot } from 'react-dom/client'
import { Toaster } from "react-hot-toast";
import { RouterProvider } from "react-router-dom";
import './index.css'
import router from "./router";

const Main = observer(() => (
  <>
    <RouterProvider router={router} />
    <Toaster position="top-right" />
  </>
));

(async () => {
  const container = document.getElementById("root");
  const root = createRoot(container as HTMLElement);
  root.render(<Main />);
})();
