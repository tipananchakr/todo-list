
"use client"

import { Container, HStack } from "@chakra-ui/react"
import Navbar from "./features/todos/components/Navbar"
import TodoForm from "./features/todos/components/TodoForm"
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/react-query";
import TodoList from "@/features/todos/components/TodoList";

function App() {

  return (
    <QueryClientProvider client={queryClient}>
      <HStack h="100vh" flexDirection={"column"} alignItems="stretch"  >
        <Container maxW={"900px"} h="100%" display={"flex"} flexDirection="column" gap={4} >
          <Navbar />
          <TodoForm />
          <TodoList />
        </Container>
      </HStack>
    </QueryClientProvider>
  )
}

export default App
