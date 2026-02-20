import { Outlet } from "react-router-dom";
import DashboardSidebar from "../../components/dashboard/DashboardSidebar";

const UserDashboard = () => {
  return (
    <div className="flex">
      <DashboardSidebar />

      <main className="flex-grow">
        <Outlet />
      </main>
    </div>
  );
};

export default UserDashboard;