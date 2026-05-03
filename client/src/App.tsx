
"use client"

import { Container, HStack, Spinner, Stack } from "@chakra-ui/react"
import Navbar from "./features/todos/components/Navbar"
import TodoForm from "./features/todos/components/TodoForm"
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/react-query";
import TodoList from "@/features/todos/components/TodoList";
import { AuthProvider, useAuth } from "@/features/auth/context/AuthContext";
import AuthForm from "@/features/auth/components/AuthForm";
import { Toaster } from "@/components/ui/toaster";

function TodoApp() {
  const { user, isCheckingSession } = useAuth();

  if (isCheckingSession) {
    return (
      <Stack h="100vh" align="center" justify="center">
        <Spinner size="lg" />
      </Stack>
    );
  }

  if (!user) {
    return (
      <Container maxW="900px" h="100vh">
        <AuthForm />
      </Container>
    );
  }

  return (
    <HStack h="100vh" flexDirection={"column"} alignItems="stretch">
      <Container maxW={"900px"} h="100%" display={"flex"} flexDirection="column" gap={4}>
        <Navbar />
        <TodoForm />
        <TodoList />
      </Container>
    </HStack>
  );
}

function App() {

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <TodoApp />
        <Toaster />
      </AuthProvider>
    </QueryClientProvider>
  )
}

export default App
