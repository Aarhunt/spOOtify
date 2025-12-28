import { create } from 'zustand'
import { getPlaylist, postPlaylist, postPlaylistPublish } from '@/client';
import type { ModelPlaylistResponse } from "@/client/types.gen"


interface PlaylistState {
    data: ModelPlaylistResponse[];
    current: string;
    loading: boolean;
    error: boolean;
    fetch: () => Promise<void>;
    setCurrent: (name: string) => void;
    createPlaylist: (name: string) => Promise<void>;
}

export const usePlaylistStore = create<PlaylistState>((set) => ({
    data: [],
    current: "",
    loading: false,
    error: false,

    setCurrent: (name) => set({ current: name }),

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
    },

    publishPlaylist: async (id: string) => {
        set({ loading: true });
        try {
            const response = await postPlaylistPublish({
                body: { spotifyID: id } // Hey API structure
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
    },


}));
