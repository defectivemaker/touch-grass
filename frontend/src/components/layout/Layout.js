import React from "react";
import { TouchGrassLogo } from "@/components/Logo";

export default function Layout({ children }) {
    return (
        <div className="bg-[#0a3622] text-white min-h-screen">
            <header className="p-4 flex justify-between items-center">
                <div className="flex items-center space-x-4">
                    <TouchGrassLogo />
                    <h1 className="text-3xl md:text-4xl lg:text-5xl font-bold">
                        Touch grass
                    </h1>
                </div>
            </header>

            <main className="container mx-auto p-4 md:p-8">{children}</main>
        </div>
    );
}
