import { NamedExoticComponent, FC, RefObject } from 'react';

import type * as Type from '@/common/interface';
import Request from '@/utils/request';

export enum PluginType {
  Connector = 'connector',
  Search = 'search',
  Editor = 'editor',
  EditorReplacement = 'editor_replacement',
  Route = 'route',
  Captcha = 'captcha',
  Render = 'render',
  Sidebar = 'sidebar',
}

export interface PluginInfo {
  slug_name: string;
  type: PluginType;
  name?: string;
  description?: string;
  route?: string;
  registrationMode?: 'multiple' | 'singleton';
}

export interface Plugin {
  info: PluginInfo;
  component: NamedExoticComponent | FC;
  i18nConfig?;
  hooks?: {
    useRender?: Array<
      (
        element: HTMLElement | RefObject<HTMLElement> | null,
        request?: typeof Request,
      ) => void
    >;
    useCaptcha?: (props: { captchaKey: Type.CaptchaKey; commonProps: any }) => {
      getCaptcha: () => Record<string, any>;
      check: (t: () => void) => void;
      handleCaptchaError: (error) => any;
      close: () => Promise<void>;
      resolveCaptchaReq: (data) => void;
    };
  };
  activated?: boolean;
}
