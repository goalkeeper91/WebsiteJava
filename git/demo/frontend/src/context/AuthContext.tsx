import React, { createContext, useContext, useEffect, useState } from 'react';

type AuthContextType = {
  isAuthenticated: boolean;
  username: string | null;
  checkAuth: () => Promise<void>;
  logout: () => Promise<void>;
  loading: boolean;
  loginError: string | null;
  authChecked: boolean;
};

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  username: null,
  checkAuth: async () => {},
  logout: async () => {},
  loading: true,
  loginError: null,
  authChecked: false,
});

// === Neue ENV-Schalter ===
const FAKE_AUTH = import.meta.env.VITE_FAKE_AUTH === "true";
const FAKE_USER = import.meta.env.VITE_FAKE_USER || "dev-user";

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [username, setUsername] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [loginError, setLoginError] = useState<string | null>(null);
  const [authChecked, setAuthChecked] = useState(false);

  const checkAuth = async (showError = false) => {
    setLoading(true);

    if (FAKE_AUTH) {
      // In Dev/Test-Umgebung: immer als eingeloggt behandeln
      setUsername(FAKE_USER);
      setAuthChecked(true);
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(`/api/auth/me`, {
        credentials: "include",
      });

      if (res.status === 401) {
        setUsername(null);
        if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
        return;
      }

      if (!res.ok) {
        console.error("Auth check failed:", res.statusText);
        setUsername(null);
        if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
        return;
      }

      const data = await res.json();
      setUsername(data.username || data.display_name);
      setLoginError(null); // Reset error
    } catch (err) {
      console.error("Login check failed:", err);
      setUsername(null);
      if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
    } finally {
      setAuthChecked(true);
      setLoading(false);
    }
  };

  const logout = async () => {
    if (FAKE_AUTH) {
      // In Fake-Mode ignorieren wir Logout und bleiben eingeloggt
      return;
    }

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
    if (!loginError) return;

    const timer = setTimeout(() => {
        setLoginError(null);
    }, 4000); // 4 Sekunden

    return () => clearTimeout(timer);
  }, [loginError]);

  useEffect(() => {
    checkAuth(false);
  }, []);

  const isAuthenticated = !!username;

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, username, checkAuth, logout, loading, loginError, authChecked }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
