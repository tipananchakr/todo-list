import { Button, Field, Input, Stack } from "@chakra-ui/react";
import { useForm } from "react-hook-form";
import { useCreateTodo } from "../hook/useCreateTodo";
import { toaster } from "@/components/ui/toaster";

interface FormData {
  body: string;
}

const TodoForm = () => {
  const mutation = useCreateTodo()

  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
    reset
  } = useForm<FormData>()

  const onSubmit = handleSubmit((data) => {
    mutation.mutate(data, {
      onSuccess: () => reset(),
      onError: (error) => {
        toaster.create({
          title: "Cannot create todo",
          description: error instanceof Error ? error.message : "Please try again",
          type: "error",
        });
      },
    });
  })

  return (
    <form onSubmit={onSubmit}>
      <Stack w={"100%"} gap="4" flex={1} direction={"row"}>
        <Field.Root invalid={!!errors.body} w={"100%"}>
          <Input 
            type="text" 
            placeholder="Add a task"
            {...register("body", {
              required: "This field is required",
              minLength: {
                value: 3,
                message: "Must be at least 3 characters"
              }
            })} 
          />
          <Field.ErrorText>{errors.body?.message}</Field.ErrorText>
        </Field.Root>


        <Button 
          type="submit" 
          disabled={!isValid || mutation.isPending}
          loading={mutation.isPending}
        >
          Add
        </Button>

      </Stack>
    </form>
  )
}

export default TodoForm
