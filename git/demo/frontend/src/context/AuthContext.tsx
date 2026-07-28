import React, { createContext, useContext, useEffect, useState } from 'react';

type AuthContextType = {
  isAuthenticated: boolean;
  isAdmin: boolean;
  username: string | null;
  email: string | null;
  // Real Twitch user ID (the users.twitch_id / session user_id key) - NOT
  // the same as username (the login name). Needed anywhere a caller must
  // match the backend's own id key, e.g. Paddle checkout customData.
  twitchId: string | null;
  checkAuth: () => Promise<void>;
  logout: () => Promise<void>;
  loading: boolean;
  loginError: string | null;
  authChecked: boolean;
};

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isAdmin: false,
  username: null,
  email: null,
  twitchId: null,
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
  const [email, setEmail] = useState<string | null>(null);
  const [twitchId, setTwitchId] = useState<string | null>(null);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);
  const [loading, setLoading] = useState(true);
  const [loginError, setLoginError] = useState<string | null>(null);
  const [authChecked, setAuthChecked] = useState(false);

  const checkAuth = async (showError = false) => {
    setLoading(true);

    if (FAKE_AUTH) {
      setUsername(FAKE_USER);
      setEmail(`${FAKE_USER}@example.com`);
      setTwitchId("dev-twitch-id");
      setIsAdmin(true);
      setAuthChecked(true);
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(`/auth/me`, {
        credentials: "include",
      });

      if (res.status === 401) {
        setUsername(null);
        setEmail(null);
        setTwitchId(null);
        setIsAdmin(false);
        if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
        return;
      }

      if (!res.ok) {
        console.error("Auth check failed:", res.statusText);
        setUsername(null);
        setEmail(null);
        setTwitchId(null);
        if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
        return;
      }

      const data = await res.json();
      setUsername(data.username || data.login);
      setEmail(data.email || null);
      setTwitchId(data.twitch_id || data.id || null);
      setIsAdmin(data.is_admin || data.isAdmin || false);
      setLoginError(null);
    } catch (err) {
      console.error("Login check failed:", err);
      setUsername(null);
      setEmail(null);
      setTwitchId(null);
      setIsAdmin(false);
      if (showError) setLoginError("Das mit dem Login hat nicht geklappt!");
    } finally {
      setAuthChecked(true);
      setLoading(false);
    }
  };

  const logout = async () => {
    if (FAKE_AUTH) {
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
      setEmail(null);
      setTwitchId(null);
      setIsAdmin(false);
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
      value={{ isAdmin, isAuthenticated, username, email, twitchId, checkAuth, logout, loading, loginError, authChecked }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
