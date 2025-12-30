import { create } from 'zustand'
import { postSpotifyAlbumTracks, postSpotifyArtistAlbums, postPlaylistItem, postPlaylistItemUndo, postSearch, getPlaylistByIdInclusions } from '@/client';
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen"


interface SummaryState {
    mainPlaylistData: ModelItemResponse[];
    mainArtistData: ModelItemResponse[];
    mainAlbumData: ModelItemResponse[];
    mainTrackData: ModelItemResponse[];
    artistData: ModelItemResponse[];
    albumData: ModelItemResponse[];
    trackData: ModelItemResponse[];
    playlistId: string;
    currentArtist: string;
    currentAlbum: string;
    summaryLoading: boolean;
    albumLoading: boolean;
    trackLoading: boolean;
    includeLoading: boolean;
    error: boolean;
    summaryType: ModelItemType;
    setPlaylistIdSummary: (val: string) => void;
    summary: () => Promise<void>;
    getAlbumsFromArtist: (artistId: string) => Promise<void>;
    getTracksFromAlbum: (artistId: string) => Promise<void>;
    setSummaryType: (val: ModelItemType) => void;
    setCurrentArtist: (val: string) => void;
    setCurrentAlbum: (val: string) => void;
    includeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
    undoIncludeItem: (itemId: string, include: boolean, type: ModelItemType, index: number) => Promise<void>;
}

export const useSummaryStore = create<SummaryState>((set, get) => ({
    mainPlaylistData: [],
    mainArtistData: [],
    mainAlbumData: [],
    mainTrackData: [],
    artistData: [],
    albumData: [],
    trackData: [],
    currentArtist: "",
    currentAlbum: "",
    playlistId: "", 
    summaryLoading: false,
    albumLoading: false,
    trackLoading: false,
    includeLoading: false,
    error: false,
    summaryType: 1,

    summary: async () => {
        set({ summaryLoading: true, artistData: [], albumData: [], trackData: [] });

        const currentPlaylistId = get().playlistId;

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
                    mainPlaylistData: playlists,
                    mainArtistData: artists, 
                    mainAlbumData: albums, 
                    mainTrackData: tracks, 
                    summaryLoading: false, 
                    error: false 
                });
            } else {
                set({ summaryLoading: false, error: true });
            }
        } catch (err) {
            set({ summaryLoading: false, error: true });
        }
    },

    setSummaryType: (val: ModelItemType) => set({ summaryType: val, artistData: [], albumData: [], trackData: []}),
    setCurrentArtist: (val: string) => set({ currentArtist: val }),
    setCurrentAlbum: (val: string) => set({ currentAlbum: val }),

    setPlaylistIdSummary: (id: string) => set({ playlistId: id }),

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

                const newState: Partial<SummaryState> = {
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

                const newState: Partial<SummaryState> = {
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
    },}));
