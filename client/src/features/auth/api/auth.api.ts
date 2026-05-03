const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

export interface User {
  id: string;
  email: string;
  createdAt?: string;
}

export interface AuthCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

interface ApiErrorBody {
  error?: string;
}

const parseApiError = async (res: Response, fallback: string) => {
  try {
    const body = (await res.json()) as ApiErrorBody;
    return body.error || fallback;
  } catch {
    return fallback;
  }
};

export const register = async (credentials: AuthCredentials) => {
  const res = await fetch(`${BASE_URL}/api/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  if (!res.ok) {
    throw new Error(await parseApiError(res, "Failed to register"));
  }

  return res.json() as Promise<AuthResponse>;
};

export const login = async (credentials: AuthCredentials) => {
  const res = await fetch(`${BASE_URL}/api/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  if (!res.ok) {
    throw new Error(await parseApiError(res, "Failed to login"));
  }

  return res.json() as Promise<AuthResponse>;
};

export const getCurrentUser = async (token: string) => {
  const res = await fetch(`${BASE_URL}/api/auth/me`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error(await parseApiError(res, "Session expired"));
  }

  return res.json() as Promise<User>;
};
