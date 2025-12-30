"use client"

import { useSummaryStore } from "@/components/stores/summary.store"
import type { ModelItemResponse, ModelItemType } from "@/client/types.gen";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, Plus, X, Music } from "lucide-react";
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

export default function Summary() {
    return (
        <>
        <SelectType />
        <ResultBox />
        </>
    )
}

export function SelectType() {
    const { setSummaryType, summary} = useSummaryStore()
   return (
       <Select onValueChange={(v) => 
           {setSummaryType(parseInt(v) as ModelItemType)
           summary()}
       } defaultValue="1">
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="0">Playlist</SelectItem>
            <SelectItem value="1">Artist</SelectItem>
            <SelectItem value="2">Album</SelectItem>
            <SelectItem value="3">Track</SelectItem>
          </SelectContent>
        </Select>
    )
}

interface SearchResultItemProps {
  item: ModelItemResponse;
  index: number;
  onAction: (id: string, include: boolean, type: ModelItemType, index: number, undo: boolean) => void;
  onExpand: (id: string, type: ModelItemType) => void;
}

export function SearchResultItem({ item, index, onAction, onExpand }: SearchResultItemProps) {
  const imageUrl = item.icon && item.icon.length > 0 ? item.icon[0].url : "";
  
  const getInclusionBadge = () => {
    switch (item.included) {
      case 1: 
        return <Badge className="bg-green-500 hover:bg-green-600">Included</Badge>;
      case 2: 
        return <Badge variant="destructive">Excluded</Badge>;
      case 3: 
        return <Badge className="bg-blue-500 hover:bg-blue-600">Included</Badge>;
      case 4: 
        return <Badge className="bg-orange-500 hover:bg-orange-600">Excluded</Badge>;
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
              item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, index, item.included === 1);}}
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
              item.spotifyID && item.itemType && onAction(item.spotifyID, false, item.itemType, index, item.included === 2);}}
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

interface TrackResultItemProps {
  item: ModelItemResponse;
  index: number;
  onAction: (id: string, include: boolean, type: ModelItemType, index: number, undo: boolean) => void;
}

export function TrackResultItem({ item, index, onAction }: TrackResultItemProps) {
  
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
            item.included === 1 && "text-green-500 bg-green-500/10 hover:bg-green-500/20 hover:text-green-600"
          )}
          onClick={(e) => {
            e.stopPropagation();
            item.spotifyID && item.itemType && onAction(item.spotifyID, true, item.itemType, index, item.included === 1);
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
            item.spotifyID && item.itemType && onAction(item.spotifyID, false, item.itemType, index, item.included === 2);
          }}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>
    </div>
  );
}

function ResultBox() {
  const { mainArtistData, mainAlbumData, mainTrackData, artistData, albumData, summaryLoading, trackData, albumLoading, trackLoading, summaryType, includeItem, undoIncludeItem, getAlbumsFromArtist, getTracksFromAlbum, setCurrentArtist, setCurrentAlbum } = useSummaryStore();
  // You'll need an action in a store to handle the actual DB update
  const handleInclusion = (id: string, include: boolean, type: ModelItemType, index: number, undo: boolean) => { 
    undo ? undoIncludeItem(id, include, type, index) : includeItem(id, include, type, index)
  };

  const handleExpand = (id: string, type: ModelItemType) => {
    switch (type) {
        case 1:
            getAlbumsFromArtist(id) 
            setCurrentArtist(id) 
            break;
        case 2:
            getTracksFromAlbum(id)
            setCurrentAlbum(id)
            break;
        default:
            break;
    } 
  }

    const artists = summaryType == 1 ? mainArtistData : artistData
    const albums = summaryType == 2 ? mainAlbumData : albumData
    const tracks = summaryType == 3 ? mainTrackData : trackData

  return (
    <ResizablePanelGroup direction="horizontal" className="min-h-[400px] rounded-lg border">
    {
        summaryType == 1 ? 
        <>
      <ResizablePanel defaultSize={mainAlbumData.length > 0 ? 50 : 50}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[600px]">
          {summaryLoading ? (
            <p>Loading results...</p>
          ) : artists.length > 0 ? (
            artists.map((item, index) => (
              <SearchResultItem 
                key={item.spotifyID} 
                item={item} 
                index = {index}
                onAction={handleInclusion} 
                onExpand={handleExpand}
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No results found.</p>
          )}
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      </> : null
    }
      {
        summaryType == 1 || summaryType == 2 ?  
        <> 
          <ResizablePanel defaultSize={albumData.length > 0 ? 50 : 25}>
            <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[600px]">
              {albumLoading ? (
                <p>Loading results...</p>
              ) : albums.length > 0 ? (
                albums.map((item, index) => (
                  <SearchResultItem 
                    key={item.spotifyID} 
                    item={item} 
                    index = {index}
                    onAction={handleInclusion} 
                    onExpand={handleExpand}
                  />
                ))
              ) : (
                <p className="text-muted-foreground text-center py-10">No results found.</p>
              )}
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle />
        </> : null
      }
      {
        summaryType == 1 || summaryType == 2 ? 
        <>
      <ResizablePanel defaultSize={trackData.length > 0 ? 50 : 25}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[600px]">
          {trackLoading ? (
            <p>Loading results...</p>
          ) : tracks.length > 0 ? (
            tracks.map((item, index) => (
            <TrackResultItem
                key={item.spotifyID} 
                item={item} 
                index={index}
                onAction={handleInclusion} 
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No results found.</p>
          )}
        </div>
      </ResizablePanel>
    </> : null
      }
      {
        summaryType == 3 ? 
        <>
      <ResizablePanel defaultSize={trackData.length > 0 ? 50 : 25}>
        <div className="flex flex-col gap-2 p-4 overflow-y-auto max-h-[600px]">
          {trackLoading ? (
            <p>Loading results...</p>
          ) : tracks.length > 0 ? (
            tracks.map((item, index) => (
            <SearchResultItem
                key={item.spotifyID} 
                item={item} 
                index={index}
                onAction={handleInclusion} 
                onExpand={() => void {}}
              />
            ))
          ) : (
            <p className="text-muted-foreground text-center py-10">No results found.</p>
          )}
        </div>
      </ResizablePanel>
    </> : null
      }
    </ResizablePanelGroup>
  );
}

