import { Box, Flex, Button, Text } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import { useColorMode, useColorModeValue } from "../../../components/ui/color-mode";
import { useAuth } from "@/features/auth/context/AuthContext";

export default function Navbar() {
	const { colorMode, toggleColorMode } = useColorMode();
	const { user, logout } = useAuth();

	return (
		<div >
			<Box bg={useColorModeValue("gray.400", "gray.700")} px={4} my={4} borderRadius={"5"}>
				<Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
					{/* LEFT SIDE */}
					<Flex
						justifyContent={"center"}
						alignItems={"center"}
						gap={3}
						display={{ base: "none", sm: "flex" }}
					>
						<Text fontSize={"lg"} fontWeight={500}>
							Daily Tasks
						</Text>
					</Flex>

					{/* RIGHT SIDE */}
					<Flex alignItems={"center"} gap={3}>
						{user && (
							<Text fontSize="sm" color="fg.muted" display={{ base: "none", md: "block" }}>
								{user.email}
							</Text>
						)}
						
						{/* Toggle Color Mode */}
						<Button onClick={toggleColorMode}>
							{colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
						</Button>
						{user && (
							<Button variant="outline" onClick={logout}>
								Logout
							</Button>
						)}
					</Flex>
				</Flex>
			</Box>
		</div>
	);
}
