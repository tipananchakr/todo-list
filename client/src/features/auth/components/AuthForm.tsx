import { Button, Field, Heading, Input, Stack, Text } from "@chakra-ui/react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toaster } from "@/components/ui/toaster";
import { useAuth } from "../context/AuthContext";

interface FormData {
  email: string;
  password: string;
}

const AuthForm = () => {
  const { login, register } = useAuth();
  const [mode, setMode] = useState<"login" | "register">("login");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register: registerField,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<FormData>({ mode: "onChange" });

  const onSubmit = handleSubmit(async (data) => {
    setIsSubmitting(true);
    try {
      if (mode === "login") {
        await login(data);
      } else {
        await register(data);
      }
    } catch (error) {
      toaster.create({
        title: mode === "login" ? "Login failed" : "Registration failed",
        description: error instanceof Error ? error.message : "Please try again",
        type: "error",
      });
    } finally {
      setIsSubmitting(false);
    }
  });

  return (
    <Stack maxW="420px" mx="auto" mt={20} gap={5} w="100%">
      <Stack gap={1}>
        <Heading size="2xl">Daily Tasks</Heading>
        <Text color="fg.muted">Sign in to manage your todo list.</Text>
      </Stack>

      <form onSubmit={onSubmit}>
        <Stack gap={4}>
          <Field.Root invalid={!!errors.email}>
            <Field.Label>Email</Field.Label>
            <Input
              type="email"
              autoComplete="email"
              {...registerField("email", {
                required: "Email is required",
              })}
            />
            <Field.ErrorText>{errors.email?.message}</Field.ErrorText>
          </Field.Root>

          <Field.Root invalid={!!errors.password}>
            <Field.Label>Password</Field.Label>
            <Input
              type="password"
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              {...registerField("password", {
                required: "Password is required",
                minLength: {
                  value: 6,
                  message: "Must be at least 6 characters",
                },
              })}
            />
            <Field.ErrorText>{errors.password?.message}</Field.ErrorText>
          </Field.Root>

          <Button type="submit" disabled={!isValid || isSubmitting} loading={isSubmitting}>
            {mode === "login" ? "Login" : "Create account"}
          </Button>
        </Stack>
      </form>

      <Button
        variant="ghost"
        onClick={() => setMode(mode === "login" ? "register" : "login")}
      >
        {mode === "login" ? "Need an account?" : "Already have an account?"}
      </Button>
    </Stack>
  );
};

export default AuthForm;
