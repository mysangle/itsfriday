import { Button, Input } from "@usememos/mui";
import { observer } from "mobx-react-lite";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useTranslate } from "@/utils/i18n";
import { authService, type SignUpRequest } from "@/api";
import useNavigateTo from "@/hooks/useNavigateTo";

const Login = observer(() => {
  const t = useTranslate();
  const navigateTo = useNavigateTo();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  // Redirect to root page if already signed in.

  const mutation = useMutation({
    mutationFn: (loginRequest: SignUpRequest) => {
      return authService.login(loginRequest)
    },
    onSuccess: () => {
      navigateTo("/dashboard");
    },
  })

  const handleUsernameInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setUsername(text);
  };

  const handlePasswordInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setPassword(text);
  };

  const handleLoginButtonClick = async () => {
    if (username === "" || password === "") {
      return;
    }

    mutation.mutate({ username: username, password: password })
  }

  return (
    <div>
      <div>
        <span>{t("common.username")}</span>
        <Input
          placeholder={t("common.username")}
          value={username}
          onChange={handleUsernameInputChanged}
          required
        />
      </div>
      <div>
        <span>{t("common.password")}</span>
        <Input
          placeholder={t("common.password")}
          value={password}
          onChange={handlePasswordInputChanged}
          required
        />
      </div>
      <div>
        <Button
          onClick={handleLoginButtonClick}
        >
          {t("common.log-in")}
        </Button>
      </div>
    </div>
  );
});

export default Login;
