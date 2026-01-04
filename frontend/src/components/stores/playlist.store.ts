import { create } from 'zustand'
import { getPlaylist, postPlaylist, postPlaylistPublish, putPlaylistByIdRename, deletePlaylistById, postPlaylistPublishall} from '@/client';
import type { ModelPlaylistResponse } from "@/client/types.gen"


interface PlaylistState {
    data: ModelPlaylistResponse[];
    current: string;
    currentId: string;
    loading: boolean;
    publishLoading: boolean;
    publishAllLoading: boolean;
    deleteLoading: boolean;
    renameLoading: boolean;
    error: boolean;
    fetch: () => Promise<void>;
    setCurrentId: (id: string, name: string) => void;
    createPlaylist: (name: string) => Promise<void>;
    publishPlaylist: () => Promise<void>;
    publishPlaylists: () => Promise<void>;
    deletePlaylist: () => Promise<void>;
    renamePlaylist: (name: string) => Promise<void>;
}

export const usePlaylistStore = create<PlaylistState>((set, get) => ({
    data: [],
    current: "",
    currentId: "",
    loading: false,
    publishLoading: false,
    publishAllLoading: false,
    deleteLoading: false,
    renameLoading: false,
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
                body: { name } 
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
                body: { spotifyID: id } 
            });

            if (response.data) {
                set(() => ({
                    publishLoading: false
                }));
            }
        } catch (err) {
            console.error("Publishing failed", err);
            set({ publishLoading: false, error: true });
        }
    },

    publishPlaylists: async () => {
        set({ publishAllLoading: true });
        try {
            const response = await postPlaylistPublishall({
            });

            if (response.data) {
                set(() => ({
                    publishAllLoading: false
                }));
            }
        } catch (err) {
            console.error("Publishing failed", err);
            set({ publishAllLoading: false, error: true });
        }
    },
    
    deletePlaylist: async () => {
        const id = get().currentId;
        set({ deleteLoading: true });
        try {
            const response = await deletePlaylistById({
                path: { id: id } 
            });

            if (response.data) {
                set((state) => ({
                    data: state.data.filter((res) => res.spotifyID != id),
                    deleteLoading: false
                }))
            }
        } catch (err) {
            console.error("Deletion failed", err);
            set({ deleteLoading: false, error: true });
        }
    },


    renamePlaylist: async (name: string) => {
        const id = get().currentId;
        set({ renameLoading: true });
        try {
            const response = await putPlaylistByIdRename({
                path: { id: id },
                body: { name }
            });

            if (response.data) {
                set((state) => ({
                    data: state.data.map((res) => { return {...res, name: res.spotifyID == id ? name : res.name }} ),
                    renameLoading: false
                }));
            }
        } catch (err) {
            console.error("Creation failed", err);
            set({ renameLoading: false, error: true });
        }
    },
}));
