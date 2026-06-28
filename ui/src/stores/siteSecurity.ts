import { create } from 'zustand';

interface SecurityStore {
  login_required: boolean;
  check_update: boolean;
  external_content_display: string;
  update: (params: {
    external_content_display: string;
    check_update: boolean;
    login_required: boolean;
  }) => void;
}

const siteSecurityStore = create<SecurityStore>((set) => ({
  login_required: false,
  check_update: true,
  external_content_display: 'always_display',
  update: (params) =>
    set((state) => {
      return {
        ...state,
        ...params,
      };
    }),
}));

export default siteSecurityStore;
