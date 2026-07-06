import { create } from 'zustand';

import { SettingsSecurity } from '@/common/interface';

interface IType {
  login: SettingsSecurity;
  update: (params: SettingsSecurity) => void;
}

const loginSetting = create<IType>((set) => ({
  login: {
    allow_new_registrations: true,
    allow_email_registrations: true,
    // allow_email_domains: [],
    allow_password_login: true,
    allow_user_recover: true,
  },
  update: (params) =>
    set(() => {
      return {
        login: params,
      };
    }),
}));

export default loginSetting;
