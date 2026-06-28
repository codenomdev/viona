import { memo, FC } from 'react';
// import { useTranslation } from 'react-i18next';

import {
  // getTransNs,
  // getTransKeyPrefix,
  PluginInfo,
} from '@/utils/pluginKit/utils';

import info from './info.yaml';
import './i18n';

const pluginInfo: PluginInfo = {
  slug_name: info.slug_name,
  type: info.type,
};

interface Props {
  className?: string;
}

const Index: FC<Props> = ({ className }) => {
  return <div className={className}>Third Party Connector</div>;
};

export default {
  info: pluginInfo,
  component: memo(Index),
};
