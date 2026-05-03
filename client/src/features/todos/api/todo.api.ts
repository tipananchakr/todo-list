const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

export interface Todo {
  id: string;
  userId: string;
  completed: boolean;
  body: string;
}

const authHeaders = (token: string) => ({
  Authorization: `Bearer ${token}`,
});

export const getTodos = async (token: string) => {
  const res = await fetch(`${BASE_URL}/api/todos`, {
    headers: authHeaders(token),
  });
  if (!res.ok) {
    throw new Error("Failed to fetch todos");
  }
  return res.json() as Promise<Todo[]>;
};

export const createTodo = async ({ token, body }: { token: string; body: string }) => {
  const res = await fetch(`${BASE_URL}/api/todos`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(token),
    },
    body: JSON.stringify({ body }),
  });

  if (!res.ok) throw new Error("failed to create todo");

  return res.json() as Promise<Todo>;
};

export const completeTodo = async ({ token, id }: { token: string; id: string }) => {
  const res = await fetch(`${BASE_URL}/api/todos/${id}`, {
    method: "PATCH",
    headers: authHeaders(token),
  });

  if (!res.ok) throw new Error("failed to complete todo");
};

export const deleteTodo = async ({ token, id }: { token: string; id: string }) => {
  const res = await fetch(`${BASE_URL}/api/todos/${id}`, {
    method: "DELETE",
    headers: authHeaders(token),
  });

  if (!res.ok) throw new Error("failed to delete todo");
};
