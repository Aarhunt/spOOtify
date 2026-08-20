"use client"

import * as React from "react"
import { CheckIcon, ChevronsUpDownIcon, Plus, ListMusic, PenLine, Trash } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import { usePlaylistStore } from "@/components/stores/playlist.store"
import { postSearch, getSpotifyArtistByIdGenres, postPlaylistItem, postPlaylistItemUndo, postPlaylistInclude, postPlaylistIncludeUndo} from "@/client"
import type { ModelItemResponse } from "@/client/types.gen"
import { SearchResultItem } from "@/components/pages/Search"

import { Spinner } from "@/components/ui/spinner"
import { GenrePill } from "../utils/GenrePill"


const sidebarActionClass = "w-full justify-start gap-3 px-3 py-6 text-gray-400 hover:text-white hover:bg-white/10 border-none transition-all";

export default function Playlist() {
    const fetchPlaylists = usePlaylistStore((state) => state.fetchSelectionData);

    React.useEffect(() => {
        fetchPlaylists();
    }, [fetchPlaylists]);

    return (
        <div className="flex flex-col h-full w-full bg-transparent overflow-hidden">
            
            <div className="flex-none pb-6">
                <PlaylistSearch />
            </div>

            <div className="flex-1 overflow-y-auto min-h-0 pr-2 custom-scrollbar">
                <div className="flex flex-col gap-1">
                    <p className="text-[11px] font-bold text-gray-500 px-3 mb-2 uppercase tracking-wider">
                        Playlist Management
                    </p>
                    <CreateDialog />
                    <RenameDialog />
                    <DeleteDialog />
                    <PrefixDialog />
                </div>
            </div>

            <div className="flex-none pt-6 mt-4 border-t border-[#282828]">
                <div className="flex flex-col gap-3">
                    <PublishButton />
                    <PublishAllButton />
                </div>
            </div>
        </div>
    )
}

export function PublishButton() {
    const { publishPlaylist, publishLoading, currentPlaylistId: currentId } = usePlaylistStore();

    return (
        <Button 
            disabled={currentId === ""} 
            onClick={publishPlaylist}
            className="w-full rounded-full bg-[#1DB954] hover:bg-[#1ed760] text-black font-bold gap-2 py-6"
        >
            { publishLoading ? <Spinner className="text-black" /> : <ListMusic size={18} /> } 
            Publish Selection
        </Button>
    )
}

export function PublishAllButton() {
    const { publishPlaylists, publishAllLoading } = usePlaylistStore();

    return (
        <Button 
            variant="outline"
            onClick={publishPlaylists}
            className="w-full rounded-full border-gray-700 bg-transparent hover:bg-white/5 text-white font-bold gap-2 py-6"
        >
            { publishAllLoading ? <Spinner /> : <ListMusic size={18} /> } 
            Publish All
        </Button>
    )
}

export function DeleteDialog() {
    const deletePlaylist = usePlaylistStore((state) => state.deletePlaylist);
    
    const { currentPlaylistId: currentId, currentPlaylistName: current } = usePlaylistStore()

    const [playlistName, setPlaylistName] = React.useState("My Playlist");
    const inputId = React.useId(); 

    const handleDeletion = async () => {
        if (current != playlistName) return
        await deletePlaylist();
    };


    return (
        <Dialog>
        <DialogTrigger asChild>
                <Button disabled={currentId === ""} variant="ghost" className={cn(sidebarActionClass, "hover:text-destructive")}>
                    <Trash size={18} /> Delete Playlist
                </Button>
            </DialogTrigger>
<DialogContent className="max-w-[700px] sm:max-w-md bg-[#181818] border-[#282828] text-white shadow-2xl">
  <DialogHeader>
    <DialogTitle className="text-xl font-bold">Delete Playlist</DialogTitle>
    <DialogDescription className="text-gray-400 text-sm pt-2">
      This action cannot be undone. To confirm, please type{" "}
      <span className="text-white font-semibold">"{current}"</span> below.
    </DialogDescription>
  </DialogHeader>

  <div className="flex flex-col gap-4 py-4">
    <div className="grid gap-2">
      <Label htmlFor={inputId} className="text-xs font-bold uppercase tracking-wider text-gray-500">
        Confirm Name
      </Label>
      <Input
        id={inputId}
        value={playlistName}
        onChange={(e) => setPlaylistName(e.target.value)}
        placeholder={current}
        className="bg-[#242424] border-[#3e3e3e] text-white placeholder:text-gray-600 focus:ring-1 focus:ring-destructive focus:border-destructive transition-all h-11"
      />
    </div>
  </div>

  <DialogFooter className="sm:justify-end gap-2">
    <DialogClose asChild>
      <Button 
        type="button" 
        variant="ghost" 
        className="text-gray-400 hover:text-white hover:bg-white/5"
      >
        Cancel
      </Button>
    </DialogClose>
    
    <DialogClose asChild>
      <Button 
        type="button" 
        variant="destructive" 
        onClick={handleDeletion}
        disabled={current !== playlistName}
        className={cn(
          "font-bold px-8 transition-all",
          current === playlistName 
            ? "bg-red-600 hover:bg-red-500 opacity-100" 
            : "bg-red-900/20 text-red-900 opacity-50 cursor-not-allowed"
        )}
      >
        Delete Forever
      </Button>
    </DialogClose>
  </DialogFooter>
</DialogContent>
        </Dialog>
    )
}


