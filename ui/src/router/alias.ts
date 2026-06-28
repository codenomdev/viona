export const REACT_BASE_PATH = process.env.REACT_APP_BASE_URL || '';
export const BASE_ORIGIN = `${window.location.origin}${REACT_BASE_PATH}`;

export const RouteAlias = {
  home: '/',
  login: '/auth/login',
  signUp: '/auth/register',
  activationFailed: '/failed-activation',
  inactive: '/user/inactive',
  suspended: 'user/suspend',
};
