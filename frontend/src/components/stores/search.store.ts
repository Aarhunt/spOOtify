import { create } from 'zustand'
import { postSearch } from '@/client';
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface SearchState {
    data: ModelItemResponse[];
    playlistId: string;
    loading: boolean;
    error: boolean;
    setPlaylistId: (val: string) => void;
    search: (query: string, type: ModelItemType) => Promise<void>;
}

export const useSearchStore = create<SearchState>((set, get) => ({
    data: [],
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
            set({ data: response.data, loading: false, error: false });
        } else {
            set({ loading: false, error: true });
        }
    },
    
    setPlaylistId: (id: string) => set({ playlistId: id }),
}));
