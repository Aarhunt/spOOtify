import { create } from 'zustand'
import { getPlaylist, postPlaylist } from '@/client';
import type { ModelPlaylist } from "@/client/types.gen"


interface PlaylistState {
    data: ModelPlaylist[];
    loading: boolean;
    error: boolean;
    fetch: () => Promise<void>;
    createPlaylist: (name: string) => Promise<void>;
}

export const usePlaylistStore = create<PlaylistState>((set) => ({
    data: [],
    loading: false,
    error: false,

    fetch: async () => {
        set({ loading: true });
        const response = await getPlaylist();
        if (response.data) {
            set({ data: response.data, loading: false });
        } else {
            set({ loading: false, error: true });
        }
    },

    createPlaylist: async (name: string) => {
        set({ loading: true });
        try {
            const response = await postPlaylist({
                body: { name } // Hey API structure
            });

            if (response.data) {
                set((state) => ({
                    data: [...state.data, response.data],
                    loading: false
                }));
            }
        } catch (err) {
            console.error("Creation failed", err);
            set({ loading: false, error: true });
        }
    }
}));
