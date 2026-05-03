import { useQuery } from "@tanstack/react-query"
import { useAuth } from "@/features/auth/context/AuthContext"
import { getTodos } from "../api/todo.api"

export const useTodos = () => {
  const { token } = useAuth();

  return useQuery({
    queryKey: ["todos"],
    queryFn: () => getTodos(token ?? ""),
    enabled: Boolean(token),
  })
}
