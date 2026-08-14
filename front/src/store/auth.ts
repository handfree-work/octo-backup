import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { notification } from 'antdv-next';
import { defineStore } from 'pinia';

import { loginApi } from '#/api';
import { $t } from '#/locales';

const USER_STORAGE_KEY = 'web-restic-user';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();
  const loginLoading = ref(false);

  async function authLogin(params: Recordable<any>, onSuccess?: () => Promise<void> | void) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const result = await loginApi({
        password: params.password,
        username: params.username,
      });
      userInfo = toUserInfo(result.user, result.token);
      accessStore.setAccessToken(result.token);
      accessStore.setAccessCodes([result.user.role]);
      accessStore.setIsAccessChecked(false);
      userStore.setUserInfo(userInfo);
      localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userInfo));

      if (accessStore.loginExpired) {
        accessStore.setLoginExpired(false);
      } else if (onSuccess) {
        await onSuccess();
      } else {
        await router.push(userInfo.homePath || preferences.app.defaultHomePath);
      }

      notification.success({
        description: `${$t('authentication.loginSuccessDesc')}:${userInfo.realName}`,
        duration: 3,
        title: $t('authentication.loginSuccess'),
      });
    } catch (error: any) {
      notification.error({
        description: error?.error || '登录失败，请检查用户名和密码',
        duration: 3,
        title: '登录失败',
      });
      throw error;
    } finally {
      loginLoading.value = false;
    }

    return { userInfo };
  }

  async function logout(redirect: boolean = true) {
    resetAllStores();
    localStorage.removeItem(USER_STORAGE_KEY);
    accessStore.setLoginExpired(false);

    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  async function fetchUserInfo() {
    const currentUser = userStore.userInfo as UserInfo | null;
    if (currentUser) {
      return currentUser;
    }
    const cachedUser = localStorage.getItem(USER_STORAGE_KEY);
    if (!cachedUser) {
      throw new Error('登录信息不存在');
    }
    const userInfo = JSON.parse(cachedUser) as UserInfo;
    userStore.setUserInfo(userInfo);
    accessStore.setAccessCodes(userInfo.roles ?? []);
    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});

function toUserInfo(
  user: {
    avatar?: string;
    id: number;
    nickName: string;
    role: string;
    username: string;
  },
  token: string,
): UserInfo {
  return {
    avatar: user.avatar || preferences.app.defaultAvatar,
    desc: user.role,
    homePath: preferences.app.defaultHomePath,
    realName: user.nickName || user.username,
    roles: [user.role],
    token,
    userId: String(user.id),
    username: user.username,
  };
}
