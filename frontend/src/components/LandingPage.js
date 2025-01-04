"use client";
import React, { useEffect, useRef, useState } from "react";
import { Button, Card, Input } from "@nextui-org/react";
import { MapPin, Trophy, ArrowRight } from "lucide-react";

import { motion, AnimatePresence } from "framer-motion";
import { TouchGrassLogo } from "@/components/Logo";
import InteractiveAustraliaMap from "./InteractiveMap";
import GlassCard from "./ui/GlassCard";
import { ParallaxProvider, Parallax } from "react-scroll-parallax";
import { useRouter } from "next/navigation";

const defaultPhoto = {
    thumbnail: "/aus_travel_photos/030038550019_19A.jpg",
    fullSize: "/aus_travel_photos/030038550019_19A.jpg",
    alt: `Photo 1`,
    latitude: 0,
    longitude: 0,
};

const LandingPage = ({ photos }) => {
    const router = useRouter();
    const [focusedPhoto, setFocusedPhoto] = useState(defaultPhoto);
    const [isUserInteracting, setIsUserInteracting] = useState(false);
    const [email, setEmail] = useState("");
    const timeoutRef = useRef(null);

    const redirectToSignup = (e) => {
        e.preventDefault();
        router.push(`/login?tab=signup&email=${encodeURIComponent(email)}`);
    };

    useEffect(() => {
        const cyclePhotos = () => {
            if (!isUserInteracting) {
                const randomIndex = Math.floor(Math.random() * photos.length);
                setFocusedPhoto(photos[randomIndex]);
            }
        };
        timeoutRef.current = setInterval(cyclePhotos, 3000);
        return () => clearInterval(timeoutRef.current);
    }, [photos, isUserInteracting]);

    const handlePhotoFocus = (photo) => {
        setIsUserInteracting(true);
        setFocusedPhoto(photo);
        clearInterval(timeoutRef.current);
    };

    const handleMapClick = () => {
        setIsUserInteracting(false);
    };
    return (
        <ParallaxProvider>
            <div className="bg-[#0a3622] text-white min-h-screen relative overflow-hidden">
                <Parallax speed={-20}>
                    <div className="absolute inset-0 z-0 hidden md:block">
                        <InteractiveAustraliaMap
                            photos={photos}
                            focusedPhoto={focusedPhoto}
                            setFocusedPhoto={setFocusedPhoto}
                            isUserInteracting={isUserInteracting}
                            setIsUserInteracting={setIsUserInteracting}
                            timeoutRef={timeoutRef}
                            handlePhotoFocus={handlePhotoFocus}
                            handleMapClick={handleMapClick}
                        />
                    </div>
                </Parallax>

                <div className="relative z-10">
                    <header className="p-4 md:p-8 flex items-center max-w-7xl mx-auto">
                        <TouchGrassLogo />
                        <h1 className="text-2xl md:text-3xl font-bold ml-4">
                            Touch Grass
                        </h1>
                    </header>

                    <main className="max-w-7xl mx-auto p-4 md:p-8 space-y-16 md:space-y-32 mt-8">
                        <Parallax speed={10}>
                            <section className="flex flex-col md:flex-row md:items-center w-full">
                                <div className="space-y-4 md:space-y-6 md:w-1/2">
                                    <h2 className="text-3xl md:text-5xl font-bold leading-tight">
                                        Discover Australia, One Hotspot at a
                                        Time
                                    </h2>
                                    <p className="text-lg md:text-2xl">
                                        Embark on a nationwide adventure,
                                        connect to Wi-Fi hotspots, and compete
                                        with others!
                                    </p>
                                    <Button
                                        size="lg"
                                        endContent={<ArrowRight />}
                                        onClick={redirectToSignup}
                                        className="text-lg px-6 py-3 md:px-8 md:py-6 bg-[#a3e635] text-black mb-8"
                                    >
                                        Start Exploring
                                    </Button>
                                </div>
                                <div className="mt-8 md:mt-0 md:w-1/2 relative">
                                    <div className="relative flex justify-end ">
                                        <div className="relative w-full h-48 md:w-96 md:h-96 rounded-lg overflow-hidden shadow-lg mb-16 md:mb-0 md:mr-10 mt-6 md:mt-0">
                                            <AnimatePresence mode="sync">
                                                {focusedPhoto && (
                                                    <motion.div
                                                        key={
                                                            focusedPhoto.fullSize
                                                        }
                                                        initial={{
                                                            opacity: 0,
                                                            x: 100,
                                                            scale: 0.9,
                                                        }}
                                                        animate={{
                                                            opacity: 1,
                                                            x: 0,
                                                            scale: 1,
                                                            transition: {
                                                                type: "spring",
                                                                stiffness: 300,
                                                                damping: 30,
                                                                mass: 1,
                                                                duration: 0.6,
                                                            },
                                                        }}
                                                        exit={{
                                                            opacity: 0,
                                                            x: -100,
                                                            scale: 0.9,
                                                            transition: {
                                                                type: "spring",
                                                                stiffness: 300,
                                                                damping: 30,
                                                                mass: 1,
                                                                duration: 0.6,
                                                            },
                                                        }}
                                                        className="absolute top-0 left-0 w-full h-full"
                                                    >
                                                        <motion.img
                                                            src={
                                                                focusedPhoto.fullSize
                                                            }
                                                            alt={
                                                                focusedPhoto.alt
                                                            }
                                                            className="w-full h-full object-cover"
                                                            initial={{
                                                                scale: 1.1,
                                                            }}
                                                            animate={{
                                                                scale: 1,
                                                            }}
                                                            exit={{
                                                                scale: 1.1,
                                                            }}
                                                            transition={{
                                                                duration: 0.6,
                                                            }}
                                                        />
                                                    </motion.div>
                                                )}
                                            </AnimatePresence>
                                        </div>
                                        <div className="md:hidden absolute -bottom-40 left-0 right-0 h-64 z-10 w-full">
                                            <InteractiveAustraliaMap
                                                photos={photos}
                                                focusedPhoto={focusedPhoto}
                                                setFocusedPhoto={
                                                    setFocusedPhoto
                                                }
                                                isUserInteracting={
                                                    isUserInteracting
                                                }
                                                setIsUserInteracting={
                                                    setIsUserInteracting
                                                }
                                                timeoutRef={timeoutRef}
                                                handlePhotoFocus={
                                                    handlePhotoFocus
                                                }
                                                handleMapClick={handleMapClick}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </section>
                        </Parallax>
                        <Parallax speed={5}>
                            <section className="space-y-12 relative top-10">
                                <h3 className="text-3xl lg:text-4xl font-bold text-center">
                                    How It Works
                                </h3>
                                <div className="grid md:grid-cols-3 gap-8 lg:gap-12">
                                    <GlassCard
                                        className="shadow-2xl p-6"
                                        glowColor="rgba(52, 211, 153, 0.5)"
                                    >
                                        <MapPin className="w-16 h-16 mb-4 text-emerald-400" />
                                        <h4 className="text-2xl font-bold mb-4">
                                            Discover Hotspots
                                        </h4>
                                        <p className="text-lg">
                                            Find unique Wi-Fi locations across
                                            Australia
                                        </p>
                                    </GlassCard>
                                    <GlassCard
                                        className="shadow-2xl p-6"
                                        glowColor="rgba(251, 191, 36, 0.5)"
                                    >
                                        <Trophy className="w-16 h-16 mb-4 text-yellow-400" />
                                        <h4 className="text-2xl font-bold mb-4">
                                            Earn Points
                                        </h4>
                                        <p className="text-lg">
                                            Accumulate points for each new
                                            connection
                                        </p>
                                    </GlassCard>
                                    <GlassCard
                                        className="shadow-2xl p-6"
                                        glowColor="rgba(147, 197, 253, 0.5)"
                                    >
                                        <MapPin className="w-16 h-16 mb-4 text-blue-400" />
                                        <h4 className="text-2xl font-bold mb-4">
                                            Compete
                                        </h4>
                                        <p className="text-lg">
                                            Rise through the ranks on our global
                                            leaderboard
                                        </p>
                                    </GlassCard>
                                </div>
                            </section>
                        </Parallax>

                        <Parallax speed={15}>
                            <section className="space-y-8">
                                <h3 className="text-3xl lg:text-4xl font-bold text-center">
                                    {" "}
                                    Join the Movement
                                </h3>

                                <div className="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4">
                                    <Input
                                        placeholder="Enter email"
                                        type="email"
                                        className="flex-grow "
                                        classNames={{
                                            // base: "max-w-full",
                                            input: [
                                                "bg-transparent",
                                                "text-white",
                                                "placeholder:text-white/60",
                                            ],
                                            inputWrapper: [
                                                "bg-white/20",
                                                "backdrop-blur-md",
                                                "border-none",
                                                "group-data-[focused=true]:bg-white/20",
                                                "!ring-0",
                                                "!ring-offset-0",
                                                "focus:!bg-transparent",
                                            ],
                                        }}
                                        onChange={(e) =>
                                            setEmail(e.target.value)
                                        }
                                    />
                                    <Button
                                        color="success"
                                        className="bg-[#a3e635]"
                                        onClick={(e) => redirectToSignup(e)}
                                    >
                                        Sign Up
                                    </Button>
                                </div>
                            </section>
                        </Parallax>
                    </main>

                    <footer className="mt-24 p-8 text-center text-lg">
                        <p>&copy; 2024 Touch Grass. All rights reserved.</p>
                    </footer>
                    {/* </div> */}
                </div>
            </div>
        </ParallaxProvider>
    );
};

export default LandingPage;
