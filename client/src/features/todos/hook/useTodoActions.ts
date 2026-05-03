import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/features/auth/context/AuthContext";
import { completeTodo, deleteTodo } from "../api/todo.api";

export const useCompleteTodo = () => {
  const queryClient = useQueryClient();
  const { token } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (!token) {
        throw new Error("Missing auth token");
      }

      return completeTodo({ token, id });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });
};

export const useDeleteTodo = () => {
  const queryClient = useQueryClient();
  const { token } = useAuth();

  return useMutation({
    mutationFn: (id: string) => {
      if (!token) {
        throw new Error("Missing auth token");
      }

      return deleteTodo({ token, id });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });
};
