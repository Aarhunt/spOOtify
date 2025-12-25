"use client"

import * as React from "react"
import { CheckIcon, ChevronsUpDownIcon, Plus } from "lucide-react"

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

export default function Playlist() {
    const fetchPlaylists = usePlaylistStore((state) => state.fetch);

    React.useEffect(() => {
        fetchPlaylists();
    }, [fetchPlaylists]);

    return (
        <div className="flex items-center gap-2">
        <PlaylistSearch />
        <DialogCloseButton />
        </div>
    )
}

export function DialogCloseButton() {
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

function PlaylistSearch() {
    const [open, setOpen] = React.useState(false)
    
    const { data, loading, current, setCurrent } = usePlaylistStore();

    const safeData = data || [];

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
                        ? safeData.find((p) => p.name === current)?.name || "Select playlist..."
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
                            {safeData.map((p) => (
                                <CommandItem
                                    key={p.id || p.name}
                                    value={p.name}
                                    onSelect={(currentValue) => {
                                        setCurrent(currentValue === current ? "" : currentValue)
                                        setOpen(false)
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


