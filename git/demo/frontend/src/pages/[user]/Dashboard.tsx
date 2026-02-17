import { Outlet } from "react-router-dom";

const UserDashboard = () => {
  return (
    <div className="flex">
      <nav>Sidebar mit Links zu /dashboard und /dashboard/n8n</nav>

      <main className="flex-grow">
        <Outlet />
      </main>
    </div>
  );
};

export default UserDashboard;