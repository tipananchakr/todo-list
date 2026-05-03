import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/features/auth/context/AuthContext";
import { createTodo } from "../api/todo.api";

export const useCreateTodo = () => {
  const queryClient = useQueryClient();
  const { token } = useAuth();

  return useMutation({
    mutationFn: (data: { body: string }) => {
      if (!token) {
        throw new Error("Missing auth token");
      }

      return createTodo({ token, body: data.body });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });
};
