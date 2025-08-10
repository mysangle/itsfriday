import { observer } from "mobx-react-lite";
import useCurrentUser from "@/hooks/useCurrentUser";
import { cn } from "@/lib/utils";
import UserBanner from "./UserBanner";

interface Props {
  className?: string;
}

const Header = observer((props: Props) => {
  const { className } = props;
  const currentUser = useCurrentUser();

  return (
    <header className={cn("w-full text-center text-xs font-bold underline", className)}>
      {currentUser && (
        <div className={cn("w-full flex flex-col justify-end", "items-center")}>
          <UserBanner />
        </div>
      )}
    </header>
  );
});

export default Header;
