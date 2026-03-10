import { useEffect, useState } from 'react'
import { getAuthStatus } from '@/client'
import './App.css'
import './styles/globals.css'
import { client } from '@/client/client.gen';

// Page Imports
import Playlist from './components/pages/Playlist'
import Search from './components/pages/Search'
import Summary from './components/pages/Summary'
import Login from './components/pages/Login'
import Flow from './components/pages/Flow'
import { Home, Share2 } from 'lucide-react';

// Service Import
function App() {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
    const [isLoading, setIsLoading] = useState<boolean>(true)
    const [view, setView] = useState<'summary' | 'graph'>('summary');

    client.setConfig({
        baseUrl: import.meta.env.VITE_API_URL,
    });

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const response = await getAuthStatus();
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
                {/* SIDEBAR */}
                <aside className="w-80 bg-[#121212] rounded-lg border border-[#282828] h-full overflow-hidden flex flex-col">
                    {/* NAV SECTION */}
                    <nav className="p-4 space-y-2">
                        <button 
                            onClick={() => setView('summary')}
                            className={`w-full flex items-center gap-4 px-4 py-3 rounded-md font-bold transition-colors ${view === 'summary' ? 'text-white bg-[#282828]' : 'text-gray-400 hover:text-white'}`}
                        >
                            <Home size={24} />
                            Summary
                        </button>
                        <button 
                            onClick={() => setView('graph')}
                            className={`w-full flex items-center gap-4 px-4 py-3 rounded-md font-bold transition-colors ${view === 'graph' ? 'text-white bg-[#282828]' : 'text-gray-400 hover:text-white'}`}
                        >
                            <Share2 size={24} />
                            Graph View
                        </button>
                    </nav>

                    <hr className="border-[#282828] mx-4 mb-2" />

                    <h2 className="text-xl font-bold p-6 pb-2 text-white">Your Library</h2>
                    <div className="flex-1 overflow-hidden p-4 pt-0"> 
                        <Playlist />
                    </div>
                </aside>

                {/* MAIN CONTENT */}
                <main className={`flex-1 bg-gradient-to-b from-[#222222] to-[#121212] rounded-lg border border-[#282828] flex flex-col ${view === 'graph' ? 'overflow-hidden' : 'overflow-y-auto'}`}>
                    <div className="flex-1 relative">
                        {view === 'summary' ? (
                            <>
                            <header className="p-4 flex justify-between items-center bg-transparent z-10">
                                <Search />
                            </header>
                            <section className="p-4 mt-8 animate-in fade-in duration-500">
                                <Summary /> 
                            </section>
                            </>
                        ) : (
                            <section className="absolute inset-0 animate-in zoom-in-95 duration-300">
                                <Flow />
                            </section>
                        )}
                    </div>
                </main>
            </div>
        </div>
    );
}



export default App
