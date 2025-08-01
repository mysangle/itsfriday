import { observer } from "mobx-react-lite";
import { Outlet } from "react-router-dom";
import './App.css'

const App = observer(() => {
  return (
    <div>
      <h1>It's Friday</h1>
    </div>
  );
});

export default App
