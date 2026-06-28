import { create } from 'zustand';

// import type { UserInfoRes } from '@/common/interface';
import Storage from '@/utils/storage';
import { LOGGED_TOKEN_STORAGE_KEY } from '@/common/constants';

// interface UserInfoStore {
//   user: UserInfoRes;
//   update: (params: UserInfoRes) => void;
//   clear: (removeToken?: boolean) => void;
// }

// const initUser: UserInfoRes = {
//   access_token: '',
//   username: '',
//   avatar: '',
//   rank: 0,
//   bio: '',
//   bio_html: '',
//   display_name: '',
//   location: '',
//   website: '',
//   status: 'normal',
//   mail_status: 1,
//   language: 'Default',
//   color_scheme: 'default',
//   is_admin: false,
//   have_password: true,
//   role_id: 1,
// };

const initUser = {
  language: 'Default',
  color_scheme: 'default',
  access_token: '',
};

const loggedUserInfo = create<any>((set) => ({
  user: initUser,
  update: (params) => {
    if (typeof params !== 'object' || !params) {
      return;
    }
    if (!params?.language) {
      params.language = 'Default';
    }
    if (!params?.color_scheme) {
      params.color_scheme = 'default';
    }
    set(() => {
      Storage.set(LOGGED_TOKEN_STORAGE_KEY, params.access_token);
      return { user: params };
    });
  },
  clear: (removeToken = true) =>
    set(() => {
      if (removeToken) {
        Storage.remove(LOGGED_TOKEN_STORAGE_KEY);
      }
      return { user: initUser };
    }),
}));

export default loggedUserInfo;
