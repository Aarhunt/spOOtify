"use client"

import type { ModelItemResponse, ModelItemType } from "@/client/types.gen";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, Plus, Music } from "lucide-react";
import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import { usePlaylistStore } from "../stores/playlist.store";

export default function Summary() {
    return (
        <div className="flex flex-col h-full gap-4">
            <div className="flex items-center justify-between px-2">
                <h2 className="text-xl font-bold tracking-tight">Library Summary</h2>
                <SelectType />
            </div>
<div className="w-full flex-1 min-h-0 px-4 mt-0 overflow-hidden">
                <ResultBox />
            </div>
        </div>
    )
}

export function SelectType() {
    const { setSummaryType, summary } = usePlaylistStore();
    return (
        <Select onValueChange={(v) => {
            setSummaryType(parseInt(v) as ModelItemType);
            summary();
        }} defaultValue="4">
            <SelectTrigger className="w-[140px] bg-[#282828] border-none text-xs font-bold uppercase tracking-wider h-9 rounded-full focus:ring-0">
                <SelectValue placeholder="View" />
            </SelectTrigger>
            <SelectContent className="bg-[#282828] border-[#3e3e3e] text-white">
                <SelectItem value="4">All Overview</SelectItem>
                <SelectItem value="0">Playlists Only</SelectItem>
                <SelectItem value="1">Artists View</SelectItem>
                <SelectItem value="2">Albums View</SelectItem>
                <SelectItem value="3">Tracks View</SelectItem>
            </SelectContent>
        </Select>
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
        return null;
      case 0: 
        return null
      default:
        return <Badge className="bg-[blue]/10 text-[blue] border-none text-[10px] h-5 px-1.5">Included</Badge>;
    }
  };

return (
    <div 
      className="group flex items-center justify-between p-2 transition-all hover:bg-[#ffffff1a] rounded-md cursor-pointer"
      onClick={() => item.spotifyID && item.itemType && onExpand(item.spotifyID, item.itemType)}
    >
      <div className="flex items-center gap-3 overflow-hidden">
        <Avatar className="h-10 w-10 rounded shadow-lg">
          <AvatarImage src={imageUrl} className="object-cover" />
          <AvatarFallback className="bg-[#333] text-[10px]">SP</AvatarFallback>
        </Avatar>
        
        <div className="flex flex-col min-w-0">
          <span className="font-semibold text-sm truncate text-white group-hover:text-[#1DB954] transition-colors">
            {item.name}
          </span>
          <div className="flex items-center gap-2">
            {item.itemType === 2 && <span className="text-[10px] text-gray-400 uppercase">{item.sortdata}</span>}
            {getInclusionBadge()}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <Button 
          size="icon" 
          variant="ghost" 
          className={cn(
            "h-8 w-8 rounded-full",
            item.included === 1 ? "bg-[#1DB954] text-black hover:bg-[#1ed760]" : "hover:bg-[#333] text-gray-400"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);
          }}
        >
          {item.included === 1 ? <Check className="h-4 w-4 stroke-[3]" /> : <Plus className="h-4 w-4" />}
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
        return null;
      case 2: 
        return null;
      case 3: 
        return <Badge className="bg-[blue]/10 text-[blue] border-none text-[10px] h-5 px-1.5">Included</Badge>;
      case 4: 
        return <Badge className="bg-[orange]/10 text-[orange] border-none text-[10px] h-5 px-1.5">Included</Badge>;
      default:
        return null;
    }
  };

