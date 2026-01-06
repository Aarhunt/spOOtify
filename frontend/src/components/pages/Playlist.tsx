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

import { Spinner } from "@/components/ui/spinner"

export default function Playlist() {
    const fetchPlaylists = usePlaylistStore((state) => state.fetchSelectionData);

    React.useEffect(() => {
        fetchPlaylists();
    }, [fetchPlaylists]);

    return (
        <div className="flex items-center gap-2">
        <PlaylistSearch />
        <CreateDialog />
        <PublishButton />
        <PublishAllButton />
        <RenameDialog />
        <DeleteDialog />
        </div>
    )
}

export function PublishButton() {
    const publishPlaylist = usePlaylistStore((state) => state.publishPlaylist);

    const { publishLoading, currentPlaylistId: currentId } = usePlaylistStore()

    const handlePublish = async () => {
        await publishPlaylist();
    };

    return (
        <Button disabled={currentId == ""} variant="green" onClick={handlePublish} >{ publishLoading ? <Spinner /> : <ListMusic /> } Publish Playlist</Button>
    )
}

export function PublishAllButton() {
    const publishPlaylists = usePlaylistStore((state) => state.publishPlaylists);

    const { publishAllLoading } = usePlaylistStore()

    const handlePublish = async () => {
        await publishPlaylists();
    };

    return (
        <Button variant="green" onClick={handlePublish} >{ publishAllLoading ? <Spinner /> : <ListMusic /> } Publish All Playlists</Button>
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
                <Button disabled={currentId == ""} variant="destructive"><Trash />Delete Playlist</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Delete Playlist</DialogTitle>
                    <DialogDescription>
                        Are you sure you want to delete this playlist?
                        Type the name of the playlist to confirm deletion.
                    </DialogDescription>
                </DialogHeader>
                <div className="flex items-center gap-2">
                    <div className="grid flex-1 gap-2">
                        <Label htmlFor={inputId} className="sr-only">
                            Name
                        </Label>
                        <Input
                            id={inputId}
                            value={playlistName}
                            onChange={(e) => setPlaylistName(e.target.value)}
                            placeholder="My Playlist"
                        />
                    </div>
                </div>
                <DialogFooter className="sm:justify-start">
                    <DialogClose asChild>
                        <Button 
                            type="button" 
                            variant="destructive" 
                            onClick={handleDeletion}
                            disabled={current != playlistName}
                        >
                            Delete
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
                <Button disabled={currentId == ""} variant="green" onClick={() => setPlaylistName(current)}><PenLine />Rename Playlist</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Rename Playlist</DialogTitle>
                    <DialogDescription>
                        Fill in the new name of the playlist here.
                    </DialogDescription>
                </DialogHeader>
                <div className="flex items-center gap-2">
                    <div className="grid flex-1 gap-2">
                        <Label htmlFor={inputId} className="sr-only">
                            Name
                        </Label>
                        <Input
                            id={inputId}
                            value={playlistName}
                            onChange={(e) => setPlaylistName(e.target.value)}
                            placeholder="My Playlist"
                        />
                    </div>
                </div>
                <DialogFooter className="sm:justify-start">
                    <DialogClose asChild>
                        <Button 
                            type="button" 
                            variant="outline" 
                            onClick={handleRename}
                        >
                            Rename
                        </Button>
                    </DialogClose>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

export function CreateDialog() {
    const createPlaylist = usePlaylistStore((state) => state.createPlaylist);
    
    const [playlistName, setPlaylistName] = React.useState("My Playlist");
    const inputId = React.useId(); 

    const handleCreate = async () => {
        if (!playlistName.trim()) return;
        await createPlaylist(playlistName);
        // Optional: Reset field after creation
        setPlaylistName("My Playlist");
    };


    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="green"><Plus />Create Playlist</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Create Playlist</DialogTitle>
                    <DialogDescription>
                        Fill in the name of the playlist here.
                    </DialogDescription>
                </DialogHeader>
                <div className="flex items-center gap-2">
                    <div className="grid flex-1 gap-2">
                        <Label htmlFor={inputId} className="sr-only">
                            Name
                        </Label>
                        <Input
                            id={inputId}
                            value={playlistName}
                            onChange={(e) => setPlaylistName(e.target.value)}
                            placeholder="My Playlist"
                        />
                    </div>
                </div>
                <DialogFooter className="sm:justify-start">
                    <DialogClose asChild>
                        <Button 
                            type="button" 
                            variant="outline" 
                            onClick={handleCreate}
                        >
                            Create
                        </Button>
                    </DialogClose>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

export function PlaylistSearch() {
    const [open, setOpen] = React.useState(false)
    
    const { playlistSelectionData, selectionLoading, currentPlaylistName, currentPlaylistId, setCurrentPlaylist, summary, clearSearchData, clearSummaryData } = usePlaylistStore();

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="w-[200px] justify-between"
                    disabled={selectionLoading} // Disable button while loading
                >
                    { currentPlaylistName
                        ? playlistSelectionData.find((p) => p.spotifyID === currentPlaylistId)?.name || "Select playlist..."
                        : "Select playlist..."}
                    <ChevronsUpDownIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[200px] p-0">
                <Command>
                    <CommandInput placeholder="Search playlists..." />
                    <CommandList>
                        <CommandEmpty>No playlist found.</CommandEmpty>
                            <CommandGroup>
                                {playlistSelectionData.map((p) => (
                                    <CommandItem
                                        key={p.spotifyID} 
                                        value={p.spotifyID}
                                        onSelect={() => {
                                            const isSelected = currentPlaylistId === p.spotifyID;
                                            setCurrentPlaylist(isSelected ? "" : (p.spotifyID as string), (p.name as string));
                                            clearSummaryData();
                                            clearSearchData();
                                            summary();
                                            setOpen(false);
                                        }}
                                    >
                                        <CheckIcon
                                            className={cn(
                                                "mr-2 h-4 w-4",
                                                currentPlaylistId === p.spotifyID ? "opacity-100" : "opacity-0"
                                            )}
                                        />
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


