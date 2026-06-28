import axios, { AxiosResponse } from 'axios';
import type {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  InternalAxiosRequestConfig,
} from 'axios';
import { toast } from 'sonner';

import { LOGGED_TOKEN_STORAGE_KEY } from '@/common/constants';
import { RouteAlias } from '@/router/alias';
import { errorCodeStore, loggedUserInfoStore } from '@/stores';
import { getCurrentLang } from '@/utils/localize';
import { ApiResponse } from '@/common/interface';

import { floppyNavigation } from './floppyNavigation';
import { isIgnoredPath, IGNORE_PATH_LIST } from './guard';
import Storage from './storage';

const baseConfig = {
  baseURL:
    process.env.NODE_ENV === 'development' ? '' : process.env.REACT_APP_API_URL,
  timeout: 10000,
  withCredentials: true,
};

interface ApiConfig extends AxiosRequestConfig {
  allow404?: boolean;
  ignoreError?: '403' | '50X';
  passingError?: boolean;
}

class Request {
  instance: AxiosInstance;

  constructor(config: AxiosRequestConfig) {
    this.instance = axios.create(config);

    this.instance.interceptors.request.use(
      (requestConfig: InternalAxiosRequestConfig) => {
        const token = Storage.get(LOGGED_TOKEN_STORAGE_KEY) || '';
        const lang = getCurrentLang();

        requestConfig.headers?.set('Authorization', token);
        requestConfig.headers.set('Accept-Language', lang);

        return requestConfig;
      },
      (err: AxiosError) => {
        console.error('request interceptors error:', err);
      },
    );

    this.instance.interceptors.response.use(
      (res: AxiosResponse<ApiResponse>) => {
        const { meta, data, error } = res.data;

        if (meta.status_code === 204) {
          return true;
        }

        if (error) {
          return Promise.reject({
            code: meta.status_code,
            msg: meta.message,
            data,
          });
        }

        return data;
      },
      (error) => {
        const {
          status,
          data: errBody,
          config: errConfig,
        } = error.response || {};

        const response = errBody as ApiResponse;
        const data = response?.data || {};
        const msg = response?.meta?.message || '';

        const errorObject: {
          code: any;
          msg: string;
          data: any;
          isError?: boolean;
          list?: any[];
        } = {
          code: status,
          msg,
          data,
        };

        if (status === 400) {
          if (data?.err_type && errConfig?.passingError) {
            return Promise.reject(errorObject);
          }

          if (data?.err_type) {
            if (data.err_type === 'toast') {
              toast.error(msg);
            }

            if (data.err_type === 'alert') {
              return Promise.reject({
                msg,
                ...data,
              });
            }

            if (data.err_type === 'modal') {
              console.warn('Modal: ', msg);
            }

            return Promise.reject(false);
          }

          if (Array.isArray(data) && data.length > 0) {
            errorObject.isError = true;
            errorObject.list = data;
            return Promise.reject(errorObject);
          }

          if (!data || Object.keys(data).length <= 0) {
            console.warn('Modal: ', msg);
            return Promise.reject(false);
          }
        }

        if (status === 401) {
          errorCodeStore.getState().reset();
          loggedUserInfoStore.getState().clear();
          floppyNavigation.navigateToLogin();

          return Promise.reject(false);
        }

        if (status === 403) {
          if (data?.type === 'url_expired') {
            floppyNavigation.navigate(RouteAlias.activationFailed, {
              handler: 'replace',
            });

            return Promise.reject(false);
          }

          if (data?.type === 'inactive') {
            floppyNavigation.navigate(RouteAlias.inactive);

            return Promise.reject(false);
          }

          if (data?.type === 'suspended') {
            loggedUserInfoStore.getState().clear();

            floppyNavigation.navigate(RouteAlias.suspended, {
              handler: 'replace',
            });

            return Promise.reject(false);
          }

          if (isIgnoredPath(IGNORE_PATH_LIST)) {
            return Promise.reject(false);
          }

          if (error.config?.url?.includes('/admin/api')) {
            errorCodeStore.getState().update('403');
            return Promise.reject(false);
          }

          if (msg) {
            toast.error(msg);
          }

          return Promise.reject(false);
        }

        if (status === 404 && error.config?.allow404) {
          if (isIgnoredPath(IGNORE_PATH_LIST)) {
            return Promise.reject(false);
          }

          errorCodeStore.getState().update('404');

          return Promise.reject(false);
        }

        if (status && status >= 500) {
          if (isIgnoredPath(IGNORE_PATH_LIST)) {
            return Promise.reject(false);
          }

          if (error.config?.ignoreError !== '50X') {
            errorCodeStore.getState().update('50X');
          }

          toast.error(msg || 'Internal server error');

          console.error(
            `Request failed with status code ${status}, ${msg || ''}`,
          );
        }

        return Promise.reject(errorObject);
      },
    );
  }

  public request<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.instance.request<any, T>(config);
  }

  public get<T = any>(url: string, config?: ApiConfig): Promise<T> {
    return this.instance.get<any, T>(url, config);
  }

  public post<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T> {
    return this.instance.post(url, data, config);
  }

  public put<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T> {
    return this.instance.put(url, data, config);
  }

  public delete<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T> {
    return this.instance.delete(url, {
      data,
      ...config,
    });
  }
}

export default new Request(baseConfig);
