"use client";
import { useState } from "react";
import { Tabs, Tab } from "@nextui-org/react";
import Layout from "@/components/layout/Layout";
import dynamic from "next/dynamic";
// import DiscoverPage from "@/components/pages/DiscoverPage";
// import StatsPage from "@/components/pages/StatsPage";
// import ProfilePage from "@/components/pages/ProfilePage";
import { MapPin, Trophy, User } from "lucide-react";

const DiscoverPage = dynamic(() => import("@/components/pages/DiscoverPage"), {
    ssr: false,
});
const StatsPage = dynamic(() => import("@/components/pages/StatsPage"), {
    ssr: false,
});
const ProfilePage = dynamic(() => import("@/components/pages/ProfilePage"), {
    ssr: false,
});

export default function Home() {
    const [activeTab, setActiveTab] = useState("discover");

    return (
        <Layout>
            <Tabs
                selectedKey={activeTab}
                onSelectionChange={setActiveTab}
                aria-label="Touch Grass Navigation"
                className="mb-4 flex flex-grow"
                classNames={{
                    tabList:
                        "bg-opacity-10 backdrop-filter backdrop-blur-lg w-full md:w-auto",
                }}
                color="success"
            >
                <Tab
                    key="discover"
                    title={
                        <div className="flex items-center space-x-2">
                            <MapPin />
                            <span>Discover</span>
                        </div>
                    }
                >
                    <DiscoverPage />
                </Tab>
                <Tab
                    key="stats"
                    title={
                        <div className="flex items-center space-x-2">
                            <Trophy />
                            <span>Stats</span>
                        </div>
                    }
                >
                    <StatsPage />
                </Tab>
                <Tab
                    key="profile"
                    title={
                        <div className="flex items-center space-x-2">
                            <User />
                            <span>Profile</span>
                        </div>
                    }
                >
                    <ProfilePage />
                </Tab>
            </Tabs>
        </Layout>
    );
}
