import { create } from 'zustand'
import { postSearch, postSpotifyArtistAlbums } from '@/client';
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface SearchState {
    searchData: ModelItemResponse[];
    playlistId: string;
    loading: boolean;
    error: boolean;
    setPlaylistId: (val: string) => void;
    search: (query: string, type: ModelItemType) => Promise<void>;
    search: () => Promise<void>;
}

export const useSearchStore = create<SearchState>((set, get) => ({
    searchData: [],
    playlistId: "", 
    loading: false,
    error: false,

    search: async (query: string, type: ModelItemType) => {
        set({ loading: true });

        const currentPlaylistId = get().playlistId;

        const response = await postSearch({
            body: {
                playlistid: currentPlaylistId,
                query: query,
                type: type,
            }
        });

        if (response.data) {
            set({ searchData: response.data, loading: false, error: false });
        } else {
            set({ loading: false, error: true });
        }
    },
    
    setPlaylistId: (id: string) => set({ playlistId: id }),

    getAlbumsFromArtist: async (artistId: string) => {
        set({ loading: true })

        const currentPlaylistId = get().playlistId;

        const response = await postSpotifyArtistAlbums({
            body: {
                artistid: artistId,
                playlistid: get().playlistId,
                type: 1,
            }
        })
    }
}));
