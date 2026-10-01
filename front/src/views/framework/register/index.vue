<template>
  <div class="main">
    <a-form ref="formRef" class="user-layout-login" name="custom-validation" :model="formState" :rules="rules" v-bind="layout" @finish="handleFinish" @finish-failed="handleFinishFailed">
      <a-form-item required has-feedback name="username">
        <a-input v-model:value="formState.username" size="large" autocomplete="off">
          <template #prefix>
            <span class="iconify" data-icon="ion:person" data-inline="false"></span>
          </template>
        </a-input>
      </a-form-item>
      <a-form-item has-feedback name="password">
        <a-input-password v-model:value="formState.password" size="large" autocomplete="off">
          <template #prefix>
            <span class="iconify" data-icon="ion:lock-closed" data-inline="false"></span>
          </template>
        </a-input-password>
      </a-form-item>
      <a-form-item has-feedback name="confirmPassword">
        <a-input-password v-model:value="formState.confirmPassword" size="large" autocomplete="off">
          <template #prefix>
            <span class="iconify" data-icon="ion:lock-closed" data-inline="false"></span>
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" size="large" html-type="submit" :loading="loading" class="login-button">注册</a-button>
      </a-form-item>

      <a-form-item class="user-login-other">
        <router-link class="register" :to="{ name: 'login' }"> 登录 </router-link>
      </a-form-item>
    </a-form>
  </div>
</template>
<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import * as UserApi from "/@/views/framework/auth/api";
import router from "/@/router";
import { utils } from "@fast-crud/fast-crud";
import { message } from "ant-design-vue";
export default defineComponent({
  name: "RegisterPage",
  setup() {
    const loading = ref(false);
    const formRef = ref();
    const formState = reactive({
      username: "",
      password: "",
      confirmPassword: ""
    });

    const rules = {
      username: [
        {
          required: true,
          trigger: "change",
          message: "请输入用户名"
        }
      ],
      password: [
        {
          required: true,
          trigger: "change",
          message: "请输入密码"
        }
      ],
      confirmPassword: [
        {
          required: true,
          trigger: "change",
          message: "请确认密码"
        },
        {
          validator: async (_rule: unknown, value: string) => {
            if (value && value !== formState.password) {
              throw new Error("两次输入的密码不一致");
            }
          },
          trigger: "change"
        }
      ]
    };
    const layout = {
      labelCol: {
        span: 0
      },
      wrapperCol: {
        span: 24
      }
    };

    const handleFinish = async () => {
      loading.value = true;
      try {
        await UserApi.register({ username: formState.username, password: formState.password });
        message.success("注册成功，请登录");
        await router.replace({ name: "login" });
      } catch (error: any) {
        message.error(error?.message || "注册失败，请稍后重试");
      } finally {
        loading.value = false;
      }
    };

    const handleFinishFailed = (errors: any) => {
      utils.logger.log(errors);
    };

    const resetForm = () => {
      formRef.value.resetFields();
    };

    return {
      formState,
      loading,
      formRef,
      rules,
      layout,
      handleFinishFailed,
      handleFinish,
      resetForm
    };
  }
});
</script>

<style lang="less">
@import "../../../style/theme/index.less";
.user-layout-login {
  label {
    font-size: 14px;
  }

  .login-title {
    color: @primary-color;
    font-size: 18px;
    text-align: center;
    margin: 20px;
  }
  .getCaptcha {
    display: block;
    width: 100%;
    height: 40px;
  }

  .forge-password {
    font-size: 14px;
  }

  button.login-button {
    padding: 0 15px;
    font-size: 16px;
    height: 40px;
    width: 100%;
  }

  .user-login-other {
    text-align: left;
    margin-top: 30px;
    margin-bottom: 30px;
    line-height: 22px;

    .item-icon {
      font-size: 24px;
      color: rgba(0, 0, 0, 0.2);
      margin-left: 16px;
      vertical-align: middle;
      cursor: pointer;
      transition: color 0.3s;

      &:hover {
        color: @primary-color;
      }
    }

    .register {
      float: right;
    }
  }
  .iconify {
    color: rgba(0, 0, 0, 0.45);
  }
}
</style>
