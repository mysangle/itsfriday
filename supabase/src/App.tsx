import { observer } from "mobx-react-lite";
import './App.css'

const App = observer(() => {
  return (
    <div>
      <h1 className="text-3xl font-bold underline">It's Friday</h1>
    </div>
  );
});

export default App
