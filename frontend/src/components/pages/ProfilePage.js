import React, { useState, useEffect } from "react";
import { Button, Input, Card, Badge } from "@nextui-org/react";
import GlassCard from "@/components/ui/GlassCard";
import { getProfile } from "@/pages/api/getProfile";
import { ShoppingCart } from "lucide-react";
import { useRouter } from "next/navigation";

export default function ProfilePage() {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [token, setToken] = useState("");
    const [uuid, setUuid] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [hasDevice, setHasDevice] = useState(false);
    const router = useRouter();

    useEffect(() => {
        fetchProfile();
    }, []);

    const fetchProfile = async () => {
        try {
            setIsLoading(true);
            const profileData = await getProfile();
            setUsername(profileData.username || "");
            setEmail(profileData.email || "");
            setToken(profileData.token || "");
            setUuid(profileData.uuid || "");
            setHasDevice(profileData.hasDevice || false);
            setError(null);
        } catch (err) {
            setError("Failed to fetch profile data");
            console.error("Error fetching profile:", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleUpdateProfile = async () => {
        // TODO: Implement profile update logic
        console.log("Profile update requested", { username });
        // After updating, re-fetch the profile
        await fetchProfile();
    };

    const handleBuyDevice = () => {
        // Replace with your actual Stripe checkout URL
        router.push("/home");
    };

    if (isLoading) {
        return <div>Loading profile...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className="space-y-6">
            <GlassCard>
                <h2 className="text-lg font-bold mb-4">Buy device</h2>
                <Card className="p-4 bg-opacity-20 backdrop-filter backdrop-blur-lg">
                    <div className="flex items-center justify-between">
                        <div>
                            {/* <h3 className="text-xl font-semibold mb-2">
                                Touch Grass Explorer
                            </h3> */}
                            <p className="mb-4">
                                Discover new hotspots and track your adventures!
                            </p>
                        </div>
                        <ShoppingCart
                            size={48}
                            className="text-green-500 ml-4 mr-2 mb-2"
                        />
                    </div>
                    {hasDevice ? (
                        <Badge
                            color="success"
                            content="Owned"
                            placement="bottom-right"
                        >
                            <Button
                                color="success"
                                variant="bordered"
                                isDisabled
                                className="w-full"
                            >
                                Device Purchased
                            </Button>
                        </Badge>
                    ) : (
                        <Button
                            color="success"
                            onClick={handleBuyDevice}
                            className="w-full"
                        >
                            Buy Device
                        </Button>
                    )}
                </Card>
            </GlassCard>

            <GlassCard>
                <h2 className="text-lg font-bold mb-4">Your Profile</h2>
                <div className="space-y-4">
                    <div className="bg-opacity-10">
                        <label
                            htmlFor="username"
                            className="block text-sm font-medium mb-1"
                        >
                            Username
                        </label>
                        <Input
                            id="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Enter your username"
                            className="bg-opacity-10"
                        />
                    </div>
                    <div>
                        <label
                            htmlFor="email"
                            className="block text-sm font-medium mb-1"
                        >
                            Email
                        </label>
                        <Input
                            id="email"
                            type="email"
                            value={email}
                            readOnly
                            className="bg-opacity-10"
                        />
                    </div>
                    <div>
                        <label
                            htmlFor="token"
                            className="block text-sm font-medium mb-1"
                        >
                            Token
                        </label>
                        <Input
                            id="token"
                            value={token}
                            readOnly
                            className="bg-opacity-10"
                        />
                    </div>
                    <div>
                        <label
                            htmlFor="uuid"
                            className="block text-sm font-medium mb-1"
                        >
                            UUID
                        </label>
                        <Input
                            id="uuid"
                            value={uuid}
                            readOnly
                            className="bg-opacity-10"
                        />
                    </div>
                    <Button
                        color="primary"
                        className="bg-opacity-30"
                        onClick={handleUpdateProfile}
                    >
                        Update Profile
                    </Button>
                </div>
            </GlassCard>
        </div>
    );
}
