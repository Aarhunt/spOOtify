import { create } from 'zustand'
import { postSearch, postSpotifyAlbumTracks, postSpotifyArtistAlbums, postPlaylistItem, postPlaylistItemUndo, postPlaylistInclude, postPlaylistIncludeUndo } from '@/client';
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface SearchState {
    playlistData: ModelItemResponse[];
    artistData: ModelItemResponse[];
    albumData: ModelItemResponse[];
    trackData: ModelItemResponse[];
    playlistId: string;
    currentArtist: string;
    currentAlbum: string;
    searchLoading: boolean;
    albumLoading: boolean;
    trackLoading: boolean;
    includeLoading: boolean;
    error: boolean;
    searchType: ModelItemType;
    setPlaylistId: (val: string) => void;
    search: (query: string) => Promise<void>;
    clearData: () => void;
    getAlbumsFromArtist: (artistId: string) => Promise<void>;
    getTracksFromAlbum: (artistId: string) => Promise<void>;
    setSearchType: (val: ModelItemType) => void;
    setCurrentArtist: (val: string) => void;
    setCurrentAlbum: (val: string) => void;
    includeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
    includePlaylist: (itemId: string, index: number) => Promise<void>;
    undoIncludeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
    undoIncludePlaylist: (itemId: string, index: number) => Promise<void>;
}

export const useSearchStore = create<SearchState>((set, get) => ({
    playlistData: [],
    artistData: [],
    albumData: [],
    trackData: [],
    currentArtist: "",
    currentAlbum: "",
    playlistId: "", 
    searchLoading: false,
    albumLoading: false,
    trackLoading: false,
    includeLoading: false,
    error: false,
    searchType: 1,

    clearData: () => { set({ playlistData: [], artistData: [], albumData: [], trackData: [] })},

    search: async (query: string) => {
        set({ searchLoading: true});
        get().clearData();

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
            } else if (searchType == 0) {
                set({ playlistData: response.data })
            }
            set({ searchLoading: false, error: false })
        } else {
            set({ searchLoading: false, error: true });
        }
    },

    setSearchType: (val: ModelItemType) => {
        set({ searchType: val });
        get().clearData();
    },
    setCurrentArtist: (val: string) => set({ currentArtist: val }),
    setCurrentAlbum: (val: string) => set({ currentAlbum: val }),

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
                const targetValue = include ? 1 : 2;
                const proxyValue = include ? 3 : 4;
                const state = get();

                const dataKey = type === 1 ? 'artistData' : type === 2 ? 'albumData' : 'trackData';
                const currentList = [...state[dataKey]];

                if (currentList[index]?.spotifyID === itemId) {
                    currentList[index] = { ...currentList[index], included: targetValue };
                } else {
                    const foundIndex = currentList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) currentList[foundIndex] = { ...currentList[foundIndex], included: targetValue };
                }

                const newState: Partial<SearchState> = {
                    [dataKey]: currentList,
                    includeLoading: false,
                    error: false
                };

                if (type === 1) {
                    newState.albumData = state.albumData.map(a => itemId == get().currentArtist ? (a.included == 1 || a.included == 2 ? {...a} : { ...a, included: proxyValue } ) : {...a});
                    // newState.trackData = state.trackData.map(t => itemId == get().currentArtist ? (t.included == 1 || t.included == 2 ? {...t} : { ...t, included: proxyValue } ) : {});
                } else if (type === 2) {
                    newState.trackData = state.trackData.map(t => itemId == get().currentAlbum ? (t.included == 1 || t.included == 2 ? {...t} : { ...t, included: proxyValue } ) : {...t});
                }

                set(newState);
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },


    includePlaylist: async (itemId: string, index: number) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistInclude({
                body: {
                    pspotid: get().playlistId,
                    cspotid: itemId,
                }
            });

            if (response.data) {
                const currentList = get().playlistData;
                const targetValue = 1;

                if (currentList[index]?.spotifyID === itemId) {
                    currentList[index] = { ...currentList[index], included: targetValue };
                } else {
                    const foundIndex = currentList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) currentList[foundIndex] = { ...currentList[foundIndex], included: targetValue };
                }
                set({ includeLoading: false, playlistData: currentList });
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
                const targetValue = 0;
                const state = get();

                const dataKey = type === 1 ? 'artistData' : type === 2 ? 'albumData' : 'trackData';
                const currentList = [...state[dataKey]];

                if (currentList[index]?.spotifyID === itemId) {
                    currentList[index] = { ...currentList[index], included: targetValue };
                } else {
                    const foundIndex = currentList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) currentList[foundIndex] = { ...currentList[foundIndex], included: targetValue };
                }

                const newState: Partial<SearchState> = {
                    [dataKey]: currentList,
                    includeLoading: false,
                    error: false
                };

                if (type === 1) {
                    newState.albumData = state.albumData.map(a => itemId == get().currentArtist ? (a.included == 3 || a.included == 4 ? { ...a, included: targetValue } : {...a}) : {...a});
                    // newState.trackData = state.trackData.map(t => itemId == get().currentArtist ? (t.included == 3 || t.included == 4 ? { ...t, included: targetValue } : {...t}) : {});
                } else if (type === 2) {
                    newState.trackData = state.trackData.map(t => itemId == get().currentArtist ? (t.included == 3 || t.included == 4 ? { ...t, included: targetValue } : {...t}) : {...t});
                }

                set(newState);
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },


    undoIncludePlaylist: async (itemId: string, index: number) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistIncludeUndo({
                body: {
                    pspotid: get().playlistId,
                    cspotid: itemId,
                }
            });

            if (response.data) {
                const currentList = get().playlistData;
                const targetValue = 0;

                if (currentList[index]?.spotifyID === itemId) {
                    currentList[index] = { ...currentList[index], included: targetValue };
                } else {
                    const foundIndex = currentList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) currentList[foundIndex] = { ...currentList[foundIndex], included: targetValue };
                }
                set({ includeLoading: false, playlistData: currentList });
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },
}));