export function RenameDialog() {
    const renamePlaylist = usePlaylistStore((state) => state.renamePlaylist);
    
    const { currentPlaylistName: current, currentPlaylistId: currentId } = usePlaylistStore()

    const [playlistName, setPlaylistName] = React.useState("My Playlist");
    const inputId = React.useId(); 

    const handleRename = async () => {
        if (!playlistName.trim()) return;
        await renamePlaylist(playlistName);
    };


    return (
        <Dialog>
<DialogTrigger asChild>
                <Button disabled={currentId === ""} variant="ghost" className={sidebarActionClass} onClick={() => setPlaylistName(current)}><PenLine size={18} /> Rename Playlist</Button>
            </DialogTrigger>
<DialogContent className="max-w-[700px] sm:max-w-md bg-[#181818] border-[#282828] text-white shadow-2xl">
  <DialogHeader>
    <DialogTitle className="text-xl font-bold">Rename Playlist</DialogTitle>
    <DialogDescription className="text-gray-400 text-sm pt-2">
      Enter a new name for your playlist. This will be updated across your library.
    </DialogDescription>
  </DialogHeader>

  <div className="flex flex-col gap-4 py-4">
    <div className="grid gap-2">
      <Label htmlFor={inputId} className="text-xs font-bold uppercase tracking-wider text-gray-500">
        New Name
      </Label>
      <Input
        id={inputId}
        value={playlistName}
        onChange={(e) => setPlaylistName(e.target.value)}
        placeholder={current}
        className="bg-[#242424] border-[#3e3e3e] text-white placeholder:text-gray-600 focus:ring-1 focus:ring-[#1DB954] focus:border-[#1DB954] transition-all h-11"
      />
    </div>
  </div>

  <DialogFooter className="sm:justify-end gap-2">
    <DialogClose asChild>
      <Button variant="ghost" className="text-gray-400 hover:text-white hover:bg-white/5">
        Cancel
      </Button>
    </DialogClose>
    
    <DialogClose asChild>
      <Button 
        type="button" 
        onClick={handleRename}
        disabled={!playlistName.trim() || playlistName === current}
        className="bg-[#1DB954] hover:bg-[#1ed760] text-black font-bold px-8 rounded-full transition-all disabled:opacity-50 disabled:bg-gray-700 disabled:text-gray-400"
      >
        Save Changes
      </Button>
    </DialogClose>
  </DialogFooter>
</DialogContent>
        </Dialog>
    )
}

