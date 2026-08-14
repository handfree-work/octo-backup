import path from 'node:path';

export default {
  resolve: {
    alias: {
      '#': path.resolve(__dirname, 'src'),
    },
  },
  test: {
    include: ['src/api/core/auth.test.ts'],
  },
};
