import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import './i18n/init';

import '@/utils/pluginKit';
import { useMergeRoutes } from '@/router';

function App() {
  const routes = useMergeRoutes();
  if (routes.length === 0) {
    return <div>Loading routes...</div>;
  }
  // console.log(routes);
  const router = createBrowserRouter(routes, {
    basename: process.env.REACT_APP_BASE_URL,
  });
  return <RouterProvider router={router} />;
}

export default App;
