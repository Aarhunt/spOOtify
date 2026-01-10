import { create } from 'zustand'
import { getPlaylist, postPlaylist, postPlaylistPublish, putPlaylistByIdRename, deletePlaylistById, postPlaylistPublishall, postSpotifyArtistAlbums, postSearch, getPlaylistByIdInclusions, postPlaylistInclude, postPlaylistIncludeUndo, postPlaylistItem, postPlaylistItemUndo, postSpotifyAlbumTracks} from '@/client';
import type { ModelPlaylistResponse, ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface PlaylistState {
    playlistSelectionData: ModelPlaylistResponse[];
    currentPlaylistName: string;
    currentPlaylistId: string;
    selectionLoading: boolean;
    publishLoading: boolean;
    publishAllLoading: boolean;
    deleteLoading: boolean;
    renameLoading: boolean;
    error: boolean;
    playlistSearchData: ModelItemResponse[];
    artistSearchData: ModelItemResponse[];
    albumSearchData: ModelItemResponse[];
    trackSearchData: ModelItemResponse[];
    currentSearchSelectedArtist: string;
    currentSearchSelectedAlbum: string;
    searchLoading: boolean;
    albumSearchLoading: boolean;
    trackSearchLoading: boolean;
    includeLoading: boolean;
    searchType: ModelItemType;
    playlistSummaryData: ModelItemResponse[];
    artistSummaryData: ModelItemResponse[];
    albumSummaryData: ModelItemResponse[];
    trackSummaryData: ModelItemResponse[];
    artistSummaryExpandData: ModelItemResponse[];
    albumSummaryExpandData: ModelItemResponse[];
    trackSummaryExpandData: ModelItemResponse[];
    currentSummarySelectedArtist: string;
    currentSummarySelectedAlbum: string;
    summaryPlaylistsLoading: boolean;
    summaryArtistsLoading: boolean;
    summaryAlbumsLoading: boolean;
    summaryTracksLoading: boolean;
    albumLoading: boolean;
    trackLoading: boolean;
    summaryType: number;
    fetchSelectionData: () => Promise<void>;
    setCurrentPlaylist: (id: string, name: string) => void;
    createPlaylist: (name: string) => Promise<void>;
    publishPlaylist: () => Promise<void>;
    publishPlaylists: () => Promise<void>;
    deletePlaylist: () => Promise<void>;
    renamePlaylist: (name: string) => Promise<void>;
    search: (query: string) => Promise<void>;
    clearSearchData: () => void;
    getSearchAlbumsFromArtist: (artistId: string) => Promise<void>;
    getSearchTracksFromAlbum: (artistId: string) => Promise<void>;
    getSummaryAlbumsFromArtist: (artistId: string) => Promise<void>;
    getSummaryTracksFromAlbum: (artistId: string) => Promise<void>;
    setSearchType: (val: ModelItemType) => void;
    setCurrentSearchArtist: (val: string) => void;
    setCurrentSearchAlbum: (val: string) => void;
    includeItem: (itemId: string, include: boolean, type: ModelItemType) => Promise<void>;
    includePlaylist: (itemId: string) => Promise<void>;
    undoIncludeItem: (itemId: string, include: boolean, type: ModelItemType) => Promise<void>;
    undoIncludePlaylist: (itemId: string) => Promise<void>;
    summary: () => Promise<void>;
    clearSummaryData: () => Promise<void>;
    setSummaryType: (val: ModelItemType) => void;
    setCurrentSummaryArtist: (val: string) => void;
    setCurrentSummaryAlbum: (val: string) => void;
    setSummaryLoading: (loading: boolean) => void;
}

export const usePlaylistStore = create<PlaylistState>((set, get) => ({
    playlistSelectionData: [],
    currentPlaylistName: "",
    currentPlaylistId: "",
    selectionLoading: false,
    publishLoading: false,
    publishAllLoading: false,
    deleteLoading: false,
    renameLoading: false,
    error: false,
    playlistSearchData: [],
    artistSearchData: [],
    albumSearchData: [],
    trackSearchData: [],
    currentSearchSelectedArtist: "",
    currentSearchSelectedAlbum: "",
    searchLoading: false,
    albumSearchLoading: false,
    trackSearchLoading: false,
    includeLoading: false,
    searchType: 1,
    playlistSummaryData: [],
    artistSummaryData: [],
    albumSummaryData: [],
    trackSummaryData: [],
    artistSummaryExpandData: [],
    albumSummaryExpandData: [],
    trackSummaryExpandData: [],
    currentSummarySelectedArtist: "",
    currentSummarySelectedAlbum: "",
    playlistId: "", 
    summaryLoading: false,
    albumLoading: false,
    trackLoading: false,
    summaryType: 4,
    summaryPlaylistsLoading: false,
    summaryArtistsLoading: false,
    summaryAlbumsLoading: false,
    summaryTracksLoading: false,

    setCurrentPlaylist: (id, name) => set({ currentPlaylistId: id, currentPlaylistName: name}),

    fetchSelectionData: async () => {
        set({ selectionLoading: true });
        const response = await getPlaylist();
        if (response.data) {
            set({ playlistSelectionData: response.data, selectionLoading: false });
        } else {
            set({ selectionLoading: false, error: true });
        }
    },

    createPlaylist: async (name: string) => {
        set({ selectionLoading: true });
        try {
            const response = await postPlaylist({
                body: { name } 
            });

            if (response.data) {
                set((state) => ({
                    playlistSelectionData: [...state.playlistSelectionData, response.data],
                    selectionLoading: false
                }));
            }
        } catch (err) {
            console.error("Creation failed", err);
            set({ selectionLoading: false, error: true });
        }
    },

    publishPlaylist: async () => {
        set({ publishLoading: true });
        try {
            const id = get().currentPlaylistId
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
        const id = get().currentPlaylistId;
        set({ deleteLoading: true });
        try {
            const response = await deletePlaylistById({
                path: { id: id } 
            });

            if (response.data) {
                set((state) => ({
                    playlistSelectionData: state.playlistSelectionData.filter((res) => res.spotifyID != id),
                    deleteLoading: false
                }))
            }
        } catch (err) {
            console.error("Deletion failed", err);
            set({ deleteLoading: false, error: true });
        }
    },


    renamePlaylist: async (name: string) => {
        const id = get().currentPlaylistId;
        set({ renameLoading: true });
        try {
            const response = await putPlaylistByIdRename({
                path: { id: id },
                body: { name }
            });

            if (response.data) {
                set((state) => ({
                    playlistSelectionData: state.playlistSelectionData.map((res) => { return {...res, name: res.spotifyID == id ? name : res.name }} ),
                    renameLoading: false
                }));
            }
        } catch (err) {
            console.error("Creation failed", err);
            set({ renameLoading: false, error: true });
        }
    },

    clearSearchData: () => { set({ playlistSearchData: [], artistSearchData: [], albumSearchData: [], trackSearchData: [] })},

    search: async (query: string) => {
        set({ searchLoading: true});
        get().clearSearchData();

        const currentPlaylistId = get().currentPlaylistId;
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
                set({ artistSearchData: response.data });
            } else if (searchType == 2) {
                set({ albumSearchData: response.data });
            } else if (searchType == 3) {
                set({ trackSearchData: response.data });
            } else if (searchType == 0) {
                set({ playlistSearchData: response.data })
            }
            set({ searchLoading: false, error: false })
        } else {
            set({ searchLoading: false, error: true });
        }
    },

    setSearchType: (val: ModelItemType) => {
        set({ searchType: val });
        get().clearSearchData();
    },

    setCurrentSearchArtist: (val: string) => set({ currentSearchSelectedArtist: val }),
    setCurrentSearchAlbum: (val: string) => set({ currentSearchSelectedAlbum: val }),

    getSearchAlbumsFromArtist: async (artistId: string) => {
        set({ albumSearchLoading: true, albumSearchData: [], trackSearchData: [] })

        const response = await postSpotifyArtistAlbums({
            body: {
                parentid: artistId,
                playlistid: get().currentPlaylistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ albumSearchData: response.data, albumSearchLoading: false, error: false });
        } else {
            set({ albumSearchLoading: false, error: true });
        }
    },

    getSearchTracksFromAlbum: async (albumId: string) => {
        set({ trackSearchLoading: true, trackSearchData: [] })

        const response = await postSpotifyAlbumTracks({
            body: {
                parentid: albumId,
                playlistid: get().currentPlaylistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ trackSearchData: response.data, trackSearchLoading: false, error: false });
        } else {
            set({ trackSearchLoading: false, error: true });
        }
    },

    getSummaryAlbumsFromArtist: async (artistId: string) => {
        set({ summaryAlbumsLoading: true, albumSummaryExpandData: [], trackSummaryExpandData: [] })

        const response = await postSpotifyArtistAlbums({
            body: {
                parentid: artistId,
                playlistid: get().currentPlaylistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ albumSummaryExpandData: response.data, summaryAlbumsLoading: false, error: false });
        } else {
            set({ summaryAlbumsLoading: false, error: true });
        }
    },

    getSummaryTracksFromAlbum: async (albumId: string) => {
        set({ summaryTracksLoading: true, trackSummaryExpandData: [] })

        const response = await postSpotifyAlbumTracks({
            body: {
                parentid: albumId,
                playlistid: get().currentPlaylistId,
                type: 1,
            }
        })

        if (response.data) {
            set({ trackSummaryExpandData: response.data, summaryTracksLoading: false, error: false });
        } else {
            set({ summaryTracksLoading: false, error: true });
        }
    },

    includeItem: async (itemId: string, include: boolean, type: ModelItemType) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistItem({
                body: {
                    include: include,
                    playlistid: get().currentPlaylistId,
                    spotid: itemId,
                    type: type,
                }
            });

            if (response.data) {
                const targetValue = include ? 1 : 2;
                const proxyValue = include ? 3 : 4;
                const state = get();

                const searchDataKey = type === 1 ? 'artistSearchData' : type === 2 ? 'albumSearchData' : 'trackSearchData'
                const summaryDataKey = type === 2 ? 'albumSummaryExpandData' : type === 3 ? 'trackSummaryExpandData' : ''
                const searchList = [...state[searchDataKey]];
                const summaryList = summaryDataKey != '' ? [...state[summaryDataKey]] : []

                // Update the frontend element that is now included
                var newItem;
                {
                    const foundIndex = searchList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) searchList[foundIndex] = { ...searchList[foundIndex], included: targetValue };
                    newItem = searchList[foundIndex]
                }

                // If there is a frontend element in an expanded column in summary, update that
                {
                    const foundIndex = summaryList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) summaryList[foundIndex] = { ...summaryList[foundIndex], included: targetValue}
                }

                const newState: Partial<PlaylistState> = {
                    [searchDataKey]: searchList,
                    [summaryDataKey]: summaryList,
                    includeLoading: false,
                    error: false
                };

                const albumExclusions = ["live", "acoustic", "instrumental"]
                const trackExclusions = ["- live", "- acoustic", "- instrumental"]

                // Update proxy elements if they exist
                if (type === 1) {
                    if (itemId == get().currentSearchSelectedArtist) {
                        newState.albumSearchData = state.albumSearchData.map(a => a.inclusionByProxy ? a.included == 1 || a.included == 2 ? {...a} : { ...a, included: proxyValue } : {...a} );
                        newState.albumSearchData = newState.albumSearchData.map(a => a.inclusionByProxy ? a.included != 1 && matchSubstrings(a.name!, albumExclusions) ? {...a, included: 2} : {...a} : {...a});
                        // Proxy element of summary view
                        newState.albumSummaryExpandData = state.albumSummaryExpandData.map(a => a.included == 1 || a.included == 2 ? {...a} : { ...a, included: proxyValue } );
                    }
                    // Proxy element of search view
                    // Edit summary view 
                    newState.artistSummaryData = include ? [...state.artistSummaryData, newItem] : state.artistSummaryData.filter(a => itemId != a.spotifyID);
                } else if (type === 2) {
                    if (itemId == get().currentSearchSelectedAlbum) {
                        newState.trackSearchData = state.trackSearchData.map(t => t.included == 1 || t.included == 2 ? {...t} : { ...t, included: proxyValue } );
                        newState.trackSearchData = newState.trackSearchData.map(t => t.included != 1 && matchSubstrings(t.name!, trackExclusions) ? {...t, included: 2} : {...t});
                        newState.trackSummaryExpandData = state.trackSummaryExpandData.map(t => t.included == 1 || t.included == 2 ? {...t} : { ...t, included: proxyValue });
                    }
                    newState.albumSummaryData = include ? [...state.albumSummaryData, newItem] : state.albumSummaryData.filter(a => itemId != a.spotifyID);
                } else if (type === 3) {
                    newState.trackSummaryData = include ? [...state.trackSummaryData, newItem] : state.trackSummaryData.filter(a => itemId != a.spotifyID);
                }

                set(newState);
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },

    undoIncludeItem: async (itemId: string, include: boolean, type: ModelItemType) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistItemUndo({
                body: {
                    include: include,
                    playlistid: get().currentPlaylistId,
                    spotid: itemId,
                    type: type,
                }
            });

            if (response.data) {
                const targetValue = 0;
                const state = get();

                const dataKey = type === 1 ? 'artistSearchData' : type === 2 ? 'albumSearchData' : 'trackSearchData';
                const summaryDataKey = type === 2 ? 'albumSummaryExpandData' : type === 3 ? 'trackSummaryExpandData' : ''
                const currentList = [...state[dataKey]];
                const summaryList = summaryDataKey != '' ? [...state[summaryDataKey]] : []

                // If there is a frontend element in search, update that
                {
                    const foundIndex = currentList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) currentList[foundIndex] = { ...currentList[foundIndex], included: targetValue };
                }

                // If there is a frontend element in an expanded column in summary, update that
                {
                    const foundIndex = summaryList.findIndex(i => i.spotifyID === itemId);
                    if (foundIndex !== -1) summaryList[foundIndex] = { ...summaryList[foundIndex], included: targetValue}
                }

                const newState: Partial<PlaylistState> = {
                    [dataKey]: currentList,
                    [summaryDataKey]: summaryList,
                    includeLoading: false,
                    error: false
                };

                const albumExclusions = ["live", "acoustic", "instrumental"]
                const trackExclusions = ["- live", "- acoustic", "- instrumental", "- single"]

                if (type === 1) {
                    if (itemId == get().currentSearchSelectedArtist) {
                        newState.albumSearchData = state.albumSearchData.map(a => a.included == 3 || a.included == 4 ? { ...a, included: targetValue } : {...a});
                        newState.albumSearchData = newState.albumSearchData.map(a => a.included == 2 && matchSubstrings(a.name!, albumExclusions) ? {...a, included: targetValue} : {...a})
                        newState.albumSummaryExpandData = state.albumSummaryExpandData.map(a => a.included == 1 || a.included == 2 ? {...a} : { ...a, included: targetValue });
                    }
                    newState.artistSummaryData = state.artistSummaryData.filter(a => itemId != a.spotifyID);
                    if (itemId == get().currentSummarySelectedArtist) {newState.albumSummaryExpandData = [], newState.trackSummaryExpandData = []}
                } else if (type === 2) {
                    if (itemId == get().currentSearchSelectedAlbum) {
                        newState.trackSearchData = state.trackSearchData.map(t => t.included == 3 || t.included == 4 ? { ...t, included: targetValue } : {...t});
                        newState.trackSearchData = newState.trackSearchData.map(t => t.included == 2 && matchSubstrings(t.name!, trackExclusions) ? {...t, included: targetValue} : {...t})
                        newState.trackSummaryExpandData = state.trackSummaryExpandData.map(t => t.included == 1 || t.included == 2 ? {...t} : { ...t, included: targetValue });
                        newState.trackSummaryExpandData = []
                    }
                    newState.albumSummaryData = state.albumSummaryData.filter(a => itemId != a.spotifyID);
                } else if (type === 3) {
                    newState.trackSummaryData = state.trackSummaryData.filter(a => itemId != a.spotifyID);
                }

                set(newState);
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },

    includePlaylist: async (itemId: string) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistInclude({
                body: {
                    pspotid: get().currentPlaylistId,
                    cspotid: itemId,
                }
            });

            if (response.data) {
                const targetValue = 1;
                const state = get();
                const playlistSearchData = state.playlistSearchData

                // If there is a frontend element in search, update that
                const foundIndex = playlistSearchData.findIndex(i => i.spotifyID === itemId);
                if (foundIndex !== -1) playlistSearchData[foundIndex] = {...playlistSearchData[foundIndex], included: targetValue}
                const newItem = playlistSearchData[foundIndex]

                // Add the playlist to the summary view
                const playlistSummaryData = [...state.playlistSummaryData, newItem]

                set({ includeLoading: false, playlistSearchData: playlistSearchData, playlistSummaryData: playlistSummaryData });
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },

    undoIncludePlaylist: async (itemId: string) => {
        set({ includeLoading: true });

        try {
            const response = await postPlaylistIncludeUndo({
                body: {
                    pspotid: get().currentPlaylistId,
                    cspotid: itemId,
                }
            });

            if (response.data) {
                const targetValue = 0;
                const state = get();
                const playlistSearchData = state.playlistSearchData

                // If there is a frontend element in search, update that
                const foundIndex = playlistSearchData.findIndex(i => i.spotifyID === itemId);
                if (foundIndex !== -1) playlistSearchData[foundIndex] = { ...playlistSearchData[foundIndex], included: targetValue };

                const playlistSummaryData = state.playlistSummaryData.filter(a => itemId != a.spotifyID);
                set({ includeLoading: false, playlistSearchData: playlistSearchData, playlistSummaryData: playlistSummaryData});
            }
        } catch (err) {
            set({ includeLoading: false, error: true });
        }
    },

    clearSummaryData: async () => {
        set({ albumSummaryExpandData: [], trackSummaryExpandData: [], artistSummaryData: [], albumSummaryData: [], trackSummaryData: []})
    },

    summary: async () => {
        get().setSummaryLoading(true)
        get().clearSummaryData();

        const currentPlaylistId = get().currentPlaylistId;

        if (currentPlaylistId == "") { 
            get().setSummaryLoading(false)
            set({ error: true });
            return
        }

        try {
            const response = await getPlaylistByIdInclusions({ 
                path: { id: currentPlaylistId } 
            });

            if (response.data) {
                const playlists: ModelItemResponse[] = [];
                const artists: ModelItemResponse[] = [];
                const albums: ModelItemResponse[] = [];
                const tracks: ModelItemResponse[] = [];

                response.data.forEach((item) => {

                    switch (item.itemType) {
                        case 0: // PlaylistItem
                            playlists.push(item);
                            break;
                        case 1: // Artist
                            artists.push(item);
                            break;
                        case 2: // Album
                            albums.push(item);
                            break;
                        case 3: // Track
                            tracks.push(item);
                            break;
                        default:
                            break;
                    }
                });

                set({ 
                    playlistSummaryData: playlists,
                    artistSummaryData: artists, 
                    albumSummaryData: albums, 
                    trackSummaryData: tracks, 
                    error: false 
                });
                get().setSummaryLoading(false);
            } else {
                set({ error: true });
            get().setSummaryLoading(false)
            }
        } catch (err) {
            set({ error: true });
            get().setSummaryLoading(false)
        }
    },

    setSummaryType: (val: ModelItemType) => set({ summaryType: val, artistSummaryData: [], albumSummaryData: [], trackSummaryData: [], albumSummaryExpandData: [], trackSummaryExpandData: [] }),
    setCurrentSummaryArtist: (val: string) => set({ currentSummarySelectedArtist: val }),
    setCurrentSummaryAlbum: (val: string) => set({ currentSummarySelectedAlbum: val }),
    setSummaryLoading: (loading: boolean) => set({ summaryPlaylistsLoading: loading, summaryArtistsLoading: loading, summaryAlbumsLoading: loading, summaryTracksLoading: loading})
}));

function matchSubstrings(s: string, match: string[]): boolean {
    const lowerStr = s.toLowerCase()
    return match.some(el => lowerStr.includes(el))
}
