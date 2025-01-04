// utils/apiClient.js
import { createClient as createSupabaseClient } from "./supabase/client";

export const createApiClient = (serverSupabase = null) => {
    const fetchWithAuth = async (url, options = {}) => {
        const baseUrl = process.env.NEXT_PUBLIC_GOLANG_URL;
        const protocol = "https://";
        const fullUrl = `${protocol}${baseUrl}${url}`;

        if (!options.headers) {
            options.headers = {};
        }

        let supabase;
        if (typeof window === "undefined" && serverSupabase) {
            // Server-side
            supabase = serverSupabase;
        } else {
            // Client-side
            supabase = createSupabaseClient();
        }

        const {
            data: { session },
        } = await supabase.auth.getSession();

        if (session?.access_token) {
            options.headers["Authorization"] = `Bearer ${session.access_token}`;
        }

        try {
            console.log(
                `Sending ${options.method || "GET"} request to ${fullUrl}`
            );
            const response = await fetch(fullUrl, options);

            if (response.status === 401) {
                console.log("Received 401, attempting to refresh token");
                // Token might be expired, try to refresh
                const { data, error } = await supabase.auth.refreshSession();
                if (error) throw error;
                options.headers[
                    "Authorization"
                ] = `Bearer ${data.session.access_token}`;
                console.log("Token refreshed, retrying request");
                return fetch(fullUrl, options);
            }

            if (!response.ok) {
                console.error(`API call failed with status ${response.status}`);
                const errorText = await response.text();
                console.error(`Error response: ${errorText}`);
                throw new Error(
                    `API call failed: ${response.status} ${errorText}`
                );
            }

            return response.json();
        } catch (error) {
            console.error("API call error:", error);
            throw error;
        }
    };

    return {
        get: (url) => fetchWithAuth(url),
        post: (url, data) =>
            fetchWithAuth(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(data),
            }),
        put: (url, data) =>
            fetchWithAuth(url, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(data),
            }),
        delete: (url) => fetchWithAuth(url, { method: "DELETE" }),
    };
};
