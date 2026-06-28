import { RouteAlias } from '@/router/alias';
import { getAppSettings } from '@/services';
import {
  interfaceStore,
  loginSettingStore,
  siteInfoStore,
  siteSecurityStore,
} from '@/stores';

import { floppyNavigation } from './floppyNavigation';
import { setupAppLanguage, setupAppTimeZone } from './localize';
import { pullUcAgent } from './userCenter';

export type TGuardResult = {
  ok: boolean;
  redirect?: string;
  error?: {
    code?: number | string;
    msg?: string;
  };
};

export const IGNORE_PATH_LIST = [
  RouteAlias.login,
  RouteAlias.signUp,
  // RouteAlias.accountRecovery,
  // RouteAlias.changeEmail,
  // RouteAlias.passwordReset,
  // RouteAlias.accountActivation,
  // RouteAlias.confirmNewEmail,
  // RouteAlias.confirmEmail,
  // RouteAlias.authLanding,
  // '/user-center/',
];

export const isIgnoredPath = (ignoredPath?: string | string[]) => {
  if (!ignoredPath) {
    ignoredPath = IGNORE_PATH_LIST;
  }
  if (!Array.isArray(ignoredPath)) {
    ignoredPath = [ignoredPath];
  }
  const matchingPath = ignoredPath.find((p) => {
    return floppyNavigation.matchToCurrentHref(p);
  });
  return !!matchingPath;
};

/**
 * Initialize app configuration
 */
export const initAppSettingsStore = async () => {
  const appSettings = await getAppSettings();
  if (appSettings) {
    siteInfoStore.getState().update(appSettings.general);
    siteInfoStore
      .getState()
      .updateVersion(appSettings.version, appSettings.revision);
    // siteInfoStore.getState().updateUsers(appSettings.site_users);
    interfaceStore.getState().update(appSettings.interface);
    // pageTagStore.getState().update({
    //   title: appSettings.general?.name,
    //   description: appSettings.general?.description,
    // });
    // brandingStore.getState().update(appSettings.branding);
    loginSettingStore.getState().update(appSettings.login);
    // customizeStore.getState().update(appSettings.custom_css_html);
    // themeSettingStore.getState().update(appSettings.theme);
    // seoSettingStore.getState().update(appSettings.site_seo);
    // // writeSettingStore.getState().update({
    //   ...appSettings.site_advanced,
    //   ...appSettings.site_questions,
    //   ...appSettings.site_tags,
    // });
    // aiControlStore.getState().update({
    //   ai_enabled: appSettings.ai_enabled,
    // });
    siteSecurityStore.getState().update(appSettings.site_security);
  }
};

let appInitialized = false;
export const setupApp = async () => {
  /**
   * This cannot be removed:
   * clicking on the current navigation link will trigger a call to the routing loader,
   * even though the page is not refreshed.
   */
  if (appInitialized) {
    return;
  }
  /**
   * WARN:
   * 1. must pre init logged user info for router guard
   * 2. must pre init app settings for app render
   */
  await Promise.allSettled([initAppSettingsStore()]);
  await Promise.allSettled([pullUcAgent()]);
  setupAppLanguage();
  setupAppTimeZone();
  // setupAppTheme();
  /**
   * WARN:
   * Initialization must be completed after all initialization actions,
   * otherwise the problem of rendering twice in React development mode can lead to inaccurate data or flickering pages
   */
  appInitialized = true;
};

export const shouldLoginRequired = () => {
  const gr: TGuardResult = { ok: true };
  // const { login_required } = siteSecurityStore.getState();
  // if (!login_required) {
  //   return gr;
  // }
  // const us = deriveLoginState();
  // if (us.isLogged) {
  //   return gr;
  // }
  if (isIgnoredPath(IGNORE_PATH_LIST)) {
    return gr;
  }
  gr.ok = false;
  gr.redirect = RouteAlias.login;
  return gr;
};

export const notLogged = () => {
  const gr: TGuardResult = { ok: true };
  // const us = deriveLoginState();
  // if (us.isLogged) {
  //   gr.ok = false;
  //   gr.redirect = RouteAlias.home;
  // }
  return gr;
};

export const notActivated = () => {
  const gr: TGuardResult = { ok: true };
  // const us = deriveLoginState();
  // if (us.isActivated) {
  //   gr.ok = false;
  //   gr.redirect = RouteAlias.home;
  // }
  return gr;
};

export const activated = () => {
  // const gr = logged();
  // const us = deriveLoginState();
  // if (us.isNotActivated) {
  //   gr.ok = false;
  //   gr.redirect = RouteAlias.inactive;
  // }
  // return gr;
};

export type TGuardFunc = (args: {
  loaderData?: any;
  path?: string;
  page?: string;
}) => TGuardResult;
