import { create } from 'zustand';

import { SettingsGeneral } from '@/common/interface';
import {
  DEFAULT_LANG,
  DEFAULT_SITE_NAME,
  DEFAULT_TIMEZONE,
} from '@/common/constants';

interface SiteInfoType {
  siteInfo: SettingsGeneral;
  version: string;
  revision: string;
  update: (params: SettingsGeneral) => void;
  updateVersion: (ver: string, revision: string) => void;
  users: any;
  updateUsers: (users: SiteInfoType['users']) => void;
}

const defaultUsersConf: any = {
  allow_update_avatar: false,
  allow_update_bio: false,
  allow_update_display_name: false,
  allow_update_location: false,
  allow_update_username: false,
  allow_update_website: false,
  default_avatar: 'system',
  gravatar_base_url: '',
};

const generalInfo = create<SiteInfoType>((set) => ({
  siteInfo: {
    language: DEFAULT_LANG,
    timezone: DEFAULT_TIMEZONE,
    site_name: DEFAULT_SITE_NAME,
    description: '',
    // name: DEFAULT_SITE_NAME,
    // short_description: '',
    // site_url: '',
    // contact_email: '',
    // permalink: 1,
  },
  users: defaultUsersConf,
  version: '',
  revision: '',
  update: (params) =>
    set((_) => {
      const o = { ..._.siteInfo, ...params };
      if (!o.site_name) {
        o.site_name = DEFAULT_SITE_NAME;
      }
      return {
        siteInfo: o,
      };
    }),
  updateVersion: (ver, revision) => {
    set(() => {
      return { version: ver, revision };
    });
  },
  updateUsers: (users) => {
    set(() => {
      users ||= defaultUsersConf;
      return { users };
    });
  },
}));

export default generalInfo;
