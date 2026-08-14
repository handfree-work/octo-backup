<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, h, ref } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationRegister, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'antdv-next';

import { registerApi } from '#/api';

defineOptions({ name: 'Register' });

const loading = ref(false);
const router = useRouter();

const formSchema = computed((): VbenFormSchema[] => [
  {
    component: 'VbenInput',
    componentProps: {
      autocomplete: 'username',
      placeholder: $t('authentication.usernameTip'),
    },
    fieldName: 'username',
    label: $t('authentication.username'),
    rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
  },
  {
    component: 'VbenInput',
    componentProps: {
      autocomplete: 'nickname',
      placeholder: '显示名称（可选）',
    },
    fieldName: 'nickName',
    label: '显示名称',
  },
  {
    component: 'VbenInputPassword',
    componentProps: {
      autocomplete: 'new-password',
      passwordStrength: true,
      placeholder: $t('authentication.password'),
    },
    fieldName: 'password',
    label: $t('authentication.password'),
    renderComponentContent() {
      return {
        strengthText: () => $t('authentication.passwordStrength'),
      };
    },
    rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
  },
  {
    component: 'VbenInputPassword',
    componentProps: {
      autocomplete: 'new-password',
      placeholder: $t('authentication.confirmPassword'),
    },
    dependencies: {
      rules(values) {
        return z
          .string({ required_error: $t('authentication.passwordTip') })
          .min(1, { message: $t('authentication.passwordTip') })
          .refine((value) => value === values.password, {
            message: $t('authentication.confirmPasswordTip'),
          });
      },
      triggerFields: ['password'],
    },
    fieldName: 'confirmPassword',
    label: $t('authentication.confirmPassword'),
  },
  {
    component: 'VbenCheckbox',
    fieldName: 'agreePolicy',
    renderComponentContent: () => ({
      default: () =>
        h('span', [
          $t('authentication.agree'),
          h(
            'a',
            {
              class: 'vben-link ml-1',
              href: '',
            },
            `${$t('authentication.privacyPolicy')} & ${$t('authentication.terms')}`,
          ),
        ]),
    }),
    rules: z.boolean().refine((value) => !!value, {
      message: $t('authentication.agreeTip'),
    }),
  },
]);

async function handleSubmit(value: Recordable<any>) {
  try {
    loading.value = true;
    await registerApi({
      nickName: value.nickName || undefined,
      password: value.password,
      username: value.username,
    });
    message.success('注册成功，请登录');
    await router.push({ name: 'Login', query: { username: value.username } });
  } catch (error: any) {
    message.error(error?.error || '注册失败，请稍后重试');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <AuthenticationRegister :form-schema="formSchema" :loading="loading" @submit="handleSubmit" />
</template>
