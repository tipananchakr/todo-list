import { useTodos } from "../hook/useTodos";

const TodoList = () => {
  const { data, isLoading, error } = useTodos();

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Error</p>;

  return (
    <ul>
      {data.map((todo: any) => (
        <li key={todo.id}>{todo.body}</li>
      ))}
    </ul>
  );
};

export default TodoList;