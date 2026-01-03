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
import { useSearchStore } from "@/components/stores/search.store"

import { Spinner } from "@/components/ui/spinner"
import { useSummaryStore } from "../stores/summary.store"

export default function Playlist() {
    const fetchPlaylists = usePlaylistStore((state) => state.fetch);

    React.useEffect(() => {
        fetchPlaylists();
    }, [fetchPlaylists]);

    return (
        <div className="flex items-center gap-2">
        <PlaylistSearch />
        <CreateDialog />
        <PublishButton />
        <RenameDialog />
        <DeleteDialog />
        </div>
    )
}

export function PublishButton() {
    const publishPlaylist = usePlaylistStore((state) => state.publishPlaylist);

    const { publishLoading, currentId } = usePlaylistStore()

    const handlePublish = async () => {
        await publishPlaylist();
    };

    return (
        <Button disabled={currentId == ""} variant="green" onClick={handlePublish} >{ publishLoading ? <Spinner /> : <ListMusic /> } Publish Playlist</Button>
    )
}

export function DeleteDialog() {
    const deletePlaylist = usePlaylistStore((state) => state.deletePlaylist);
    
    const { currentId, current } = usePlaylistStore()

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
    
    const { current } = usePlaylistStore()

    const [playlistName, setPlaylistName] = React.useState(current);
    const inputId = React.useId(); 

    const handleRename = async () => {
        if (!playlistName.trim()) return;
        await renamePlaylist(playlistName);
        setPlaylistName("My Playlist");
    };


    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="green"><PenLine />Rename Playlist</Button>
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
    
    const { data, loading, current, setCurrentId } = usePlaylistStore();
    const { setPlaylistId, clearData } = useSearchStore();
    const { setPlaylistIdSummary } = useSummaryStore();

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="w-[200px] justify-between"
                    disabled={loading} // Disable button while loading
                >
                    { current
                        ? data.find((p) => p.name === current)?.name || "Select playlist..."
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
                                {data.map((p) => (
                                    <CommandItem
                                        key={p.spotifyID} 
                                        value={p.spotifyID}
                                        onSelect={() => {
                                            const isSelected = current === p.name;
                                            setCurrentId(isSelected ? "" : (p.spotifyID as string), (p.name as string));
                                            setPlaylistId(isSelected ? "" : (p.spotifyID as string));
                                            setPlaylistIdSummary(isSelected ? "" : (p.spotifyID as string));
                                            clearData();
                                            setOpen(false);
                                        }}
                                    >
                                        <CheckIcon
                                            className={cn(
                                                "mr-2 h-4 w-4",
                                                current === p.name ? "opacity-100" : "opacity-0"
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


