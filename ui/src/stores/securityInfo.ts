import { create } from 'zustand';

import { SettingsSecurity } from '@/common/interface';

interface InterfaceType {
  interface: SettingsSecurity;
  update: (params: SettingsSecurity) => void;
}

const securityInfo = create<InterfaceType>((set) => ({
  interface: {
    allow_email_registrations: true,
    allow_new_registrations: true,
    allow_password_login: true,
    allow_user_recover: true,
  },
  update: (params) =>
    set((state) => {
      return {
        interface: { ...state.interface, ...params },
      };
    }),
}));

export default securityInfo;
