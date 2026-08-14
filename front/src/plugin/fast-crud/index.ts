import type {
  CrudOptions,
  PageQuery,
  PageRes,
  TransformResProps,
  UserPageQuery,
} from '@fast-crud/fast-crud';
import { FastCrud } from '@fast-crud/fast-crud';
import '@fast-crud/fast-crud/dist/style.css';
import UiAntdvNext from '@fast-crud/ui-antdv-next';
import '@fast-crud/ui-antdv-next/dist/style.css';

function install(app: any) {
  app.use(UiAntdvNext);
  app.use(FastCrud, {
    commonOptions(): CrudOptions {
      return {
        table: {
          size: 'small',
          pagination: false,
          scroll: {
            x: 700,
          },
        },
        rowHandle: {
          fixed: 'right',
          width: 200,
          buttons: {
            view: { show: true ,icon: 'EyeOutlined', text: '', type: 'link' },
            copy: { show: true ,icon: 'CopyOutlined', text: '', type: 'link' },
            edit: { icon: 'EditOutlined', text: '', type: 'link' },
            remove: {
              danger: true,
              icon: 'DeleteOutlined',
              text: '',
              type: 'link',
            },
          },
          dropdown: {
            more: {
              type: 'link',
            },
          },
        },
        request: {
          transformQuery: ({ page, form, sort }: PageQuery): UserPageQuery => {
            const limit = page?.pageSize ?? 10;
            const currentPage = page?.currentPage ?? 1;
            const offset = limit * (currentPage - 1);

            return {
              page: {
                limit,
                offset,
              },
              query: form,
              sort: sort ?? {},
            };
          },
          transformRes: ({ res }: TransformResProps): PageRes => {
            const pageSize = res.limit;
            const currentPage =
              pageSize && res.offset % pageSize === 0
                ? res.offset / pageSize + 1
                : 1;

            return {
              ...res,
              currentPage,
              pageSize,
              records: res.records,
              total: res.total,
            };
          },
        },
      };
    },
  });
}

export default {
  install,
};
