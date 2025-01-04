// api/stats.js
import { createApiClient } from "@/utils/apiClient";

export async function fetchStats(serverSupabase = null) {
    const apiClient = createApiClient(serverSupabase);
    try {
        return await apiClient.get("/statistics");
    } catch (error) {
        console.error("Failed to fetch stats:", error);
        throw new Error("Failed to fetch stats");
    }
}
