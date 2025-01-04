import React, { useState, useEffect, useRef } from "react";
import {
    ComposableMap,
    Geographies,
    Geography,
    Marker,
} from "react-simple-maps";
import { motion, AnimatePresence } from "framer-motion";
import { GetRandomCoordInAustralia } from "@/utils/displayMap";

const geoUrl =
    "https://raw.githubusercontent.com/tonywr71/GeoJson-Data/master/australian-states.json";

const InteractiveAustraliaMap = ({
    photos,
    focusedPhoto,
    setFocusedPhoto,
    isUserInteracting,
    setIsUserInteracting,
    timeoutRef,
    handlePhotoFocus,
    handleMapClick,
}) => {
    return (
        <div className="relative w-full " onClick={handleMapClick}>
            <div className="absolute inset-0 z-0">
                <ComposableMap
                    projection="geoMercator"
                    projectionConfig={{
                        scale: 700,
                        center: [134, -28],
                    }}
                    // className="absolute right-5 down-5"
                >
                    <Geographies geography={geoUrl}>
                        {({ geographies }) =>
                            geographies.map((geo) => (
                                <Geography
                                    key={geo.rsmKey}
                                    geography={geo}
                                    fill="#2d6a4f"
                                    stroke="#1b4332"
                                    strokeWidth={0.5}
                                    opacity={0.1}
                                    className="w-full"
                                />
                            ))
                        }
                    </Geographies>
                    {photos &&
                        photos.map((photo, index) => (
                            <Marker
                                key={index}
                                coordinates={[photo.longitude, photo.latitude]}
                            >
                                <motion.circle
                                    r={10}
                                    fill={
                                        focusedPhoto === photo
                                            ? "#fbbf24"
                                            : "#a3e635"
                                    }
                                    onMouseEnter={() => handlePhotoFocus(photo)}
                                    onTouchStart={() => handlePhotoFocus(photo)}
                                    whileHover={{ scale: 1.5 }}
                                    whileTap={{ scale: 1.5 }}
                                    opacity={0.05}
                                />
                            </Marker>
                        ))}
                </ComposableMap>
            </div>
        </div>
    );
};

export default InteractiveAustraliaMap;