export function CreateDialog() {
    const createPlaylist = usePlaylistStore((state) => state.createPlaylist);

    const [open, setOpen] = React.useState(false);
    const [step, setStep] = React.useState<1 | 2 | 3>(1);
    const [playlistName, setPlaylistName] = React.useState("My Playlist");
    const [newPlaylistId, setNewPlaylistId] = React.useState("");
    const [suggestions, setSuggestions] = React.useState<ModelItemResponse[]>([]);
    const [suggestionsLoading, setSuggestionsLoading] = React.useState(false);
    const [genreSuggestions, setGenreSuggestions] = React.useState<string[]>([]);
    const [genreSuggestionsLoading, setGenreSuggestionsLoading] = React.useState(false);
    const [includedArtist, setIncludedArtist] = React.useState("");
    const [customGenre, setCustomGenre] = React.useState("");
    const [creating, setCreating] = React.useState(false);
    const inputId = React.useId();

    const { playlistSelectionData, summary, setCurrentPlaylist, publishPlaylist } = usePlaylistStore();

    const resetDialog = () => {
        setStep(1);
        setPlaylistName("My Playlist");
        setSuggestions([]);
        setNewPlaylistId("");
        setGenreSuggestions([]);
    };

    const handleOpenChange = (isOpen: boolean) => {
        setOpen(isOpen);
        if (!isOpen) resetDialog();
    };

    const handleNext = async () => {
        if (!playlistName.trim()) return;
        setCreating(true);
        const newId = await createPlaylist(playlistName);
        setCreating(false);
        if (!newId) return;

        setNewPlaylistId(newId);
        setStep(2);
        setSuggestionsLoading(true);
        try {
            const res = await postSearch({ body: { playlistid: newId, query: playlistName, type: 1 } });
            setSuggestions(res.data ?? []);
        } finally {
            setSuggestionsLoading(false);
        }
    };

    const handleNextNext = async () => {
        setStep(3);
        setGenreSuggestionsLoading(true);
        try {
            const res = await getSpotifyArtistByIdGenres({ path: {id: includedArtist}})
            setGenreSuggestions(res.data ?? []);
        } finally {
            setGenreSuggestionsLoading(false);
        }
    };

    const handleFinishDialog = async () => {
        setCurrentPlaylist(newPlaylistId, playlistName)
        await publishPlaylist()
        await summary()
        resetDialog()
        setOpen(false)
    }

    const handleCustomGenre = async (name: string) => {
        setGenreSuggestions([...genreSuggestions, name])
    }

    const handleIncludeToggle = async (itemId: string, include: boolean, undo: boolean) => {
        const newIncluded = undo ? 0 : (include ? 1 : 2);
        setSuggestions(prev => prev.map(s => s.spotifyID === itemId ? { ...s, included: newIncluded } : s));
        try {
            if (undo) {
                await postPlaylistItemUndo({
                    body: {
                        include: include,
                        playlistid: newPlaylistId,
                        spotid: itemId,
                        type: 1,
                    }
                });
            } else {
                await postPlaylistItem({
                    body: {
                        include: include,
                        playlistid: newPlaylistId,
                        spotid: itemId,
                        type: 1,
                    }
                });
            }
        } catch {
            setSuggestions(prev => prev.map(s => s.spotifyID === itemId ? { ...s, included: undo ? newIncluded : 0 } : s));
        }
    };

    const handleIncludePlaylist = async (name: string, isSelected: boolean) => {
        var genrePlaylist = playlistSelectionData.find(i => i.name?.toLowerCase() == name.toLowerCase())?.spotifyID

        if (!genrePlaylist) {
            genrePlaylist = await createPlaylist(name)
            if (!genrePlaylist) return;
        }

        if (isSelected) {
            await postPlaylistInclude({
                body: {
                    pspotid: genrePlaylist,
                    cspotid: newPlaylistId,
                }
            })
        }
        else {
            await postPlaylistIncludeUndo({
                body: {
                    pspotid: genrePlaylist,
                    cspotid: newPlaylistId,
                }
            })
        }
    }

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogTrigger asChild>
                <Button variant="ghost" className={sidebarActionClass}>
                    <Plus size={18} className="text-[#1DB954]" /> Create Playlist
                </Button>
            </DialogTrigger>
            <DialogContent className="w-[50vw] sm:max-w-lg bg-[#181818] border-[#282828] text-white shadow-2xl">
                {step === 1 ? (
                    <>
                        <DialogHeader>
                            <DialogTitle className="text-xl font-bold flex items-center gap-2">
                                <Plus className="text-[#1DB954] h-5 w-5" /> Create New Playlist
                            </DialogTitle>
                            <DialogDescription className="text-gray-400 text-sm pt-2">
                                Give your new collection a name. You can add tracks to it immediately after.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="flex flex-col gap-4 py-4">
                            <div className="grid gap-2">
                                <Label htmlFor={inputId} className="text-xs font-bold uppercase tracking-wider text-gray-500">
                                    Playlist Name
                                </Label>
                                <Input
                                    id={inputId}
                                    value={playlistName}
                                    onChange={(e) => setPlaylistName(e.target.value)}
                                    onKeyDown={(e) => e.key === 'Enter' && handleNext()}
                                    placeholder="e.g. Late Night Vibes"
                                    className="bg-[#242424] border-[#3e3e3e] text-white placeholder:text-gray-600 focus:ring-1 focus:ring-[#1DB954] focus:border-[#1DB954] transition-all h-11"
                                />
                            </div>
                        </div>
                        <DialogFooter className="sm:justify-end gap-2">
                            <DialogClose asChild>
                                <Button variant="ghost" className="text-gray-400 hover:text-white hover:bg-white/5">
                                    Cancel
                                </Button>
                            </DialogClose>
                            <Button
                                type="button"
                                onClick={handleNext}
                                disabled={!playlistName.trim() || creating}
                                className="bg-white hover:bg-gray-200 text-black font-bold px-8 rounded-full transition-all disabled:bg-gray-700 disabled:text-gray-400"
                            >
                                {creating ? <Spinner className="text-black" /> : "Next →"}
                            </Button>
                        </DialogFooter>
                    </>
                ) : step === 2 ? (
                    <>
                        <DialogHeader>
                            <DialogTitle className="text-xl font-bold">Suggested Artists</DialogTitle>
                            <DialogDescription className="text-gray-400 text-sm pt-2">
                                Based on <span className="text-white font-medium">"{playlistName}"</span> — include artists to get you started.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="py-2">
                            {suggestionsLoading ? (
                                <div className="flex justify-center py-10">
                                    <Spinner />
                                </div>
                            ) : suggestions.length > 0 ? (
                                <div className="flex flex-col gap-1 max-h-72 overflow-y-auto custom-scrollbar pr-1">
                                    {suggestions.map(item => (
                                        <SearchResultItem
                                            key={item.spotifyID}
                                            item={item}
                                            onAction={(id, include, _type, undo) => {handleIncludeToggle(id, include, undo); setIncludedArtist(id)}}
                                            onExpand={() => {}}
                                        />
                                    ))}
                                </div>
                            ) : (
                                <p className="text-sm text-gray-400 text-center py-10">No suggestions found.</p>
                            )}
                        </div>

                        <DialogFooter className="sm:justify-end">
                            <DialogClose asChild>
                                <Button variant="ghost" className="text-gray-400 hover:text-white hover:bg-white/5">
                                    Cancel
                                </Button>
                            </DialogClose>
                            <Button
                                type="button"
                                onClick={handleNextNext}
                                disabled={!playlistName.trim() || creating}
                                className="bg-white hover:bg-gray-200 text-black font-bold px-8 rounded-full transition-all disabled:bg-gray-700 disabled:text-gray-400"
                            >
                                {creating ? <Spinner className="text-black" /> : "Next →"}
                            </Button>
                        </DialogFooter>
                    </>
                ) : (
                    <>
                        <DialogHeader>
                            <DialogTitle className="text-xl font-bold">Suggested Playlists</DialogTitle>
                        </DialogHeader>
                        <div className="py-2">
                            {genreSuggestionsLoading ? (
                                <div className="flex justify-center py-10">
                                    <Spinner />
                                </div>
                            ) : genreSuggestions.length > 0 ? (
                                <div className="flex flex-col gap-1 max-h-72 overflow-y-auto custom-scrollbar pr-1">
                                    {genreSuggestions.map(item => (
                                        <GenrePill
                                            key={item}
                                            name={item}
                                            onSelect={(item, isSelected) => {handleIncludePlaylist(item, isSelected)}}
                                        />
                                    ))}
                                </div>
                            ) : (
                                <p className="text-sm text-gray-400 text-center py-10">No suggestions found.</p>
                            )}
                        </div>
                        <div>
                                <Input
                                    id={inputId}
                                    value={customGenre}
                                    onChange={(e) => setCustomGenre(e.target.value)}
                                    onKeyDown={(e) => e.key === 'Enter' && handleCustomGenre(customGenre)}
                                    placeholder="e.g. Melodic Uptempo Doomdeath"
                                    className="bg-[#242424] border-[#3e3e3e] text-white placeholder:text-gray-600 focus:ring-1 focus:ring-[#1DB954] focus:border-[#1DB954] transition-all h-11"
                                />
                                </div>
                        <DialogFooter className="sm:justify-end">
                            <Button
                                type="button"
                                onClick={() => handleFinishDialog()}
                                className="bg-[#1DB954] hover:bg-[#1ed760] text-black font-bold px-8 rounded-full transition-all"
                            >
                                Done
                            </Button>
                        </DialogFooter>
                    </>
                )}
            </DialogContent>
        </Dialog>
    );
}

