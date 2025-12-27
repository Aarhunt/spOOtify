import { create } from 'zustand'
import { postSearch, postSpotifyAlbumTracks, postSpotifyArtistAlbums, postPlaylistItem, postPlaylistItemUndo } from '@/client';
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface SearchState {
    artistData: ModelItemResponse[];
    albumData: ModelItemResponse[];
    trackData: ModelItemResponse[];
    playlistId: string;
    searchLoading: boolean;
    albumLoading: boolean;
    trackLoading: boolean;
    includeLoading: boolean;
    error: boolean;
    searchType: ModelItemType;
    setPlaylistId: (val: string) => void;
    search: (query: string) => Promise<void>;
    getAlbumsFromArtist: (artistId: string) => Promise<void>;
    getTracksFromAlbum: (artistId: string) => Promise<void>;
    setSearchType: (val: ModelItemType) => void;
    includeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
    undoIncludeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
}

export const useSearchStore = create<SearchState>((set, get) => ({
    artistData: [],
    albumData: [],
    trackData: [],
    playlistId: "", 
    searchLoading: false,
    albumLoading: false,
    trackLoading: false,
    includeLoading: false,
    error: false,
    searchType: 1,

    search: async (query: string) => {
        set({ searchLoading: true, artistData: [], albumData: [], trackData: [] });

        const currentPlaylistId = get().playlistId;
        const searchType = get().searchType;

        const response = await postSearch({
            body: {
                playlistid: currentPlaylistId,
                query: query,
                type: searchType,
            }
        });

        if (response.data) {
            if (searchType == 1) {
                set({ artistData: response.data });
            } else if (searchType == 2) {
                set({ albumData: response.data });
            } else if (searchType == 3) {
                set({ trackData: response.data });
            }
            set({ searchLoading: false, error: false })
        } else {
            set({ searchLoading: false, error: true });
        }
    },

    setSearchType: (val: ModelItemType) => set({ searchType: val, artistData: [], albumData: [], trackData: []}),

        setPlaylistId: (id: string) => set({ playlistId: id }),

        getAlbumsFromArtist: async (artistId: string) => {
        set({ albumLoading: true, trackData: [] })

        const response = await postSpotifyArtistAlbums({
            body: {
                parentid: artistId,
                playlistid: get().playlistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ albumData: response.data, albumLoading: false, error: false });
        } else {
            set({ albumLoading: false, error: true });
        }
    },

    getTracksFromAlbum: async (albumId: string) => {
        set({ trackLoading: true })

        const response = await postSpotifyAlbumTracks({
            body: {
                parentid: albumId,
                playlistid: get().playlistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ trackData: response.data, trackLoading: false, error: false });
        } else {
            set({ trackLoading: false, error: true });
        }
    },

    includeItem: async (itemId: string, include: boolean, type: ModelItemType, index: number) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistItem({
                body: {
                    include: include,
                    playlistid: get().playlistId,
                    spotid: itemId,
                    type: type,
                }
            });

            if (response.data) {
                const dataKey = 
                    type === 1 ? 'artistData' : 
                    type === 2 ? 'albumData' : 
                    'trackData';
                const currentList = get()[dataKey];

                if (currentList[index] && currentList[index].spotifyID === itemId) {
                    const updatedList = [...currentList];

                    updatedList[index] = { 
                        ...updatedList[index], 
                        included: include ? 1 : 2 
                    };

                    set({ 
                        [dataKey]: updatedList, 
                        includeLoading: false, 
                        error: false 
                    });
                } else {
                    console.warn("Index mismatch, falling back to ID search");
                    const fallbackList = currentList.map(item => 
                                                         item.spotifyID === itemId ? { ...item, included: include ? 1 : 2 } : item
                                                        );
                                                        set({ [dataKey]: fallbackList, includeLoading: false });
                }
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },


    undoIncludeItem: async (itemId: string, include: boolean, type: ModelItemType, index: number) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistItemUndo({
                body: {
                    include: include,
                    playlistid: get().playlistId,
                    spotid: itemId,
                    type: type,
                }
            });

            if (response.data) {
                const dataKey = 
                    type === 1 ? 'artistData' : 
                    type === 2 ? 'albumData' : 
                    'trackData';

                const currentList = get()[dataKey];

                // Safety check for index optimization
                if (currentList[index] && currentList[index].spotifyID === itemId) {
                    const updatedList = [...currentList];

                    // FORCE state to 0 (Nothing)
                    updatedList[index] = { 
                        ...updatedList[index], 
                        included: 0 
                    };

                    set({ 
                        [dataKey]: updatedList, 
                        includeLoading: false, 
                        error: false 
                    });
                } else {
                    // Fallback ID search if index is stale
                    const fallbackList = currentList.map(item => 
                                                         item.spotifyID === itemId ? { ...item, included: 0 } : item
                                                        );
                                                        set({ [dataKey]: fallbackList, includeLoading: false });
                }
            }
        } catch (err) {
            console.error("Undo failed", err);
            set({ includeLoading: false, error: true });
        }
    },
}));
