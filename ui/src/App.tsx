import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import './i18n/init';

import '@/utils/pluginKit';
import { useMergeRoutes } from '@/router';
import LoadingPage from '@/components/LoadingPage';

function App() {
  const routes = useMergeRoutes();
  if (routes.length === 0) {
    return <LoadingPage />;
  }
  // console.log(routes);
  const router = createBrowserRouter(routes, {
    basename: process.env.REACT_APP_BASE_URL,
  });
  return <RouterProvider router={router} />;
}

export default App;
