import { CreditCardIcon, HouseIcon, LibraryIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { NavLink } from "react-router-dom";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { Routes } from "@/router";
import { useTranslate } from "@/utils/i18n";

interface NavLinkItem {
  id: string;
  path: string;
  title: string;
  icon: React.ReactNode;
}

interface Props {
  className?: string;
}

const Navigation = observer((props: Props) => {
  const { className } = props;
  const t = useTranslate();

  const homeNavLink: NavLinkItem = {
    id: "header-home",
    path: Routes.ROOT,
    title: t("common.home"),
    icon: <HouseIcon className="w-6 h-auto shrink-0" />,
  };
  const libroNavLink: NavLinkItem = {
    id: "header-libro",
    path: Routes.LIBRO,
    title: t("common.libro"),
    icon: <LibraryIcon className="w-6 h-auto shrink-0" />,
  };
  const dineroNavLink: NavLinkItem = {
    id: "header-dinero",
    path: Routes.DINERO,
    title: t("common.dinero"),
    icon: <CreditCardIcon className="w-6 h-auto shrink-0" />,
  };

  const navLinks: NavLinkItem[] = [homeNavLink, libroNavLink, dineroNavLink];

  return (
    <div className={cn("w-full px-1 py-4 flex flex-col justify-start items-start space-y-2 overflow-auto overflow-x-hidden hide-scrollbar shrink", className)}>
      {navLinks.map((navLink) => (
        <NavLink
          className={({ isActive }) =>
            cn(
              "px-2 py-2 rounded-2xl border flex flex-row items-center text-lg text-sidebar-foreground transition-colors",
              "",
              isActive
                ? "bg-sidebar-accent text-sidebar-accent-foreground border-sidebar-accent-border drop-shadow"
                : "border-transparent hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:border-sidebar-accent-border opacity-80",
            )
          }
          key={navLink.id}
          to={navLink.path}
          id={navLink.id}
          viewTransition
        >
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <div>{navLink.icon}</div>
              </TooltipTrigger>
              <TooltipContent side="right">
                <p>{navLink.title}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </NavLink>
      ))}
    </div>
  );
});

export default Navigation;
