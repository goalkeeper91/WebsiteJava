import React, { createContext, useContext, useEffect, useState } from 'react';

type AuthContextType = {
  isAuthenticated: boolean;
  username: string | null;
  checkAuth: () => Promise<void>;
  logout: () => Promise<void>;
  loading: boolean;
};

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  username: null,
  checkAuth: async () => {},
  logout: async () => {},
  loading: true,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [username, setUsername] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const checkAuth = async () => {
    setLoading(true);
    try {
      const res = await fetch(`/api/auth/me`, {
        credentials: "include",
      });

      if (res.status === 401) {
        setUsername(null);
        return;
      }

      if (!res.ok) {
        console.error("Auth check failed:", res.statusText);
        setUsername(null);
        return;
      }

      const data = await res.json();
      setUsername(data.username || data.display_name);
    } catch (err) {
      console.error("Login check failed:", err);
      setUsername(null);
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    setLoading(true);
    try {
      await fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });
    } catch (e) {
      console.error("Logout failed:", e);
    } finally {
      setUsername(null);
      setLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const isAuthenticated = !!username;

  return (
    <AuthContext.Provider value={{ isAuthenticated, username, checkAuth, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
