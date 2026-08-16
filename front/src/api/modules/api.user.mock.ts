export default [
  {
    path: "/auth/login",
    method: "post",
    handle() {
      return {
        code: 0,
        msg: "success",
        data: {
          token: "faker token",
          expiresAt: Math.floor(Date.now() / 1000) + 604800,
          user: {
            id: 1,
            username: "admin",
            nickName: "admin",
            role: "admin"
          }
        }
      };
    }
  },
  {
    path: "/auth/register",
    method: "post",
    handle() {
      return {
        code: 0,
        msg: "success",
        data: {
          id: 1,
          username: "username",
          role: "read"
        }
      };
    }
  }
];
