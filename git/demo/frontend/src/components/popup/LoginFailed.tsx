import React from "react";
import { useAuth } from "../../context/AuthContext";

export const LoginPopup: React.FC = () => {
  const { loginError, authChecked } = useAuth();

  if (!loginError || !authChecked) return null;

  return (
    <div className="fixed top-4 right-4 bg-red-500 text-white px-4 py-2 rounded shadow z-50">
      {loginError}
    </div>
  );
};
