"use client"

import {
    InputGroup,
    InputGroupAddon,
    InputGroupButton,
    InputGroupInput,
} from "@/components/ui/input-group"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { useSearchStore } from "@/components/stores/search.store"
import React from "react";
import type { ModelItemType } from "@/client/types.gen";

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"

export default function Search() {
    return (
        <>
        <SearchBar />
        <ResultBox />
        </>
    )
}

function ResultBox() {
  return (
    <ResizablePanelGroup
      direction="horizontal"
      className="min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"
    >
      <ResizablePanel defaultSize={25}>
        <div className="flex h-full items-center justify-center p-6">
          <span className="font-semibold">Sidebar</span>
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel defaultSize={75}>
        <div className="flex h-full items-center justify-center p-6">
          <span className="font-semibold">Content</span>
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  )
}

function SearchBar() {
    const search = useSearchStore((state) => state.search);

    // Assuming you have a way to get the actual Playlist Spotify ID
    const [query, setQuery] = React.useState("");
    const [type, setItemType] = React.useState<ModelItemType>(0);

    const handleSearch = async () => {
        if (!query) return;
        await search(query, type)
    }

    // Helper to handle the string-to-number conversion from the dropdown
    const handleTypeChange = (value: string) => {
        setItemType(parseInt(value) as ModelItemType);
    };

    return (
        <div className="grid w-full max-w-sm gap-6">
            <InputGroup className="[--radius:1rem]">
                <InputGroupInput 
                    placeholder="Search Spotify..." 
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                />
                <InputGroupAddon align="inline-end">
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <InputGroupButton variant="ghost" className="!pr-1.5 text-xs">
                                {["Playlist", "Artist", "Album", "Track"][type]}
                            </InputGroupButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56">
                            <DropdownMenuLabel>Search Type</DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuRadioGroup value={type.toString()} onValueChange={handleTypeChange}>
                                <DropdownMenuRadioItem value="0">Playlist</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="1">Artist</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="2">Album</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="3">Track</DropdownMenuRadioItem>
                            </DropdownMenuRadioGroup>
                        </DropdownMenuContent>
                    </DropdownMenu>
                    <InputGroupButton onClick={handleSearch}>Go</InputGroupButton>
                </InputGroupAddon>
            </InputGroup>
        </div>
    )
}
