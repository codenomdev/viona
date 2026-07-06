import { create } from 'zustand';

import { HelmetBase, HelmetUpdate } from '@/common/interface';

import generalInfoStore from './generalInfo';

interface HelmetStore {
  items: HelmetBase;
  update: (params: HelmetUpdate) => void;
}

const makePageTitle = (title = '', subtitle = '') => {
  const { siteInfo } = generalInfoStore.getState();
  if (!subtitle) {
    subtitle = `${siteInfo.site_name}`;
  }
  let pageTitle = subtitle;
  if (title && title !== subtitle) {
    pageTitle = `${title}${subtitle ? ` - ${subtitle}` : ''}`;
  }
  return pageTitle;
};

const pageTags = create<HelmetStore>((set) => ({
  items: {
    pageTitle: '',
    description: '',
    keywords: '',
  },
  update: (params) => {
    const o: HelmetBase = {};
    if (params.title || params.subtitle) {
      o.pageTitle = makePageTitle(params.title, params.subtitle);
    }
    o.description =
      params.description ||
      generalInfoStore.getState().siteInfo?.description ||
      '';
    o.keywords = params.keywords || '';

    set({
      items: o,
    });
  },
}));

export default pageTags;
