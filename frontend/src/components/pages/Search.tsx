"use client"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import React, { useEffect } from "react";
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, Plus, X, Music, SearchIcon } from "lucide-react";
import { cn } from "@/lib/utils";

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import { usePlaylistStore } from "../stores/playlist.store";

export default function Search() {
  return (
    <div className="flex flex-col w-full mx-auto items-center">
      <div className="sticky top-0 z-10 w-full flex flex-col items-center justify-center backdrop-blur-md pb-4 px-4">
        <div className="w-full max-w-3xl">
          <SearchBar />
        </div>
      </div>

      <div className="w-full px-4 mt-0">
        <ResultBox />
      </div>
    </div>
  );
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
              <span className="text-xs text-muted-foreground capitalize text-left">
                {item.itemType == 2 ? item.sortdata : ""}
              </span>
          </div>
        </div>
      </div>

      <div className="flex gap-2">
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
        return <Badge className="bg-[#1DB954]/10 text-[#1DB954] border-none text-[10px] h-5 px-1.5">Included</Badge>;
      case 2: 
        return <Badge className="bg-destructive/10 text-destructive border-none text-[10px] h-5 px-1.5">Excluded</Badge>;
      case 3: 
        return <Badge className="bg-[blue]/10 text-[blue] border-none text-[10px] h-5 px-1.5">Included</Badge>;
      case 4: 
        return <Badge className="bg-[orange]/10 text-[orange] border-none text-[10px] h-5 px-1.5">Excluded</Badge>;
      default:
        return null;
    }
  };

    return (
    <div 
      className="group flex items-center justify-between p-2 transition-all hover:bg-[#ffffff1a] rounded-md cursor-pointer" 
      onClick={() => item.spotifyID && item.itemType && onExpand(item.spotifyID, item.itemType)}
    >
      <div className="flex items-center gap-4 overflow-hidden">
        <Avatar className="h-12 w-12 rounded shadow-lg">
          <AvatarImage src={imageUrl} alt={item.name} className="object-cover" />
          <AvatarFallback className="bg-[#282828] text-xs">
            {item.name?.substring(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        
        <div className="flex flex-col min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-semibold text-sm truncate text-white group-hover:text-[#1DB954] transition-colors">
              {item.name}
            </span>
            {getInclusionBadge()}
          </div>
          <span className="text-xs text-gray-400 capitalize text-left">
            {item.itemType === 2 ? item.sortdata : ""}
          </span>
        </div>
      </div>

      <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
        <Button 
          size="icon" 
          variant="ghost" 
          className={cn(
            "h-8 w-8 rounded-full transition-colors",
            item.included === 1 ? "bg-[#1DB954] text-black hover:bg-[#1ed760]" : "bg-[#282828] text-gray-400 hover:text-white"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);
          }}
        >
          {item.included === 1 ? <Check className="h-4 w-4 stroke-[3]" /> : <Plus className="h-4 w-4" />}
        </Button>

        <Button 
          size="icon" 
          variant="ghost" 
          className={cn(
            "h-8 w-8 rounded-full transition-colors",
            item.included === 2 ? "bg-destructive text-white" : "bg-[#282828] text-gray-400 hover:text-destructive"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, false, item.itemType, item.included === 2);
          }}
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

function TrackResultItem({ item, onAction }: TrackResultItemProps) {
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1:
        return <Badge className="bg-[#1DB954]/10 text-[#1DB954] border-none text-[10px] h-4 px-1.5">IN</Badge>;
      case 2:
        return <Badge className="bg-destructive/10 text-destructive border-none text-[10px] h-4 px-1.5">EX</Badge>;
      case 3:
        return <Badge className="bg-blue-400/10 text-blue-400 border-none text-[10px] h-4 px-1.5">IN</Badge>;
      case 4:
        return <Badge className="bg-orange-400/10 text-orange-400 border-none text-[10px] h-4 px-1.5">EX</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className="group flex items-center justify-between py-2 px-3 transition-all hover:bg-[#ffffff1a] rounded-md cursor-default">
      <div className="flex items-center gap-3 overflow-hidden">
        <div className="flex items-center justify-start w-6 text-muted-foreground text-xs font-medium tabular-nums">
          {item.sortdata ? (
            <span className="group-hover:hidden">{item.sortdata}</span>
          ) : (
            <Music className="h-3.5 w-3.5 opacity-40 group-hover:hidden" />
          )}
          <Music className="h-3.5 w-3.5 hidden group-hover:block text-white" />
        </div>
        
        <div className="flex flex-col min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-medium text-sm truncate text-gray-200 group-hover:text-white transition-colors">
              {item.name}
            </span>
            {getInclusionBadge()}
          </div>
        </div>
      </div>

      <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-all">
        <Button 
          size="icon" 
          variant="ghost"
          className={cn(
            "h-7 w-7 rounded-full",
            item.included === 1 ? "bg-[#1DB954] text-black hover:bg-[#1ed760]" : "hover:bg-[#333] text-gray-400"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);
          }}
        >
          {item.included === 1 ? <Check className="h-3.5 w-3.5 stroke-[3]" /> : <Plus className="h-3.5 w-3.5" />}
        </Button>

        <Button 
          size="icon" 
          variant="ghost"
          className={cn(
            "h-7 w-7 rounded-full",
            item.included === 2 ? "bg-destructive text-white" : "hover:bg-[#333] text-gray-400 hover:text-destructive"
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

  const renderEmpty = (label: string) => (
    <div className="flex flex-col items-center justify-center py-20 opacity-20">
      <Music className="h-12 w-12 mb-4" />
      <p className="text-sm font-medium">No {label} Found</p>
    </div>
  ); 

  return (
<ResizablePanelGroup className="flex-1 rounded-xl bg-[#121212] border border-[#282828] overflow-hidden max-h-[405px] ">
    {(searchType === 0) && (
        <>
          <ResizablePanel defaultSize={50} minSize={20} className="flex flex-col">
            <div className="px-4 py-3 bg-[#181818] border-b border-[#282828]">
              <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Playlists</span>
            </div>
            <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
              {searchLoading ? <p className="text-xs p-4 animate-pulse">Searching...</p> : 
                playlistSearchData.length > 0 ? playlistSearchData.map(item => (
                  <SearchResultItem key={item.spotifyID} item={item} onAction={handleInclusion} onExpand={handleExpand} />
                )) : renderEmpty("Playlists")
              }
            </div>
          </ResizablePanel>
          <ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
        </>
      )}
    {(searchType === 1) && (
        <>
          <ResizablePanel defaultSize={50} minSize={20} className="flex flex-col">
            <div className="px-4 py-3 bg-[#181818] border-b border-[#282828]">
              <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Artists</span>
            </div>
            <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
              {searchLoading ? <p className="text-xs p-4 animate-pulse">Searching...</p> : 
                artistSearchData.length > 0 ? artistSearchData.map(item => (
                  <SearchResultItem key={item.spotifyID} item={item} onAction={handleInclusion} onExpand={handleExpand} />
                )) : renderEmpty("Artists")
              }
            </div>
          </ResizablePanel>
          <ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
        </>
      )}
        {(searchType === 1 || searchType === 2) && (
            <>
        <ResizablePanel defaultSize={albumSearchData.length > 0 ? 50 : 25} minSize={20} className="flex flex-col">
          <div className="px-4 py-3 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Albums</span>
          </div>
          <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
            {albumSearchLoading ? <p className="text-xs p-4 animate-pulse">Loading tracks...</p> :
              albumSearchData.length > 0 ? albumSearchData.map(item => (
                  <SearchResultItem key={item.spotifyID} item={item} onAction={handleInclusion} onExpand={handleExpand} />
              )) : renderEmpty("Albums")
            }
          </div>
        </ResizablePanel>
          <ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
          </>
      )}
        {(searchType === 3) && (
        <>
          <ResizablePanel defaultSize={trackSearchData.length > 0 ? 50 : 25} minSize={20} className="flex flex-col">
            <div className="px-4 py-3 bg-[#181818] border-b border-[#282828] flex justify-between items-center">
              <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Tracks</span>
            </div>
            <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
              {trackSearchLoading ? <p className="text-xs p-4 animate-pulse">Searching...</p> : 
                trackSearchData.length > 0 ? trackSearchData.map(item => (
                  <SearchResultItem key={item.spotifyID} item={item} onAction={handleInclusion} onExpand={handleExpand} />
                )) : renderEmpty("Tracks")
              }
            </div>
          </ResizablePanel>
          <ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
        </>
      )}
        {(searchType === 1 || searchType == 2) && (
        <ResizablePanel defaultSize={50} minSize={20} className="flex flex-col">
          <div className="px-4 py-3 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Tracks</span>
          </div>
          <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
            {trackSearchLoading ? <p className="text-xs p-4 animate-pulse">Loading tracks...</p> :
              trackSearchData.length > 0 ? trackSearchData.map(item => (
                <TrackResultItem key={item.spotifyID} item={item} onAction={handleInclusion} />
              )) : renderEmpty("Tracks")
            }
          </div>
        </ResizablePanel>
      )}
    </ResizablePanelGroup>
  );
}

function SearchBar() {
    const { searchType, setSearchType, search } = usePlaylistStore();
    const [query, setQuery] = React.useState("");

    useEffect(() => {
        const delayDebounceFn = setTimeout(() => {
            if (query || searchType === 0) {
                handleSearch();
            }
        }, 500);

        return () => clearTimeout(delayDebounceFn);
    }, [query, searchType]);

    const handleSearch = async () => {
        if (searchType !== 0 && !query) return;
        await search(query);
    }

    const handleTypeChange = (value: string) => {
        setSearchType(parseInt(value) as any);
    };

return (
    <div className="relative group w-full">
      <div className="absolute -inset-0.5 bg-gradient-to-r from-[#1DB954]/20 to-[#19e68c]/20 rounded-full blur opacity-0 group-focus-within:opacity-100 transition duration-500" />
      
      <div className="relative flex items-center w-full bg-[#242424] hover:bg-[#2a2a2a] focus-within:bg-[#333] rounded-full transition-all duration-300">
        <div className="pl-5 text-gray-400 group-focus-within:text-[#1DB954] transition-colors">
          <SearchIcon size={20} strokeWidth={2.5} />
        </div>

        <input
          placeholder={`Search for ${["Playlists", "Artists", "Albums", "Tracks"][searchType]}...`}
          className="w-full bg-transparent border-none py-4 px-4 text-sm text-white placeholder:text-gray-500 focus:outline-none focus:ring-0"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />

        <div className="flex items-center gap-1 pr-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button 
                variant="ghost" 
                className="rounded-full h-9 px-4 text-[10px] font-bold uppercase tracking-wider text-gray-400 hover:text-white hover:bg-[#ffffff1a]"
              >
                {["Playlist", "Artist", "Album", "Track"][searchType]}
              </Button>
            </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-48 bg-[#282828] border-[#3e3e3e] text-white">
                            <DropdownMenuLabel className="text-gray-400 text-[10px] uppercase tracking-widest">Search Category</DropdownMenuLabel>
                            <DropdownMenuSeparator className="bg-[#3e3e3e]" />
                            <DropdownMenuRadioGroup value={searchType.toString()} onValueChange={handleTypeChange}>
                                <DropdownMenuRadioItem value="1" className="focus:bg-[#333] cursor-pointer">Artist</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="2" className="focus:bg-[#333] cursor-pointer">Album</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="3" className="focus:bg-[#333] cursor-pointer">Track</DropdownMenuRadioItem>
                                <DropdownMenuRadioItem value="0" className="focus:bg-[#333] cursor-pointer">Playlist</DropdownMenuRadioItem>
                            </DropdownMenuRadioGroup>
                        </DropdownMenuContent>
</DropdownMenu>

          <Button 
            onClick={handleSearch}
            className="rounded-full bg-white hover:scale-105 active:scale-95 text-black font-bold h-9 px-6 transition-all"
          >
            Go
          </Button>
        </div>
      </div>
    </div>
  );
}
