import { FC, memo, useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';

import { SWRConfig } from 'swr';

import { Toaster } from '@/components/ui/sonner';
import { errorCodeStore } from '@/stores';
import HttpErrorContent from '@/components/HttpErrorContent';
import PageTags from '@/components/PageTags';

const Layout: FC = () => {
  const location = useLocation();
  const { code: httpStatusCode, reset: httpStatusReset } = errorCodeStore();

  useEffect(() => {
    httpStatusReset();
  }, [location]);
  return (
    <HelmetProvider>
      <PageTags />
      <SWRConfig value={{ revalidateOnFocus: false }}>
        <div>
          {httpStatusCode ? (
            <HttpErrorContent httpCode={httpStatusCode} />
          ) : (
            <Outlet />
          )}
        </div>
        <Toaster />
      </SWRConfig>
    </HelmetProvider>
  );
};

export default memo(Layout);
