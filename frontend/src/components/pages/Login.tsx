import { getAuthUrl } from '@/client'; // Your HeyAPI client

import { useState } from 'react';

export default function Login() {
    const [loading, setLoading] = useState(false);

    const handleLogin = async () => {
        setLoading(true);
        try {
            const response = await getAuthUrl();

            console.log(response.data)
            
            if (response.data && response.data.url) {
                window.location.href = response.data.url;
            }
        } catch (error) {
            console.error("Failed to get auth URL:", error);
            setLoading(false);
        }
    };

return (
        <div className="h-screen w-screen bg-gradient-to-b from-[#181818] to-black flex flex-col items-center justify-center text-white p-6">
            <div className="max-w-md w-full text-center space-y-10">
                <div className="flex justify-center">
                    <div className="w-24 h-24 bg-[#1DB954] rounded-full flex items-center justify-center shadow-2xl">
                        <svg className="w-14 h-14 text-black" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S16.6 0 12 0zm5.5 17.3c-.2.3-.5.4-.8.2-2.7-1.7-6-2.1-9.9-1.2-.3.1-.6-.1-.7-.4-.1-.3.1-.6.4-.7 4.3-1 8-.5 11.1 1.3.3.1.4.5.2.8zm1.1-2.9c-.3.4-.8.5-1.2.3-3.2-1.9-8-2.5-11.8-1.3-.4.1-1-.1-1.1-.6-.1-.4.1-1 .6-1.1 4.3-1.3 9.7-.6 13.3 1.6.4.2.5.7.2 1.1zm.1-3C15.9 8.6 9.3 8.3 5.4 9.5c-.6.2-1.3-.2-1.4-.8-.2-.6.2-1.3.8-1.4 4.5-1.4 11.9-1 15.6 1.2.6.3.8 1.1.5 1.7-.3.5-1 .7-1.5.4z" />
                        </svg>
                    </div>
                </div>

                <div className="space-y-4">
                    <h1 className="text-5xl font-extrabold tracking-tighter">Spootify Manager</h1>
                    <p className="text-[#b3b3b3] text-lg font-medium">Manage inclusions and exclusions with ease.</p>
                </div>

                <button
                    onClick={handleLogin}
                    disabled={loading}
                    className="w-full bg-[#1DB954] hover:bg-[#1ed760] text-black font-bold py-4 rounded-full text-lg tracking-wide transform transition-all hover:scale-105 active:scale-95 shadow-lg"
                >
                    Connect with Spotify
                </button>
            </div>
        </div>
    );
}
