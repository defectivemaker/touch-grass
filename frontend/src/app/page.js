import LandingPage from "@/components/LandingPage";

// const apiUrl = process.env.NEXT_PUBLIC_API_URL;
// const protocol = "https://";
// const apiUrl = `${protocol}${process.env.NEXT_PUBLIC_NEXTJS_URL}`;
const apiUrl = process.env.NEXT_PUBLIC_NEXTJS_URL;

async function getPhotos() {
    const res = await fetch(`${apiUrl}/api/photos`, {
        cache: "no-store",
    });
    if (!res.ok) {
        throw new Error("Failed to fetch photos");
    }
    return res.json();
}

export default async function Home() {
    const photos = await getPhotos();
    return <LandingPage photos={photos} />;
}
