"use client";

import React from "react";

import { ThemeProvider } from "next-themes";

import { NextUIProvider } from "@nextui-org/react";

export default function Providers({ children }) {
    return (
        <ThemeProvider attribute="class" defaultTheme="system">
            <NextUIProvider>{children}</NextUIProvider>
        </ThemeProvider>
    );
}
