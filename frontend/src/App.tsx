import { useEffect, useState } from 'react'
import { getAuthStatus } from '@/client'
import './App.css'
import './styles/globals.css'

// Page Imports
import Playlist from './components/pages/Playlist'
import Search from './components/pages/Search'
import Summary from './components/pages/Summary'
import Login from './components/pages/Login'

// Service Import
function App() {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
    const [isLoading, setIsLoading] = useState<boolean>(true)

    useEffect(() => {
        const checkAuth = async () => {
            try {
                // Call the new backend endpoint
                const response = await getAuthStatus();
                // Ensure we handle the response object correctly
                if (response.data && response.data.authenticated) {
                    setIsAuthenticated(true);
                }
            } catch (error) {
                console.log("Not authenticated or backend offline");
                setIsAuthenticated(false);
            } finally {
                setIsLoading(false);
            }
        };

        checkAuth();
    }, []);

    if (isLoading) {
        return <div className="min-h-screen bg-black text-white flex items-center justify-center">Loading...</div>
    }

    if (!isAuthenticated) {
        return <Login />
    }

return (
        <div className="h-screen w-screen flex flex-col bg-black overflow-hidden font-sans">
            <div className="flex flex-1 overflow-hidden p-2 gap-2">
<aside className="w-80 bg-[#121212] rounded-lg border border-[#282828] h-full overflow-hidden flex flex-col">
    <h2 className="text-xl font-bold p-6 pb-2 text-white">Your Library</h2>
    <div className="flex-1 overflow-hidden p-4 pt-0"> 
        <Playlist />
    </div>
</aside>

                <main className="flex-1 bg-gradient-to-b from-[#222222] to-[#121212] rounded-lg overflow-y-auto p-4 border border-[#282828]">
                    <header className="mb-10 flex justify-between items-center">
                        <Search />
                    </header>
                    
                    <section className="mt-8">
                        <Summary />
                    </section>
                </main>
            </div>
        </div>
    );
}

export default App
