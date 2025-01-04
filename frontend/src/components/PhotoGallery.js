import React, { useState } from "react";
import { Modal, Image } from "@nextui-org/react";

const PhotoGallery = ({ photos }) => {
    const [selectedPhoto, setSelectedPhoto] = useState(null);

    return (
        <section className="my-12">
            <h3 className="text-3xl lg:text-4xl font-bold text-center mb-8">
                Explore Australia
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {photos.map((photo, index) => (
                    <div
                        key={index}
                        className="aspect-square overflow-hidden rounded-lg cursor-pointer transition-transform hover:scale-105"
                        onClick={() => setSelectedPhoto(photo)}
                    >
                        <img
                            src={photo.thumbnail}
                            alt={photo.alt}
                            className="w-full h-full object-cover"
                        />
                    </div>
                ))}
            </div>
            <Modal
                isOpen={!!selectedPhoto}
                onClose={() => setSelectedPhoto(null)}
                size="xl"
            >
                {selectedPhoto && (
                    <Image
                        src={selectedPhoto.fullSize}
                        alt={selectedPhoto.alt}
                        className="max-w-full max-h-[80vh] object-contain"
                    />
                )}
            </Modal>
        </section>
    );
};

export default PhotoGallery;
