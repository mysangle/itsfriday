import { observer } from "mobx-react-lite";
import { Outlet } from "react-router-dom";
import './App.css'

const App = observer(() => {
  return <Outlet />;
});

export default App
