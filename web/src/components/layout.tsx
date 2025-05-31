import { Outlet } from 'react-router-dom';
import { Navbar } from './navbar';

export const Layout = () => (
  <div>
    <Navbar />
    <div>
      <Outlet />
    </div>
  </div>
);
