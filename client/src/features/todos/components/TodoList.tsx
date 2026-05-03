import { Badge, Button, HStack, Stack, Text } from "@chakra-ui/react";
import { useTodos } from "../hook/useTodos";
import { useCompleteTodo, useDeleteTodo } from "../hook/useTodoActions";

const TodoList = () => {
  const { data, isLoading, error } = useTodos();
  const completeMutation = useCompleteTodo();
  const deleteMutation = useDeleteTodo();

  if (isLoading) return <Text color="fg.muted">Loading...</Text>;
  if (error) return <Text color="red.500">Cannot load todos</Text>;
  if (!data?.length) return <Text color="fg.muted" textAlign={"center"}>No todos yet.</Text>;

  return (
    <Stack as="ul" gap={3} p={0} listStyleType="none">
      {data.map((todo) => (
        <HStack
          as="li"
          key={todo.id}
          borderWidth="1px"
          borderRadius="md"
          px={4}
          py={3}
          justify="space-between"
          gap={3}
        >
          <HStack gap={3} minW={0}>
            <Badge colorPalette={todo.completed ? "green" : "gray"}>
              {todo.completed ? "Done" : "Open"}
            </Badge>
            <Text
              textDecoration={todo.completed ? "line-through" : "none"}
              color={todo.completed ? "fg.muted" : "fg"}
              truncate
            >
              {todo.body}
            </Text>
          </HStack>

          <HStack gap={2}>
            {!todo.completed && (
              <Button
                size="sm"
                variant="outline"
                onClick={() => completeMutation.mutate(todo.id)}
                loading={completeMutation.isPending}
              >
                Done
              </Button>
            )}
            <Button
              size="sm"
              colorPalette="red"
              variant="outline"
              onClick={() => deleteMutation.mutate(todo.id)}
              loading={deleteMutation.isPending}
            >
              Delete
            </Button>
          </HStack>
        </HStack>
      ))}
    </Stack>
  );
};

export default TodoList;
