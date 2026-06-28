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
            <div className="flex flex-row">
              <main className="view flex-col line-height-[initial] flex w-full flex-wrap min-w-307.5">
                <div className="relative flex w-full flex-wrap justify-center">
                  <section className="relative flex shrink-0 flex-col w-full text-left">
                    <div className="mt-0 justify-items-start pt-5 md:grid md:min-h-[calc(100vh-66px)] md:grid-rows-[max-content] md:justify-center md:items-start md:pt-8 md:pb-2">
                      <Outlet />
                    </div>
                  </section>
                </div>
              </main>
            </div>
          )}
        </div>
        <Toaster />
      </SWRConfig>
    </HelmetProvider>
  );
};

export default memo(Layout);
