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
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, Plus, X } from "lucide-react";

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


interface SearchResultItemProps {
  item: ModelItemResponse;
  onAction: (id: string, include: boolean) => void;
}

export function SearchResultItem({ item, onAction }: SearchResultItemProps) {
  const imageUrl = item.icon && item.icon.length > 0 ? item.icon[0].url : "";
  
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1: 
        return <Badge className="bg-green-500 hover:bg-green-600">Included</Badge>;
      case 2: 
        return <Badge variant="destructive">Excluded</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className="flex items-center justify-between p-3 transition-colors hover:bg-muted/50 rounded-lg border">
      <div className="flex items-center gap-4">
        <Avatar className="h-12 w-12 rounded-md">
          <AvatarImage src={imageUrl} alt={item.name} />
          <AvatarFallback className="rounded-md">
            {item.name?.substring(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <span className="font-medium text-sm leading-none">{item.name}</span>
            {getInclusionBadge()}
          </div>
          <div className="flex items-start gap-2">
              <span className="text-xs text-muted-foreground capitalize">
                {["Playlist", "Artist", "Album", "Track"][item.itemType ?? 3]}
              </span>
          </div>
        </div>
      </div>

      <div className="flex gap-2">
        {/* Toggle Inclusion Button */}
        <Button 
          size="icon" 
          variant={item.included === 1 ? "default" : "outline"} 
          className="h-8 w-8"
          onClick={() => item.spotifyID && onAction(item.spotifyID, true)}
        >
          {item.included === 1 ? <Check className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
        </Button>

        {/* Toggle Exclusion Button */}
        <Button 
          size="icon" 
          variant={item.included === 2 ? "destructive" : "outline"} 
          className="h-8 w-8"
          onClick={() => item.spotifyID && onAction(item.spotifyID, false)}
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

function ResultBox() {
  const { data, loading } = useSearchStore();
  // You'll need an action in a store to handle the actual DB update
  const handleInclusion = (id: string, include: boolean) => { 
    console.log(`Setting item ${id} to include: ${include}`);
    // call your includeExcludeItem service here
    // TODO
  };

  return (
    <ResizablePanelGroup direction="horizontal" className="min-h-[400px] rounded-lg border">
      <ResizablePanel defaultSize={25}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[600px]">
          {loading ? (
            <p>Loading results...</p>
          ) : data.length > 0 ? (
            data.map((item) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No results found.</p>
          )}
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel defaultSize={75}>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

function ResultElement({ el }) {
    const icon = el.icon.length > 0 ? el.icon[0].url : null
    return (
        <>
        <img src={icon} alt={el.name} width="100" height="100"/>
        </>
    )
}

function SearchBar() {
    const search = useSearchStore((state) => state.search);

    // Assuming you have a way to get the actual Playlist Spotify ID
    const [query, setQuery] = React.useState("");
    const [type, setItemType] = React.useState<ModelItemType>(1);

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
                                <DropdownMenuRadioItem value="1">Artist</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="2">Album</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="3">Track</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="0">Playlist</DropdownMenuRadioItem>
                            </DropdownMenuRadioGroup>
                        </DropdownMenuContent>
                    </DropdownMenu>
                    <InputGroupButton onClick={handleSearch}>Go</InputGroupButton>
                </InputGroupAddon>
            </InputGroup>
        </div>
    )
}
