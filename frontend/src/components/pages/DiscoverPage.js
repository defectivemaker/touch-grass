"use client";

import React, { useState, useEffect, useRef } from "react";
import { Tabs, Tab, Button, Card } from "@nextui-org/react";
import { Map as MapIcon, List, Pin } from "lucide-react";
import GlassCard from "@/components/ui/GlassCard";
import { TileLayer, Marker, Popup } from "react-leaflet";
import { icon } from "leaflet";
import "leaflet/dist/leaflet.css";
import { createApiClient } from "@/utils/apiClient";
import { createBrowserClient } from "@supabase/ssr";
import { useRouter } from "next/navigation";

import dynamic from "next/dynamic";

const MapContainer = dynamic(
    () => import("react-leaflet").then((mod) => mod.MapContainer),
    { ssr: false }
);

const apiUrl = process.env.NEXT_PUBLIC_GOLANG_URL;

export default function DiscoverPage() {
    const [entries, setEntries] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [mostRecentEntry, setMostRecentEntry] = useState(null);
    const [isPinned, setIsPinned] = useState(false);
    const [notificationSound, setNotificationSound] = useState(null);

    const router = useRouter();

    const supabase = createBrowserClient(
        process.env.NEXT_PUBLIC_SUPABASE_URL,
        process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY
    );

    const mapIcon = icon({
        iconUrl: "/map-pin-2.png",
        iconSize: [32, 32],
    });

    useEffect(() => {
        // Ensure this runs only on the client side
        if (typeof window !== "undefined") {
            setNotificationSound(new Audio("/sounds/notification.wav"));
        }
    }, []);

    const playNotificationSound = () => {
        if (notificationSound) {
            notificationSound
                .play()
                .catch((error) => console.error("Error playing sound:", error));
        }
    };

    useEffect(() => {
        const checkSession = async () => {
            const {
                data: { session },
            } = await supabase.auth.getSession();
            if (session) {
                fetchEntries();
                fetchMostRecentEntry();
                setupWebSocket(session.access_token);
            } else {
                router.push("/login"); // Redirect to login if no session
            }
        };
        checkSession();
    }, []);

    useEffect(() => {
        let intervalId;
        if (mostRecentEntry) {
            intervalId = setInterval(() => {
                checkTime(mostRecentEntry);
            }, 10000);
        }
        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, [mostRecentEntry]);

    const fetchMostRecentEntry = async () => {
        const apiClient = createApiClient();
        setIsLoading(true);
        console.log("SUCCESS");
        const response = await apiClient.get("/get-recent-entry");
        console.log("RES", response);
        console.log("RES", response.status);
        if (response && response.ID != 0) {
            setMostRecentEntry(response);
        }

        setError(null);
        setIsLoading(false);
    };

    const fetchEntries = async () => {
        const apiClient = createApiClient();
        try {
            setIsLoading(true);
            const response = await apiClient.get("/get-entries");
            if (response && Array.isArray(response)) {
                setEntries(response);
            } else {
                setEntries([]);
            }
            console.log("ENTRIES", entries);
            setError(null);
        } catch (err) {
            setError("Failed to fetch entries");
            console.error("Error fetching entries:", err);
            setEntries([]);
        } finally {
            setIsLoading(false);
        }
    };

    const setupWebSocket = (token) => {
        const socket = new WebSocket(`wss://${apiUrl}/ws?jwt=${token}`);
        socket.onmessage = (event) => {
            console.log("SOCKET: Message from server ", event.data);
            fetchMostRecentEntry();

            playNotificationSound();
        };

        socket.onclose = (event) => {
            console.log("SOCKET: Socket Closed Connection: ", event);
        };

        socket.onerror = (error) => {
            console.log("SOCKET: Socket Error: ", error);
        };

        return () => {
            socket.close();
        };
    };

    const pinLocation = async () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    const { latitude, longitude } = position.coords;
                    console.log(
                        `Latitude: ${latitude}, Longitude: ${longitude}`
                    );

                    let latString = String(latitude);
                    let lonString = String(longitude);
                    const apiClient = createApiClient();

                    try {
                        if (!isPinned) {
                            await apiClient.post("/add-location", {
                                entryId: mostRecentEntry.ID,
                                latitude: latString,
                                longitude: lonString,
                            });
                            setIsPinned(true);
                        } else {
                            await apiClient.post("/add-location", {
                                entryId: mostRecentEntry.ID,
                                latitude: "0",
                                longitude: "0",
                            });
                            setIsPinned(false);
                        }
                        fetchEntries(); // Refresh entries after updating
                    } catch (error) {
                        console.error("Error updating location:", error);
                    }
                },
                (error) => {
                    console.error(
                        "Error Code = " + error.code + " - " + error.message
                    );
                }
            );
        } else {
            console.log("Geolocation is not supported by this browser.");
        }
    };

    const calculateMinsAgo = (timestamp) => {
        const now = Math.round(Date.now() / 1000);
        const diff = Math.round((now - timestamp) / 60);
        return diff > 0 ? `${diff} mins ago` : "just now";
    };

    const checkTime = (entry) => {
        if (entry) {
            const now = Math.round(Date.now() / 1000);
            const diff = Math.round((now - entry.RecordedTime) / 60);
            if (diff >= 10) {
                setMostRecentEntry(null);
            }
        }
    };

    if (isLoading) {
        return <div>Loading entries...</div>;
    }

    return (
        <GlassCard>
            {/* <audio ref={audioRef} src="/sounds/notification.wav" /> */}
            <div className="grid md:grid-cols-2 gap-4">
                <div className="grid gap-y-4">
                    <h2 className="text-lg font-bold mb-2">
                        Most Recent Hotspot
                    </h2>
                    {mostRecentEntry ? (
                        <Card className="p-4 bg-opacity-20 backdrop-filter backdrop-blur-lg">
                            <h3 className="font-bold text-lg mb-2">
                                New Hotspot Found!
                            </h3>
                            <p>
                                <strong>Payphone ID:</strong>{" "}
                                {mostRecentEntry.PayphoneID.slice(0, 10)}...
                            </p>
                            <p>
                                <strong>Time:</strong>{" "}
                                {calculateMinsAgo(mostRecentEntry.RecordedTime)}
                            </p>
                            <Button
                                color={isPinned ? "warning" : "success"}
                                className="mt-4 text-white"
                                onClick={pinLocation}
                                startContent={<Pin />}
                            >
                                {isPinned ? "Unpin Location" : "Confirm"}
                            </Button>
                        </Card>
                    ) : (
                        <Card className="p-4 bg-opacity-20 backdrop-filter backdrop-blur-lg">
                            <p>No recent hotspots found.</p>
                        </Card>
                    )}
                </div>
                <div>
                    <h2 className="text-lg font-bold mb-2">Hotspot History</h2>
                    {error ? (
                        <div className="text-red-500">{error}</div>
                    ) : (
                        <Tabs
                            aria-label="History views"
                            classNames={{
                                tabList:
                                    "bg-opacity-10 backdrop-filter backdrop-blur-lg rounded-lg",
                            }}
                        >
                            <Tab
                                key="map"
                                title={
                                    <div className="flex items-center">
                                        <MapIcon className="mr-2" /> Map
                                    </div>
                                }
                            >
                                <div className="h-96 mt-4">
                                    <MapContainer
                                        center={[-33.07, 151]}
                                        zoom={5}
                                        style={{
                                            height: "100%",
                                            width: "100%",
                                        }}
                                    >
                                        <TileLayer
                                            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                                            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                                        />
                                        {entries.map(
                                            (entry, index) =>
                                                entry &&
                                                entry.MapLatitude !== 0 &&
                                                entry.MapLongitude !== 0 && (
                                                    <Marker
                                                        key={index}
                                                        position={[
                                                            entry.MapLatitude,
                                                            entry.MapLongitude,
                                                        ]}
                                                        icon={mapIcon}
                                                    >
                                                        <Popup>
                                                            <div>
                                                                <h3>
                                                                    Payphone ID:{" "}
                                                                    {entry.PayphoneID.slice(
                                                                        0,
                                                                        10
                                                                    )}
                                                                    ...
                                                                </h3>
                                                                {entry.MapLocation &&
                                                                    entry.MapLocation !==
                                                                        "Unknown" && (
                                                                        <p>
                                                                            Location:{" "}
                                                                            {
                                                                                entry.MapLatitude
                                                                            }
                                                                            {
                                                                                ", "
                                                                            }
                                                                            {
                                                                                entry.MapLongitude
                                                                            }
                                                                        </p>
                                                                    )}
                                                                <p>
                                                                    Time:{" "}
                                                                    {new Date(
                                                                        entry.RecordedTime *
                                                                            1000
                                                                    ).toLocaleString()}
                                                                </p>
                                                            </div>
                                                        </Popup>
                                                    </Marker>
                                                )
                                        )}
                                    </MapContainer>
                                </div>
                            </Tab>
                            <Tab
                                key="list"
                                title={
                                    <div className="flex items-center">
                                        <List className="mr-2" /> List
                                    </div>
                                }
                            >
                                <div className="h-64 overflow-scroll">
                                    {entries.length === 0 ? (
                                        <p className="mt-4">
                                            No entries available.
                                        </p>
                                    ) : (
                                        <div className="space-y-4 mt-4 max-h-96">
                                            {entries.map((entry, index) => (
                                                <Card
                                                    key={index}
                                                    className="p-4 bg-opacity-20 backdrop-filter backdrop-blur-lg"
                                                >
                                                    <p>
                                                        <strong>
                                                            Payphone ID:
                                                        </strong>{" "}
                                                        {entry.PayphoneID.slice(
                                                            0,
                                                            10
                                                        )}
                                                        ...
                                                    </p>
                                                    {entry.MapLocation &&
                                                        entry.MapLocation !==
                                                            "Unknown" && (
                                                            <p>
                                                                Location:{" "}
                                                                {
                                                                    entry.MapLatitude
                                                                }
                                                                {", "}
                                                                {
                                                                    entry.MapLongitude
                                                                }
                                                            </p>
                                                        )}
                                                    <p>
                                                        <strong>Time:</strong>{" "}
                                                        {new Date(
                                                            entry.RecordedTime *
                                                                1000
                                                        ).toLocaleString()}
                                                    </p>
                                                </Card>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </Tab>
                        </Tabs>
                    )}
                </div>
            </div>
        </GlassCard>
    );
}
