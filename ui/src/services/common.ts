import { request } from '@/utils';
import type * as Type from '@/common/interface';

export const getLanguageConfig = () => {
  return request.get('/api/v1/language/get');
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

export const register = (params: Type.RegisterReqParams) => {
  return request.post<any>('/api/v1/auth/register', params);
};

export const recoveryAccount = (params: Type.RecoveryAccountReqParams) => {
  return request.post<any>('/api/v1/auth/recovery-account', params);
};
