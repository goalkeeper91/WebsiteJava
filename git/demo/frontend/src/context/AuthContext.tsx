import React, { createContext, useContext, useEffect, useState } from 'react';

type AuthContextType = {
  isAuthenticated: boolean;
  username: string | null;
  login: () => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  username: null,
  login: async () => {},
  logout: async () => {},
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [username, setUsername] = useState<string | null>(null);

  const login = async () => {
    try {
        const res = await fetch("http://localhost:8080/api/auth/me", {
            credentials: "include",
            });
            if (!res.ok) throw new Error("Not logged in");
            const data = await  res.json();
            setUsername(data.username);
    } catch (error) {
        setUsername(null);
    }
  };

  const logout = async () => {
    try {
        await fetch("http://localhost:8080/api/auth/logout", {
            method: "POST",
            credentials: "include",
            });
    } finally {
        setUsername(null);
    }
  };

  useEffect(() => {
      login();
  }, []);

  const isAuthenticated = !!username;

  return (
    <AuthContext.Provider value={{ isAuthenticated, username, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
    if (import.meta.env.VITE_DEV_AUTH === 'true') {
        return {
          isAuthenticated: true,
          username: import.meta.env.VITE_DEV_USERNAME || 'devuser',
          logout: async () => console.log('Mock Logout'),
          login: async () => console.log('Mock Login'),
        };
    }
    useContext(AuthContext);
    }
