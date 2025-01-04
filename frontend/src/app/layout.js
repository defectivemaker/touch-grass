import Providers from "@/utils/providers";

import "./globals.css";

export default function RootLayout({ children }) {
    return (
        <html lang="en">
            <title>Touch grass</title>
            <link rel="icon" href="/icon.svg" type="image/svg+xml" />
            <body>
                <Providers>{children}</Providers>
            </body>
        </html>
    );
}