export function PlaylistSearch() {
    const [open, setOpen] = React.useState(false);
    const { playlistSelectionData, selectionLoading, currentPlaylistId, setCurrentPlaylist, summary, clearSearchData, clearSummaryData } = usePlaylistStore();

    const selectedPlaylist = playlistSelectionData.find((p) => p.spotifyID === currentPlaylistId);

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="ghost"
                    role="combobox"
                    className="w-full justify-between bg-white/5 hover:bg-white/10 text-white border-none h-12 px-4 rounded-md"
                    disabled={selectionLoading}
                >
                    <span className="truncate">
                        {selectedPlaylist ? selectedPlaylist.name : "Select a playlist..."}
                    </span>
                    <ChevronsUpDownIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[280px] p-0 bg-[#282828] border-[#3e3e3e] shadow-2xl">
                <Command className="bg-transparent">
                    <CommandInput placeholder="Search your library..." className="text-white" />
                    <CommandList className="max-h-[300px]">
                        <CommandEmpty className="py-6 text-center text-sm text-gray-400">No playlist found.</CommandEmpty>
                        <CommandGroup>
                            {playlistSelectionData.map((p) => (
                                <CommandItem
                                    key={p.spotifyID}
                                    onSelect={() => {
                                        const isSelected = currentPlaylistId === p.spotifyID;
                                        setCurrentPlaylist(isSelected ? "" : (p.spotifyID as string), (p.name as string));
                                        clearSummaryData();
                                        clearSearchData();
                                        summary();
                                        setOpen(false);
                                    }}
                                    className="aria-selected:bg-white/10 text-gray-300 hover:text-white"
                                >
                                    <CheckIcon className={cn("mr-2 h-4 w-4 text-[#1DB954]", currentPlaylistId === p.spotifyID ? "opacity-100" : "opacity-0")} />
                                    {p.name}
                                </CommandItem>
                            ))}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    )
}

