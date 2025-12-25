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

import { usePlaylistStore } from "@/components/stores/playlist.store"

export default function Playlist() {
    const fetchPlaylists = usePlaylistStore((state) => state.fetch);

    // Fetch data when the component mounts
    React.useEffect(() => {
        fetchPlaylists();
    }, [fetchPlaylists]);

    return (
        // <div className="flex items-center gap-2">
        <>
        <PlaylistSearch />
        <PlaylistCreate />
        </>
        // </div>
    )
}

function PlaylistCreate() {
    return <Button variant={'green'}> <Plus /> Create Playlist</Button>
}

function PlaylistSearch() {
    const [open, setOpen] = React.useState(false)
    const [value, setValue] = React.useState("")
    
    // 1. Destructure data and loading from the store
    const { data, loading } = usePlaylistStore();

    // 2. Prevent crash: If data is undefined or null, fallback to an empty array
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
                    {/* 3. Safe Check: Use optional chaining (?.) to prevent crashes */}
                    {value
                        ? safeData.find((p) => p.name === value)?.name || "Select playlist..."
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
                            {/* 4. Use the safeData array */}
                            {safeData.map((p) => (
                                <CommandItem
                                    key={p.id || p.name}
                                    value={p.name}
                                    onSelect={(currentValue) => {
                                        setValue(currentValue === value ? "" : currentValue)
                                        setOpen(false)
                                    }}
                                >
                                    <CheckIcon
                                        className={cn(
                                            "mr-2 h-4 w-4",
                                            value === p.name ? "opacity-100" : "opacity-0"
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
