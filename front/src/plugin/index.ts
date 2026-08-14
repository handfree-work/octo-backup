import FastCrud from './fast-crud';

export default {
  install(app: any) {
    app.use(FastCrud);
  },
};
