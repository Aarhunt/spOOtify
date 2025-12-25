import { create } from 'zustand'
import { getPlaylist } from '@/client';
import type { ModelPlaylist } from "@/client/types.gen"


interface PlaylistState {
    data: ModelPlaylist[];
    loading: boolean;
    fetch: () => Promise<void>;
}

export const usePlaylistStore = create<PlaylistState>((set) => ({
    data: [],
    loading: false,
    fetch: async () => {
        set({ loading: true });
        try {
            // "getPlaylists" is typed! It knows it returns Playlist[]
            const response = await getPlaylist(); 
            set({ data: response.data, loading: false });
        } catch (error) {
            set({ loading: false });
        }
    }
}));
