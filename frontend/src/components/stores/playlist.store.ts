import { create } from 'zustand'
import { getPlaylist, postPlaylist, postPlaylistPublish } from '@/client';
import type { ModelPlaylistResponse } from "@/client/types.gen"


interface PlaylistState {
    data: ModelPlaylistResponse[];
    current: string;
    currentId: string;
    loading: boolean;
    publishLoading: boolean;
    error: boolean;
    fetch: () => Promise<void>;
    setCurrentId: (id: string, name: string) => void;
    createPlaylist: (name: string) => Promise<void>;
    publishPlaylist: () => Promise<void>;
}

export const usePlaylistStore = create<PlaylistState>((set, get) => ({
    data: [],
    current: "",
    currentId: "",
    loading: false,
    publishLoading: false,
    error: false,

    setCurrentId: (id, name) => set({ currentId: id, current: name}),

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

    publishPlaylist: async () => {
        set({ publishLoading: true });
        try {
            const id = get().currentId
            const response = await postPlaylistPublish({
                body: { spotifyID: id } // Hey API structure
            });

            if (response.data) {
                set((state) => ({
                    data: [...state.data, response.data],
                    publishLoading: false
                }));
            }
        } catch (err) {
            console.error("Creation failed", err);
            set({ publishLoading: false, error: true });
        }
    },


}));
