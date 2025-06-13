import { Button } from "@usememos/mui";
import { observer } from "mobx-react-lite";
import { useTranslate } from "@/utils/i18n";
import { useMutation } from "@tanstack/react-query";
import { authService } from "@/api";
import useNavigateTo from "@/hooks/useNavigateTo";

const Main = observer(() => {
  const t = useTranslate();
  const navigateTo = useNavigateTo();

  const mutation = useMutation({
      mutationFn: () => {
        return authService.logout()
      },
      onSuccess: () => {
        navigateTo("/dashboard/auth/login");
      },
    })

  const handleLogoutButtonClick = async () => {
    mutation.mutate()
  }

  return (
    <div>
        <h1>Dashboard</h1>
         <div>
        <Button
          onClick={handleLogoutButtonClick}
        >
          {t("common.log-out")}
        </Button>
      </div>
    </div>
  );
});

export default Main;
