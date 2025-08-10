import { makeAutoObservable } from "mobx";
import { supabaseClient } from "@/store";
import { type User } from '@supabase/supabase-js'

class LocalState {
  currentUser?: string;
  userMapByName: Record<string, User> = {};

  constructor() {
    makeAutoObservable(this);
  }

  setPartial(partial: Partial<LocalState>) {
    Object.assign(this, partial);
  }
}

const userStore = (() => {
  const state = new LocalState();

  return {
    state,
  };
})();

export const initialUserStore = async () => {
  const { data: { user } } = await supabaseClient.auth.getUser()
  if (!user) {
      userStore.state.setPartial({
        currentUser: undefined,
        userMapByName: {},
      });
      return;
  }

  if (user.email) {
    userStore.state.setPartial({
    currentUser: user.email,
    userMapByName: {
      [user.email]: user,
    },
  });
  }
};

export default userStore;