return (
    <div 
      className="group flex items-center justify-between p-2 transition-all hover:bg-[#ffffff1a] rounded-md cursor-pointer"
      onClick={() => item.spotifyID && item.itemType && onExpand(item.spotifyID, item.itemType)}
    >
      <div className="flex items-center gap-3 overflow-hidden">
        <Avatar className="h-10 w-10 rounded shadow-lg">
          <AvatarImage src={imageUrl} className="object-cover" />
          <AvatarFallback className="bg-[#333] text-[10px]">SP</AvatarFallback>
        </Avatar>
        
        <div className="flex flex-col min-w-0">
          <span className="font-semibold text-sm truncate text-white group-hover:text-[#1DB954] transition-colors">
            {item.name}
          </span>
          <div className="flex items-center gap-2">
            {item.itemType === 2 && <span className="text-[10px] text-gray-400 uppercase">{item.sortdata}</span>}
            {getInclusionBadge()}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <Button 
          size="icon" 
          variant="ghost" 
          className={cn(
            "h-8 w-8 rounded-full",
            item.included === 1 ? "bg-[#1DB954] text-black hover:bg-[#1ed760]" : "hover:bg-[#333] text-gray-400"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, item.included === 1);
          }}
        >
          {item.included === 1 ? <Check className="h-4 w-4 stroke-[3]" /> : <Plus className="h-4 w-4" />}
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
        return null;
      case 2:
        return null;
      case 3:
        return <Badge className="bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 border-none h-5 px-1.5 text-[10px]">PR</Badge>;
      case 4:
        return <Badge className="bg-orange-500-500/10 text-orange-500 hover:bg-orange-500/20 border-none h-5 px-1.5 text-[10px]">PR</Badge>;
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
      </div>
    </div>
  );
}

function ResultBox() {
    const { summaryType } = usePlaylistStore();

    return (
        <> { summaryType == 4 ? <ResultBoxSummary /> : <ResultBoxExpand /> } </>
    )
}

function ResultBoxSummary() {
  const { playlistSummaryData, artistSummaryData, albumSummaryData, trackSummaryData, summaryPlaylistsLoading, summaryArtistsLoading, summaryAlbumsLoading, summaryTracksLoading, includeItem, undoIncludeItem, undoIncludePlaylist, includePlaylist} = usePlaylistStore();

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

  const renderEmpty = (label: string) => (
    <div className="flex flex-col items-center justify-center py-20 opacity-20">
      <Music className="h-12 w-12 mb-4" />
      <p className="text-sm font-medium">No {label} Found</p>
    </div>
  ); 



  return (
<ResizablePanelGroup className="flex-1 rounded-xl bg-[#121212] border border-[#282828] overflow-hidden max-h-[355px] ">
    <ResizablePanel defaultSize={25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Playlists</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryPlaylistsLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading playlists...</p>
          ) : playlistSummaryData.length > 0 ? (
            playlistSummaryData.map((item) => (
              <PlaylistResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => {}}
              />
            ))
          ) : renderEmpty("Playlists")}
        </div>
      </ResizablePanel>
<ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
<ResizablePanel defaultSize={25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Artists</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryArtistsLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading artists...</p>
          ) : artistSummaryData.length > 0 ? (
            artistSummaryData.map((item) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => {}}
              />
            ))
          ) : renderEmpty("Artists")}
        </div>
      </ResizablePanel>
<ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
<ResizablePanel defaultSize={25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828] ">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Albums</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryAlbumsLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading albums...</p>
          ) : albumSummaryData.length > 0 ? (
            albumSummaryData.map((item) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => {}}
              />
            ))
          ) : renderEmpty("Albums")}
        </div>
      </ResizablePanel>
<ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
<ResizablePanel defaultSize={25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Tracks</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryTracksLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading tracks...</p>
          ) : trackSummaryData.length > 0 ? (
            trackSummaryData.map((item) => (
            <SearchResultItem
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => {}}
              />
            ))
          ) : renderEmpty("Tracks")}
        </div>
      </ResizablePanel>
      </ResizablePanelGroup>
  )
}

