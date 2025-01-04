// api/stats.js
import { createApiClient } from "@/utils/apiClient";

export async function getProfile(serverSupabase = null) {
    const apiClient = createApiClient(serverSupabase);
    try {
        return await apiClient.get("/get-profile");
    } catch (error) {
        console.error("Failed to fetch token:", error);
        throw new Error("Failed to fetch token");
    }
}
