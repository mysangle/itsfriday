import { Button, Input } from "@usememos/mui";
import { observer } from "mobx-react-lite";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useTranslate } from "@/utils/i18n";
import { authService, type SignUpRequest } from "@/api";
import useNavigateTo from "@/hooks/useNavigateTo";

const SignUp = observer(() => {
  const t = useTranslate();
  const navigateTo = useNavigateTo();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const mutation = useMutation({
    mutationFn: (signUpRequest: SignUpRequest) => {
      return authService.signUp(signUpRequest)
    },
    onSuccess: () => {
      navigateTo("/dashboard/welcome");
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

  const handleSignUpButtonClick = async () => {
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
          onClick={handleSignUpButtonClick}
        >
          {t("common.sign-up")}
        </Button>
      </div>
    </div>
  );
});

export default SignUp;
