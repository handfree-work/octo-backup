import { defineStore } from "pinia";
import { store } from "../index";
import router from "../../router";
// @ts-ignore
import { LocalStorage } from "/src/utils/util.storage";
// @ts-ignore
import * as UserApi from "/src/views/framework/auth/api";
// @ts-ignore
import { LoginReq, UserInfoRes } from "/@/views/framework/auth/api";
import { Modal } from "ant-design-vue";
import { useI18n } from "vue-i18n";

import { mitter } from "/src/utils/util.mitt";
import { resetAllStores, useAccessStore } from "/@/vben/stores";

interface UserState {
  userInfo: Nullable<UserInfoRes>;
  token?: string;
}

const USER_INFO_KEY = "USER_INFO";
const TOKEN_KEY = "TOKEN";
export const useUserStore = defineStore({
  id: "app.user",
  state: (): UserState => ({
    // user info
    userInfo: null,
    // token
    token: undefined
  }),
  getters: {
    getUserInfo(): UserInfoRes {
      return this.userInfo || LocalStorage.get(USER_INFO_KEY) || {};
    },
    getToken(): string {
      return this.token || LocalStorage.get(TOKEN_KEY);
    }
  },
  actions: {
    setToken(token: string, expire = 60 * 60 * 24 * 7) {
      this.token = token;
      const accessStore = useAccessStore();
      accessStore.setAccessToken(token);
      LocalStorage.set(TOKEN_KEY, this.token, expire);
    },
    setUserInfo(info: UserInfoRes) {
      this.userInfo = info;
      LocalStorage.set(USER_INFO_KEY, info);
    },
    resetState() {
      this.userInfo = null;
      this.token = "";
      LocalStorage.remove(TOKEN_KEY);
      LocalStorage.remove(USER_INFO_KEY);
    },
    /**
     * @description: login
     */
    async login(params: LoginReq): Promise<any> {
      const data = await UserApi.login(params);
      const expiresIn = data.expiresAt ? Math.max(1, data.expiresAt - Math.floor(Date.now() / 1000)) : undefined;

      this.setToken(data.token, expiresIn);
      this.setUserInfo(data.user);
      const accessStore = useAccessStore();
      accessStore.setAccessCodes([data.user.role]);
      await router.replace("/index");
      mitter.emit("app.login", { userInfo: data.user, token: data.token });
      return data.user;
    },
    async getUserInfoAction(): Promise<UserInfoRes> {
      return this.getUserInfo;
    },
    /**
     * @description: logout
     */
    logout(goLogin = true) {
      this.resetState();
      resetAllStores();
      goLogin && router.push("/login");
      mitter.emit("app.logout");
    },

    /**
     * @description: Confirm before logging out
     */
    confirmLoginOut() {
      const { t } = useI18n();
      Modal.config({
        iconType: "warning",
        title: t("app.login.logoutTip"),
        content: t("app.login.logoutMessage"),
        onOk: async () => {
          await this.logout(true);
        }
      });
    }
  }
});
