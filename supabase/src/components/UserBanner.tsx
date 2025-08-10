import { LogOutIcon, User2Icon } from "lucide-react";
import { toast } from "react-hot-toast";
import useCurrentUser from "@/hooks/useCurrentUser";
import { cn } from "@/lib/utils";
import { Routes } from "@/router";
import { supabaseClient } from "@/store";
import { useTranslate } from "@/utils/i18n";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "./ui/dropdown-menu";

const UserBanner = () => {
  const t = useTranslate();
  const currentUser = useCurrentUser();

  const handleSignOut = async () => {
    const { error } = await supabaseClient.auth.signOut();
    if (error != null) {
      toast.error(error.message);
    } else {
      window.location.href = Routes.AUTH;
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild disabled={!currentUser}>
        <div className={cn("ml-auto w-auto flex flex-row justify-start items-center cursor-pointer text-foreground", "px-1")}>
          {currentUser.email ? (
            <User2Icon className="w-6 mx-auto h-auto text-muted-foreground" />
          ) : (
            <User2Icon className="w-6 mx-auto h-auto text-muted-foreground" />
          )}
        </div>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem onClick={handleSignOut}>
          <LogOutIcon className="size-4 text-muted-foreground" />
            {t("common.sign-out")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default UserBanner;
