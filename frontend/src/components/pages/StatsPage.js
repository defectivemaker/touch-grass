"use client";
import { useState, useEffect } from "react";
import { fetchStats } from "@/pages/api/stats";
import GlassCard from "@/components/ui/GlassCard";

export default function StatsPage() {
    const [stats, setStats] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function loadStats() {
            try {
                const data = await fetchStats();
                setStats(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        loadStats();
    }, []);

    if (loading) return <div>Loading stats...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <GlassCard>
            <h2 className="text-lg font-bold mb-4">Your Adventure Stats</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div>
                    <h3 className="font-bold">Level</h3>
                    <p className="text-3xl">3{stats.level}</p>
                </div>
                <div>
                    <h3 className="font-bold">Hotspots Found</h3>
                    <p className="text-3xl">21{stats.hotspotsFound}</p>
                </div>
                <div>
                    <h3 className="font-bold">Distance Travelled</h3>
                    <p className="text-3xl">{stats.distanceTravelled}40 km</p>
                </div>
                <div>
                    <h3 className="font-bold">Next Level</h3>
                    <div className="w-full bg-success-200 dark:bg-success-800 rounded-full h-2.5">
                        <div
                            className="bg-success-600 h-2.5 rounded-full"
                            style={{ width: `${stats.nextLevelProgress}%` }}
                        ></div>
                    </div>
                </div>
            </div>
        </GlassCard>
    );
}
