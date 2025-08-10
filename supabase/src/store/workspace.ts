import { makeAutoObservable } from "mobx";
import { isValidateLocale } from "@/utils/i18n";

class LocalState {
  locale: string = "en";
  appearance: string = "system";

  constructor() {
    makeAutoObservable(this);
  }

  setPartial(partial: Partial<LocalState>) {
    const finalState = {
      ...this,
      ...partial,
    };
    if (!isValidateLocale(finalState.locale)) {
      finalState.locale = "en";
    }
    if (!["system", "light", "dark"].includes(finalState.appearance)) {
      finalState.appearance = "system";
    }
    Object.assign(this, finalState);
  }
}

const workspaceStore = (() => {
  const state = new LocalState();

  return {
    state,
  }
})();

export default workspaceStore;
