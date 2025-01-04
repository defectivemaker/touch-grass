"use server";
import { revalidatePath } from "next/cache";
import { createClient } from "@/utils/supabase/server";

export async function login(formData) {
    const supabase = createClient();
    const data = {
        email: formData.get("email"),
        password: formData.get("password"),
    };
    const { data: signInData, error } = await supabase.auth.signInWithPassword(
        data
    );
    if (error) {
        return { error: error.message };
    }
    revalidatePath("/home");
    return { success: true, user: signInData.user };
}

export async function signup(formData) {
    const supabase = createClient();
    const email = formData.get("email");
    const password = formData.get("password");
    const confirmPassword = formData.get("confirmPassword");

    if (password !== confirmPassword) {
        return { error: "Passwords do not match" };
    }

    const { data: signUpData, error } = await supabase.auth.signUp({
        email,
        password,
    });
    if (error) {
        return { error: error.message };
    }
    revalidatePath("/home");
    return { success: true, user: signUpData.user };
}
