# Lyra

## Notes
* Need a way to tell Lyra that a song was a good match.
* Need a way to tell Lyra what a could recommend is
	* ie This song is good and IMO this other song would be nice
	* Note use user playlists as a basis.
		* ***How to get user playlists?***
* Need a way to tell Lyra which camera it should use
	* Smart check?
		* If screen is dark, or not enough detail can be made out
			* Ignore camera
		* If both screens are dark
			* Ignore both
		* or maybe it's night? \=
* Use sound somehow?
	* Loud means lots of conversation?
	* Quiet means people are listening to music?
	* What if the music itself is loud?
		* Clean out music and read what remains?

## Determining what songs to play for *me*
* Create a profile per song, per album (Global & Local)
	* Based on these following things:
		* Metrics pulled in from the CV while song was playing
			* Lighting
			* Number of people in the image
				* Low ranking, face recognition is unreliable
				* People's face may be out of image
				* Increase ranking if there's someway to identify bodies.
			* Change since last image
			* Diversity of color
			* Dominant colours
		* Time listened to
			* More specifically this is to get rid of the songs that are only
			* listened to for <1 min or <30 sec)
		* Length
			* Long songs would probably fit together better than having a long
			  song followed by a short one.
		* Next song played
		* Previous song played.
		* Genre
		* Artist
		* Weather
	* Using a song or albums profile use dating(?) algorithms to determine other
	  songs that would be a good `match` for this one
		* This needs to be done intelligently
			* Match closely:
				* All of CV
				* Time
				* Length
				* Genre
			* Different
				* Artist
				* Next song played

	* Of course these metrics need to be prioritized such that something like
	  the length of a song does not take priority over something like the genre

* As well create a user profile (Global & Local)
	* Song that they often listen to
		* Frequency. Needs to be done in such a way that the same subset of
		  songs do not get overplayed.
			* perhaps have a random element that will allow the threshold of a
			  song being considered `new` is relaxed.
	* New songs
		* that match. use above algorithm
	* Songs that they frequently skip
		* Deletion recommend?
	* Songs that they have never listened to(?)
		* (?) This could be either that we should play these songs or not
			* If they don't want these songs played perhaps we should advise
			  them to delete?
	* Location
	* Popularity
		* Perhaps secondarily based on location.

## Determining what songs to play for the *group*
* Maintain a list of intersections (Local)
	* Songs
	* Albums
	* Genres
* Pass user profile along (Global){Privacy}
* Using the list of intersections
	* Matching songs/albums are definitely played
	* Matching Genres will have to be a bit more complicated
		* Lower score than songs/albums for sure
		* If the two genres are an exact match Check global song/album profile
			* Check if any songs/albums on the host device have a good rating
			  with songs/albums on the guest device

## App itself
1. Upon start, user is greeted with the Player Screen [Fig. 1]
	* Within this screen is the bulk of Lyra's functionality
	* Based on metrics that are immediately and quickly accessible
		* Lyra will pick a song, which fills in the center icon element
	* Metrics that take a bit more time to analyze (namely the CV) begin now
		* Lyra will then use these new metrics to determine which song to
		  recommend next.
	* The user may pull the song down to inform Lyra this current song is a good
	  match.
