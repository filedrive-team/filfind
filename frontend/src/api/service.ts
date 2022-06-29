import axios, { AxiosRequestConfig, Method } from 'axios';
import { Cookies } from '@/data/cookies';

interface PendingType {
  url: string | undefined;
  method: Method | undefined;
  params: object;
  data: object;
  cancel: Function;
}

const pending: Array<PendingType> = [];
const CancelToken = axios.CancelToken;

const service = axios.create({
  withCredentials: false,
  timeout: 35000,
});

const removePending = (config: AxiosRequestConfig) => {
  for (const key in pending) {
    const item: number = +key;
    const list: PendingType = pending[key];
    if (
      list.url === config.url &&
      list.method === config.method &&
      JSON.stringify(list.params) === JSON.stringify(config.params) &&
      JSON.stringify(list.data) === JSON.stringify(config.data)
    ) {
      pending.splice(item, 1);
    }
  }
};

service.interceptors.request.use(
  (request) => {
    removePending(request);
    request.cancelToken = new CancelToken((c) => {
      pending.push({
        url: request.url,
        method: request.method,
        params: request.params,
        data: request.data,
        cancel: c,
      });
    });
    request.baseURL = process.env['REACT_APP_BASE_URL'];
    let headers = request.headers;
    const authorization = Cookies.getAuthorization();
    request.headers = { ...headers, Authorization: authorization };
    return request;
  },
  (error) => {},
);

service.interceptors.response.use(
  (response) => {
    removePending(response.config);
    // response.data = JSON.parse(response.data);
    return response;
  },
  (error) => {
    return error;
  },
);

export default service;
