import { memo, FC } from 'react';
import { useTranslation } from 'react-i18next';

import {
  getTransNs,
  getTransKeyPrefix,
  PluginInfo,
} from '@/utils/pluginKit/utils';

import { useGetStartUseOauthConnector } from './service';

import info from './info.yaml';
import './i18n';
import { Button } from '@/components/ui/button';
import clsx from 'clsx';

const pluginInfo: PluginInfo = {
  slug_name: info.slug_name,
  type: info.type,
};

interface Props {
  className?: string;
}

const Index: FC<Props> = ({ className }) => {
  const { t } = useTranslation(getTransNs(), {
    keyPrefix: getTransKeyPrefix(pluginInfo),
  });

  const { data } = useGetStartUseOauthConnector();

  if (!data?.length) return null;
  return (
    <div className={clsx('d-grid gap-2', className)}>
      {data?.map((item) => {
        return (
          <Button>
            {/* <SvgIcon base64={item.icon} svgClassName="btnSvg me-2" /> */}
            <span>{t('connect', { auth_name: item.name })}</span>
          </Button>
        );
      })}
    </div>
  );
};

export default {
  info: pluginInfo,
  component: memo(Index),
};
