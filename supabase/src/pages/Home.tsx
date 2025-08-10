import { observer } from "mobx-react-lite";

const Home = observer(() => {
  return (
    <div className="w-full min-h-full bg-background text-foreground">
        <div className="text-3xl font-bold underline">It's Friday</div>
    </div>
  );
});

export default Home;
