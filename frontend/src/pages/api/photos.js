import { GetRandomCoordInAustralia } from "@/utils/displayMap";
import { promises as fs } from "fs";
import path from "path";

const photo_path = "public/aus_travel_photos/";
const feed_path = "/aus_travel_photos/";
export default async function handler(req, res) {
    const photoDirectory = path.join(process.cwd(), photo_path);
    const photoFilenames = await fs.readdir(photoDirectory);
    const photos = photoFilenames.map((image, index) => {
        const { latitude, longitude } = GetRandomCoordInAustralia();
        return {
            thumbnail: feed_path + image,
            fullSize: feed_path + image,
            alt: `Photo ${index + 1}`,
            latitude,
            longitude,
        };
    });
    res.status(200).json(photos);
}
