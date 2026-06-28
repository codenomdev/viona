import { FC, useEffect, useLayoutEffect } from 'react';
import { Helmet } from 'react-helmet-async';

import { REACT_BASE_PATH } from '@/router/alias';
import { pageTagStore, siteInfoStore } from '@/stores';
import { getCurrentLang } from '@/utils/localize';

const doInsertCustomCSS = !document.querySelector('link[href*="custom.css"]');

const Index: FC = () => {
  // const { favicon, square_icon } = brandingStore((state) => state.branding);
  const { pageTitle, keywords, description } = pageTagStore(
    (state) => state.items,
  );
  const appVersion = siteInfoStore((_) => _.version);
  const hashVersion = siteInfoStore((_) => _.revision);
  const siteName = siteInfoStore((_) => _.siteInfo).name;
  const setAppGenerator = () => {
    if (!appVersion) {
      return;
    }
    const generatorMetaNode = document.querySelector('meta[name="generator"]');
    if (generatorMetaNode) {
      generatorMetaNode.setAttribute(
        'content',
        `Answer ${appVersion} - https://github.com/apache/answer version ${hashVersion}`,
      );
    }
  };
  const setDocTitle = () => {
    try {
      if (pageTitle) {
        document.title = pageTitle;
      }
      // eslint-disable-next-line no-empty
    } catch (ex) {}
  };
  const currentLang = getCurrentLang();
  const setDocLang = () => {
    if (currentLang) {
      document.documentElement.setAttribute(
        'lang',
        currentLang.replace('_', '-'),
      );
    }
  };
  // properties used for social media tags
  const openGraphType = 'website';
  const twitterType = 'summary';
  const { href } = window.location;
  const { hostname } = new URL(href);

  useEffect(() => {
    setDocLang();
  }, [currentLang]);
  useEffect(() => {
    setAppGenerator();
  }, [appVersion]);
  useLayoutEffect(() => {
    setDocTitle();
  }, [pageTitle]);
  return (
    <Helmet>
      <link
        rel="icon"
        type="image/png"
        href={`${REACT_BASE_PATH}/favicon.ico`}
      />
      {/* <link rel="icon" type="image/png" sizes="192x192" href={square_icon} />
      <link rel="apple-touch-icon" type="image/png" href={square_icon} /> */}
      <title>{pageTitle}</title>
      {keywords && <meta name="keywords" content={keywords} />}
      {description && <meta name="description" content={description} />}
      {doInsertCustomCSS && (
        <link
          rel="stylesheet"
          href={`${process.env.PUBLIC_URL}${REACT_BASE_PATH}/custom.css`}
        />
      )}
      {/* Social media meta share tags start here */}
      <meta property="og:type" content={openGraphType} />
      <meta property="og:title" name="twitter:title" content={pageTitle} />
      <meta property="og:site_name" content={siteName} />
      <meta property="og:url" content={href} />
      {description && <meta property="og:description" content={description} />}
      <meta
        property="og:image"
        itemProp="image primaryImageOfPage"
        content="/favicon.ico"
      />
      <meta name="twitter:card" content={twitterType} />
      <meta name="twitter:domain" content={hostname} />
      {description && <meta name="twitter:description" content={description} />}
      <meta name="twitter:image" content="/favicon.ico" />
      {/* Social media meta share tags end here */}
    </Helmet>
  );
};

export default Index;
