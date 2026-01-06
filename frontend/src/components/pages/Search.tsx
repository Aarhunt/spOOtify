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

import React from "react";
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, Plus, X, Music } from "lucide-react";
import { cn } from "@/lib/utils";

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import { usePlaylistStore } from "../stores/playlist.store";

export default function Search() {
    return (
        <>
        <SearchBar />
        <ResultBox />
        </>
    )
}

interface PlaylistResultItemProps {
  item: ModelItemResponse;
  onAction: (id: string, include: boolean, type: ModelItemType, undo: boolean) => void;
  onExpand: (id: string, type: ModelItemType) => void;
}

export function PlaylistResultItem({ item, onAction, onExpand }: PlaylistResultItemProps) {
  const imageUrl = item.icon && item.icon.length > 0 ? item.icon[0].url : "";
  
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1: 
        return <Badge className="bg-green-500 hover:bg-green-435">Included</Badge>;
      case 0: 
        return null
      default:
        return <Badge className="bg-blue-500 hover:bg-blue-435">Included By Proxy</Badge>;
    }
  };

  return (
    <div className="flex items-center justify-between p-3 transition-colors hover:bg-muted/50 rounded-lg border" onClick={() => item.spotifyID && item.itemType && onExpand(item.spotifyID, item.itemType)}>
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
                {item.itemType == 2 ? item.sortdata : ""}
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
          onClick={(e) => {
            e.stopPropagation();
              item.spotifyID && onAction(item.spotifyID, true, item.itemType!, item.included === 1);}}
        >
          {item.included === 1 ? <Check className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
        </Button>
      </div>
    </div>
  );
}

interface SearchResultItemProps {
  item: ModelItemResponse;
  onAction: (id: string, include: boolean, type: ModelItemType, undo: boolean) => void;
  onExpand: (id: string, type: ModelItemType) => void;
}

export function SearchResultItem({ item, onAction, onExpand }: SearchResultItemProps) {
  const imageUrl = item.icon && item.icon.length > 0 ? item.icon[0].url : "";
  
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1: 
        return <Badge className="bg-green-500 hover:bg-green-435">Included</Badge>;
      case 2: 
        return <Badge variant="destructive">Excluded</Badge>;
      case 3: 
        return <Badge className="bg-blue-500 hover:bg-blue-435">Included</Badge>;
      case 4: 
        return <Badge className="bg-orange-500 hover:bg-orange-435">Excluded</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className="flex items-center justify-between p-3 transition-colors hover:bg-muted/50 rounded-lg border" onClick={() => item.spotifyID && item.itemType && onExpand(item.spotifyID, item.itemType)}>
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
                {item.itemType == 2 ? item.sortdata : ""}
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
          onClick={(e) => {
            e.stopPropagation();
              item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);}}
        >
          {item.included === 1 ? <Check className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
        </Button>

        {/* Toggle Exclusion Button */}
        <Button 
          size="icon" 
          variant={item.included === 2 ? "destructive" : "outline"} 
          className="h-8 w-8"
          onClick={(e) => {
                e.stopPropagation();
              item.spotifyID && item.itemType && onAction(item.spotifyID, false, item.itemType, item.included === 2);}}
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

interface TrackResultItemProps {
  item: ModelItemResponse;
  onAction: (id: string, include: boolean, type: ModelItemType, undo: boolean) => void;
}

export function TrackResultItem({ item, onAction }: TrackResultItemProps) {
  
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1:
        return <Badge className="bg-green-500/10 text-green-500 hover:bg-green-500/20 border-none h-5 px-1.5 text-[10px]">IN</Badge>;
      case 2:
        return <Badge className="bg-destructive/10 text-destructive hover:bg-destructive/20 border-none h-5 px-1.5 text-[10px]">EX</Badge>;
      case 3:
        return <Badge className="bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 border-none h-5 px-1.5 text-[10px]">IN</Badge>;
      case 4:
        return <Badge className="bg-orange-500-500/10 text-orange-500 hover:bg-orange-500/20 border-none h-5 px-1.5 text-[10px]">EX</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className="group flex items-center justify-between py-1.5 px-3 transition-colors hover:bg-muted/40 rounded-md border border-transparent hover:border-border">
      <div className="flex items-center gap-3 overflow-hidden">
        <div className="flex items-center justify-center w-6 text-muted-foreground text-xs font-medium">
          {item.sortdata ? item.sortdata : <Music className="h-3 w-3 opacity-40" />}
        </div>
        
        <div className="flex flex-col min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-medium text-sm truncate">{item.name}</span>
            {getInclusionBadge()}
          </div>
        </div>
      </div>

      <div className="flex gap-1 opacity-40 group-hover:opacity-100 transition-opacity">
        <Button 
          size="icon" 
          variant="ghost"
          className={cn(
            "h-7 w-7",
            item.included === 1 && "text-green-500 bg-green-500/10 hover:bg-green-500/20 hover:text-green-435"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);
          }}
        >
          {item.included === 1 ? <Check className="h-3.5 w-3.5" /> : <Plus className="h-3.5 w-3.5" />}
        </Button>

        <Button 
          size="icon" 
          variant="ghost"
          className={cn(
            "h-7 w-7",
            item.included === 2 && "text-destructive bg-destructive/10 hover:bg-destructive/20"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, false, item.itemType, item.included === 2);
          }}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>
    </div>
  );
}

