import { request } from '@/utils';
import type * as Type from '@/common/interface';

export const getLanguageConfig = () => {
  return request.get('/api/v1/language/config');
};

export const getLanguageOptions = () => {
  return request.get<Type.LangsType[]>('/api/v1/language/options');
};

export const getPluginsStatus = () => {
  return request.get<Type.ActivatedPlugin[]>('/api/v1/plugin/status');
};

export const getAppSettings = () => {
  return request.get<Type.SiteSettings>('/api/v1/setting/get');
};

export const login = (params: Type.LoginReqParams) => {
  return request.post<any>('/api/v1/auth/login', params);
};
