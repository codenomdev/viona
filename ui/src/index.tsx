import React from 'react';

import ReactDOM from 'react-dom/client';

import App from './App';

// @ts-ignore: allow side-effect import of global CSS without type declarations
import './globals.css';

/**
 *Automatically jump when the href of a Link component within a matching project is not a front-end route.
 *
 */
const handleClickLink = (evt: Event) => {
  const { target } = evt;

  if (target === null || !(target instanceof Element)) {
    return;
  }
  if (!/A/i.test(target.nodeName)) {
    return;
  }

  if (target.getAttribute('href')?.includes('/api/')) {
    evt.preventDefault();
    window.location.href = target.getAttribute('href') || '';
  }
};

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement,
);

document.addEventListener('click', handleClickLink, true);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
