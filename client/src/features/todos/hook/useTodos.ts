import { useQuery } from "@tanstack/react-query"
import { getTodos } from "../api/todo.api"

export const useTodos = () => {
  return useQuery({
    queryKey: ["todos"],
    queryFn: getTodos
  })
}