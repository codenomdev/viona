import { create } from 'zustand';

import { AdminSettingsInterface } from '@/common/interface';
import { DEFAULT_LANG } from '@/common/constants';

interface InterfaceType {
  interface: AdminSettingsInterface;
  update: (params: AdminSettingsInterface) => void;
}

const interfaceSetting = create<InterfaceType>((set) => ({
  interface: {
    language: DEFAULT_LANG,
    time_zone: '',
    default_avatar: 'system',
    gravatar_base_url: '',
  },
  update: (params) =>
    set((state) => {
      return {
        interface: { ...state.interface, ...params },
      };
    }),
}));

export default interfaceSetting;
