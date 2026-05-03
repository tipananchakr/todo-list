/* eslint-disable react-hooks/set-state-in-effect, react-refresh/only-export-components */

import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { getCurrentUser, login, register, type AuthCredentials, type User } from "../api/auth.api";

const TOKEN_STORAGE_KEY = "todo_access_token";

interface AuthContextValue {
  user: User | null;
  token: string | null;
  isCheckingSession: boolean;
  login: (credentials: AuthCredentials) => Promise<void>;
  register: (credentials: AuthCredentials) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient();
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_STORAGE_KEY));
  const [user, setUser] = useState<User | null>(null);
  const [isCheckingSession, setIsCheckingSession] = useState(Boolean(token));

  const clearSession = useCallback(() => {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    setToken(null);
    setUser(null);
    queryClient.removeQueries({ queryKey: ["todos"] });
  }, [queryClient]);

  const persistSession = useCallback(
    async (nextToken: string, nextUser: User) => {
      localStorage.setItem(TOKEN_STORAGE_KEY, nextToken);
      setToken(nextToken);
      setUser(nextUser);
      await queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
    [queryClient],
  );

  useEffect(() => {
    if (!token) {
      return;
    }

    let active = true;
    setIsCheckingSession(true);

    getCurrentUser(token)
      .then((currentUser) => {
        if (active) {
          setUser(currentUser);
        }
      })
      .catch(() => {
        if (active) {
          clearSession();
        }
      })
      .finally(() => {
        if (active) {
          setIsCheckingSession(false);
        }
      });

    return () => {
      active = false;
    };
  }, [clearSession, token]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isCheckingSession,
      login: async (credentials) => {
        const result = await login(credentials);
        await persistSession(result.token, result.user);
      },
      register: async (credentials) => {
        const result = await register(credentials);
        await persistSession(result.token, result.user);
      },
      logout: clearSession,
    }),
    [clearSession, isCheckingSession, persistSession, token, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
};