function ResultBox() {
  const { playlistSearchData, artistSearchData, searchLoading, albumSearchData, trackSearchData, albumSearchLoading, trackSearchLoading, searchType, includeItem, undoIncludeItem, getSearchAlbumsFromArtist, getSearchTracksFromAlbum, setCurrentSearchArtist, setCurrentSearchAlbum, includePlaylist, undoIncludePlaylist } = usePlaylistStore();
  const handleInclusion = (id: string, include: boolean, type: ModelItemType, undo: boolean) => { 
    type == 0 ? 
        undo ? 
            undoIncludePlaylist(id) : 
            includePlaylist(id) 
    :
        undo ? 
            undoIncludeItem(id, include, type) : 
            includeItem(id, include, type)
  };

  const handleExpand = (id: string, type: ModelItemType) => {
    switch (type) {
        case 1:
            getSearchAlbumsFromArtist(id) 
            setCurrentSearchArtist(id) 
            break;
        case 2:
            getSearchTracksFromAlbum(id)
            setCurrentSearchAlbum(id)
            break;
        default:
            break;
    } 
  }

  return (
    <ResizablePanelGroup direction="horizontal" className="min-h-[435px] rounded-lg border">
    {
        searchType == 0 ? <>
          <ResizablePanel defaultSize={playlistSearchData.length > 0 ? 50 : 50}>
            <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[435px]">
              {searchLoading ? (
                <p>Loading playlists...</p>
              ) : playlistSearchData.length > 0 ? (
                playlistSearchData.map((item) => (
                  <PlaylistResultItem 
                    key={item.spotifyID} 
                    item={item} 
                    onAction={handleInclusion} 
                    onExpand={() => {}}
                  />
                ))
              ) : (
                <p className="text-muted-foreground text-center py-10">No playlists found.</p>
              )}
            </div>
          </ResizablePanel>
        </> : null
    }
    {
        searchType == 1 ? 
        <>
      <ResizablePanel defaultSize={albumSearchData.length > 0 ? 50 : 50}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[435px]">
          {searchLoading ? (
            <p>Loading artists...</p>
          ) : artistSearchData.length > 0 ? (
            artistSearchData.map((item) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={handleExpand}
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No artists found.</p>
          )}
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      </> : null
    }
      {
        searchType == 1 || searchType == 2 ?  
        <> 
          <ResizablePanel defaultSize={albumSearchData.length > 0 ? 50 : 25}>
            <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[435px]">
              {albumSearchLoading ? (
                <p>Loading albums...</p>
              ) : albumSearchData.length > 0 ? (
                albumSearchData.map((item) => (
                  <SearchResultItem 
                    key={item.spotifyID} 
                    item={item} 
                    onAction={handleInclusion} 
                    onExpand={handleExpand}
                  />
                ))
              ) : (
                <p className="text-muted-foreground text-center py-10">No albums found.</p>
              )}
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle />
        </> : null
      }
      {
        searchType == 1 || searchType == 2 ? 
        <>
      <ResizablePanel defaultSize={trackSearchData.length > 0 ? 50 : 25}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[435px]">
          {trackSearchLoading ? (
            <p>Loading tracks...</p>
          ) : trackSearchData.length > 0 ? (
            trackSearchData.map((item) => (
            <TrackResultItem
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No tracks found.</p>
          )}
        </div>
      </ResizablePanel>
    </> : null
      }
      {
        searchType == 3 ? 
        <>
      <ResizablePanel defaultSize={trackSearchData.length > 0 ? 50 : 25}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[435px]">
          {trackSearchLoading ? (
            <p>Loading tracks...</p>
          ) : trackSearchData.length > 0 ? (
            trackSearchData.map((item) => (
            <SearchResultItem
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => void {}}
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No tracks found.</p>
          )}
        </div>
      </ResizablePanel>
    </> : null
      }
    </ResizablePanelGroup>
  );
}

function SearchBar() {
    const { searchType, setSearchType, search } = usePlaylistStore();

    const [query, setQuery] = React.useState("");

    const handleSearch = async () => {
        if (searchType != 0 && !query) return;
        await search(query)
    }

    const handleTypeChange = (value: string) => {
        setSearchType(parseInt(value) as ModelItemType);
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
                                {["Playlist", "Artist", "Album", "Track"][searchType]}
                            </InputGroupButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56">
                            <DropdownMenuLabel>Search Type</DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuRadioGroup value={searchType.toString()} onValueChange={handleTypeChange}>
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
