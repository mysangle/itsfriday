import { makeAutoObservable } from "mobx";

class LocalState {
  currentUser?: string;

  constructor() {
    makeAutoObservable(this);
  }

  setPartial(partial: Partial<LocalState>) {
    Object.assign(this, partial);
  }
};

const userStore = (() => {
  const state = new LocalState();
});

export default userStore;