function ResultBoxExpand() {
  const { playlistSummaryData, artistSummaryData, albumSummaryData, trackSummaryData, artistSummaryExpandData, albumSummaryExpandData, summaryPlaylistsLoading, summaryArtistsLoading, summaryAlbumsLoading, summaryTracksLoading, trackSummaryExpandData, summaryType, includeItem, undoIncludeItem, getSummaryAlbumsFromArtist, getSummaryTracksFromAlbum, setCurrentSummaryArtist, setCurrentSummaryAlbum, undoIncludePlaylist, includePlaylist } = usePlaylistStore();
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
            getSummaryAlbumsFromArtist(id) 
            setCurrentSummaryArtist(id) 
            break;
        case 2:
            getSummaryTracksFromAlbum(id)
            setCurrentSummaryAlbum(id)
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

    const artists = summaryType == 1 ? artistSummaryData : artistSummaryExpandData
    const albums = summaryType == 2 ? albumSummaryData : albumSummaryExpandData
    const tracks = summaryType == 3 ? trackSummaryData : trackSummaryExpandData

  return (
<ResizablePanelGroup className="flex-1 rounded-xl bg-[#121212] border border-[#282828] overflow-hidden max-h-[355px] ">
    {
        summaryType == 0 ? 
        <>
<ResizablePanel defaultSize={albumSummaryData.length > 0 ? 50 : 50} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Playlists</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryPlaylistsLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading playlists...</p>
          ) : playlistSummaryData.length > 0 ? (
            playlistSummaryData.map((item) => (
              <PlaylistResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => {}}
              />
            ))
          ) : renderEmpty("Playlists")}
        </div>
      </ResizablePanel>
      </> : null
    }
    {
        summaryType == 1 ? 
        <>
<ResizablePanel defaultSize={50} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Artists</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryArtistsLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading artists...</p>
          ) : artists.length > 0 ? (
            artists.map((item) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={handleExpand}
              />
            ))
          ) :  renderEmpty("Artists")}
        </div>
      </ResizablePanel>
    <ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
      </> : null
    }
      {
        summaryType == 1 || summaryType == 2 ?  
        <> 
<ResizablePanel defaultSize={albumSummaryExpandData.length > 0 ? 50 : 25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Albums</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
              {summaryAlbumsLoading ? (
                <p className="text-xs p-4 animate-pulse">Loading albums...</p>
              ) : albums.length > 0 ? (
                albums.map((item) => (
                  <SearchResultItem 
                    key={item.spotifyID} 
                    item={item} 
                    onAction={handleInclusion} 
                    onExpand={handleExpand}
                  />
                ))
              ) :  renderEmpty("Albums")}
            </div>
          </ResizablePanel>
<ResizableHandle className="w-[1px] bg-[#282828] hover:bg-[#1DB954] transition-colors" />
        </> : null
      }
      {
        summaryType == 1 || summaryType == 2 ? 
        <>
<ResizablePanel defaultSize={trackSummaryExpandData.length > 0 ? 50 : 25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Tracks</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryTracksLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading tracks...</p>
          ) : tracks.length > 0 ? (
            tracks.map((item) => (
            <TrackResultItem
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
              />
            ))
          ) :  renderEmpty("Tracks")}
        </div>
      </ResizablePanel>
    </> : null
      }
      {
        summaryType == 3 ? 
        <>
<ResizablePanel defaultSize={trackSummaryExpandData.length > 0 ? 50 : 25} className="flex flex-col">
        <div className="px-4 py-2 bg-[#181818] border-b border-[#282828]">
            <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">Tracks</span>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
          {summaryTracksLoading ? (
            <p className="text-xs p-4 animate-pulse">Loading tracks...</p>
          ) : tracks.length > 0 ? (
            tracks.map((item) => (
            <SearchResultItem
                key={item.spotifyID} 
                item={item} 
                onAction={handleInclusion} 
                onExpand={() => void {}}
              />
            ))
          ) :  renderEmpty("Tracks")}
        </div>
      </ResizablePanel>
    </> : null
      }
    </ResizablePanelGroup>
  );
}

