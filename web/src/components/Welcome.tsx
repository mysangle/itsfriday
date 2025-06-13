import { observer } from "mobx-react-lite";

const Welcome = observer(() => {
  return (
    <div>
        <h1>Welcome to It's Friday</h1>
         <a href="./auth/login">Log in</a> 
    </div>
  );
});

export default Welcome;
