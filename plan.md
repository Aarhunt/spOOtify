# Autistify

## Backend

## Frontend
- Playlists van Spotify worden ingeladen
- Gaat vooral om nieuwe playlist maken. 
- 

## Functionalities

### Making a new playlist
- You make a playlist. In this playlist you can include any songs you want.
- You can select whole artists or albums as a whole. This should exclude "remastered, instrumental, etc"
- Ideally, this updates automatically if new stuff gets added. 
- You can also specifically exclude albums, or exclude songs from albums. 

### Merging playlists
- Some playlists are simply mergings of other playlists.
- You can make a playlist that instead of taking in an artist and its songs, it takes your own playlist (and updates automagically)

## In total, this feels like a container design pattern. Every playlist, album, artist, etc is a container that contains certain containers, but also excludes certain containers. 
That is what I should make. 
