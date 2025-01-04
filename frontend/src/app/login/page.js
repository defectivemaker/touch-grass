"use client";
import React, { useState, useEffect, Suspense } from "react";
import Layout from "@/components/layout/Layout";
import { login, signup } from "./actions";
import { Tabs, Tab, Input, Card, CardBody, Button } from "@nextui-org/react";
import { useRouter, useSearchParams } from "next/navigation";
import { createApiClient } from "@/utils/apiClient";

const Login = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [error, setError] = useState(null);
    const [activeTab, setActiveTab] = useState("login");
    const [email, setEmail] = useState("");

    useEffect(() => {
        const tabParam = searchParams.get("tab");
        const emailParam = searchParams.get("email");

        if (tabParam === "signup") {
            setActiveTab("signup");
        }

        if (emailParam) {
            setEmail(emailParam);
        }
    }, [searchParams]);

    const handleSubmit = async (event) => {
        event.preventDefault();
        const formData = new FormData(event.target);

        try {
            const result = await (activeTab === "signup"
                ? signup(formData)
                : login(formData));
            if (result.error) {
                setError(result.error);
                return;
            }

            if (result.success) {
                if (activeTab === "signup") {
                    await addUser(result.user);
                }
                router.push("/home");
            }
        } catch (error) {
            setError(`An error occurred during ${activeTab}`);
            console.error(`${activeTab} error:`, error);
        }
    };

    const addUser = async (user) => {
        const apiClient = createApiClient();
        try {
            await apiClient.post("/add-user", {
                id: user.id,
                email: user.email,
                // Add any other necessary fields
            });
            console.log("User added successfully");
        } catch (error) {
            console.error("Failed to add user:", error);
            setError("Failed to add user to the system");
        }
    };
    return (
        <div className="flex justify-center items-center min-h-[calc(100vh-100px)] dark">
            <Card className="w-full max-w-md bgneutral-900 text-gray-800">
                <CardBody className="p-8">
                    <h2 className="text-2xl font-bold mb-6 text-center text-neutral-300">
                        Welcome to Touch Grass
                    </h2>
                    <Tabs
                        aria-label="Login options"
                        className="mb-4"
                        selectedKey={activeTab}
                        onSelectionChange={setActiveTab}
                    >
                        <Tab key="login" title="Login">
                            <form
                                className="space-y-4 mt-4"
                                onSubmit={handleSubmit}
                            >
                                <Input
                                    id="email"
                                    name="email"
                                    type="email"
                                    label="Email"
                                    placeholder="Enter your email"
                                    required
                                    className="w-full"
                                    variant="flat"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                                <Input
                                    id="password"
                                    name="password"
                                    type="password"
                                    label="Password"
                                    placeholder="Enter your password"
                                    required
                                    className="w-full"
                                />
                                <Button
                                    className="w-full text-white"
                                    variant="bordered"
                                    type="submit"
                                >
                                    Log in
                                </Button>
                            </form>
                        </Tab>
                        <Tab key="signup" title="Sign Up">
                            <form
                                className="space-y-4 mt-4"
                                onSubmit={handleSubmit}
                            >
                                <Input
                                    id="signup-email"
                                    name="email"
                                    type="email"
                                    label="Email"
                                    placeholder="Enter your email"
                                    required
                                    className="w-full"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                                <Input
                                    id="signup-password"
                                    name="password"
                                    type="password"
                                    label="Password"
                                    placeholder="Choose a password"
                                    required
                                    className="w-full"
                                />
                                <Input
                                    id="confirm-password"
                                    name="confirmPassword"
                                    type="password"
                                    label="Confirm Password"
                                    placeholder="Confirm your password"
                                    required
                                    className="w-full"
                                />
                                <Button
                                    className="w-full text-white"
                                    variant="bordered"
                                    type="submit"
                                >
                                    Sign up
                                </Button>
                            </form>
                        </Tab>
                    </Tabs>
                    {error && (
                        <p className="text-red-500 mt-4 text-center">{error}</p>
                    )}
                </CardBody>
            </Card>
        </div>
    );
};

export default function LoginPage() {
    return (
        <Layout>
            <Suspense fallback={<div>Loading..</div>}>
                <Login />
            </Suspense>
        </Layout>
    );
}
