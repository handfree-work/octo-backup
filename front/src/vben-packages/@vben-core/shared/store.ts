import { useStore } from '@tanstack/vue-store';

export * from '@tanstack/vue-store';

export function useSelector<TStore extends Parameters<typeof useStore>[0]>(
  store: TStore,
  selector?: Parameters<typeof useStore>[1],
) {
  return useStore(store, selector);
}
