import { useRouter } from "next/navigation";
export const TouchGrassLogo = () => {
    const router = useRouter();
    const handleClick = (e) => {
        e.preventDefault();
        router.push("/");
    };

    return (
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 100 100"
            className="w-16 h-16 md:w-24 md:h-24"
            onClick={handleClick}
            isClickable
        >
            <path
                d="M50,10 C70,10 85,25 90,45 C95,65 85,85 65,90 C45,95 25,85 20,65 C15,45 30,10 50,10 Z"
                fill="#a3e635"
            />
            <path
                d="M50,20 C65,20 75,30 78,45 C81,60 75,75 60,78 C45,81 30,75 27,60 C24,45 35,20 50,20 Z"
                fill="#ffffff"
                fillOpacity="0.2"
            />
        </svg>
    );
};