export function PrefixDialog() {
    const { setPlaylistPrefix, getPlaylistPrefix } = usePlaylistStore()

    const [currentPrefix, setCurrentPrefix] = React.useState(getPlaylistPrefix);
    const inputId = React.useId(); 

    const handleUpdate = async () => {
        setPlaylistPrefix(currentPrefix);
    };


    return (
        <Dialog>
        <DialogTrigger asChild>
                <Button variant="ghost" className={cn(sidebarActionClass)}>
                    <Trash size={18} /> Change Prefix
                </Button>
            </DialogTrigger>
                <DialogContent className="max-w-[700px] sm:max-w-md bg-[#181818] border-[#282828] text-white shadow-2xl">
                  <DialogHeader>
                    <DialogTitle className="text-xl font-bold">Change Prefix</DialogTitle>
                    <DialogDescription className="text-gray-400 text-sm pt-2">
                       Change the prefix that is prepended to each playlist. This will not retroactively update all playlists.
                    </DialogDescription>
                  </DialogHeader>

                  <div className="flex flex-col gap-4 py-4">
                    <div className="grid gap-2">
                      <Label htmlFor={inputId} className="text-xs font-bold uppercase tracking-wider text-gray-500">
                        New Prefix
                      </Label>
                      <Input
                        id={inputId}
                        value={currentPrefix}
                        onChange={(e) => setCurrentPrefix(e.target.value)}
                        placeholder={currentPrefix}
                        className="bg-[#242424] border-[#3e3e3e] text-white placeholder:text-gray-600 focus:ring-1 focus:ring-destructive focus:border-destructive transition-all h-11"
                      />
                    </div>
                  </div>

                  <DialogFooter className="sm:justify-end gap-2">
                    <DialogClose asChild>
                      <Button 
                        type="button" 
                        variant="ghost" 
                        className="text-gray-400 hover:text-white hover:bg-white/5"
                      >
                        Cancel
                      </Button>
                    </DialogClose>
                    
                    <DialogClose asChild>
                      <Button 
                        type="button" 
                        onClick={handleUpdate}
                        className={cn(
                          "font-bold px-8 transition-all",
                        )}
                      >
                        Confirm Prefix
                      </Button>
                    </DialogClose>
                  </DialogFooter>
                </DialogContent>
        </Dialog>
    )
}
