import { requestClient } from '#/api/request';

const apiPrefix = '/mock/ComponentSelect';

export function GetList(query: any) {
  return requestClient.post(`${apiPrefix}/page`, query);
}

export function AddObj(obj: any) {
  return requestClient.post(`${apiPrefix}/add`, obj);
}

export function UpdateObj(obj: any) {
  return requestClient.post(`${apiPrefix}/update`, obj);
}

export function DelObj(id: any) {
  return requestClient.post(`${apiPrefix}/delete`, undefined, {
    params: { id },
  });
}

export function GetObj(id: any) {
  return requestClient.get(`${apiPrefix}/get`, { params: { id } });
}
