export interface ApiMeta {
  success: boolean;
  message: string;
  status_code: number;
}

export interface ApiResponse<T = any> {
  meta: ApiMeta;
  error: boolean;
  data: T;
}

export interface HelmetBase {
  pageTitle?: string;
  description?: string;
  keywords?: string;
}

export interface HelmetUpdate extends Omit<HelmetBase, 'pageTitle'> {
  title?: string;
  subtitle?: string;
}

export type CaptchaKey =
  | 'email'
  | 'password'
  | 'edit_userinfo'
  | 'question'
  | 'codenom'
  | 'comment'
  | 'edit'
  | 'invitation_codenom'
  | 'search'
  | 'report'
  | 'delete'
  | 'vote';

export interface ActivatedPlugin {
  slug_name: string;
  enabled: boolean;
}

export interface UserPluginsConfigRes {
  name: string;
  slug_name: string;
}

export interface LangsType {
  label: string;
  value: string;
}

// /**
//  * @description interface for Admin Settings
//  */
// export interface AdminSettingsGeneral {
//   name: string;
//   short_description: string;
//   description: string;
//   site_url: string;
//   contact_email: string;
//   permalink?: number;
// }

// export interface AdminSettingsLogin {
//   allow_new_registrations: boolean;
//   allow_email_registrations: boolean;
//   allow_email_domains: string[];
//   allow_password_login: boolean;
// }

export interface SettingsGeneral {
  site_name: string;
  language: string;
  timezone: string;
  description: string;
}

export interface SettingsSecurity {
  allow_email_registrations: boolean;
  allow_new_registrations: boolean;
  allow_password_login: boolean;
  allow_user_recover: boolean;
}

export interface SiteSettings {
  // branding: AdminSettingBranding;
  general: SettingsGeneral;
  security: SettingsSecurity;
  // login: AdminSettingsLogin;
  // custom_css_html: AdminSettingsCustom;
  // theme: AdminSettingsTheme;
  // site_seo: AdminSettingsSeo;
  // site_users: AdminSettingsUsers;
  // site_advanced: AdminSettingsWrite;
  // site_questions: AdminQuestionSetting;
  // site_tags: AdminTagsSetting;
  version: string;
  revision: string;
  // site_security: AdminSettingsSecurity;
  ai_enabled: boolean;
}

export interface FormValue<T = any> {
  value: T;
  isInvalid: boolean;
  errorMsg: string;
  [prop: string]: any;
}

export interface FormDataType {
  [prop: string]: FormValue;
}

export interface FieldError {
  error_field: string;
  error_msg: string;
}

export interface LoginReqParams {
  email: string;
  password: string;
}

export interface RegisterReqParams {
  email: string;
  password: string;
}

export interface RecoveryAccountReqParams {
  email: string;
}
