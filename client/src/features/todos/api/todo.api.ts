const BASE_URL = import.meta.env.NEXT_PUBLIC_API_URL || "http://localhost:3000";

export const getTodos = async () => {
  const res = await fetch(`${BASE_URL}/api/todos`);
  if (!res.ok) {
    throw new Error("Failed to fetch todos");
  }
  return res.json();
}

export const createTodo = async (data: { body: string }) => {
  const res = await fetch(`${BASE_URL}/api/todos`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!res.ok) throw new Error("failed to create todo");

  return res.json();
};
