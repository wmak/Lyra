import numpy
import sqlite3
from distance_calculator import center, distance, print_centers
# Perform user analysis
def get_songs(ids):
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    songs = c.execute("SELECT * FROM songs").fetchall()
    data = {}
    for i in ids:
        song = songs[i[0]]
        current = {}
        current["mean"]= numpy.load("mfcc/%s-mean.npy" % (song[3]))
        current["cov"]= numpy.load("mfcc/%s-cov.npy" % (song[3]))
        current["sigma"]= numpy.load("mfcc/%s-sigma.npy" % (song[3]))
        current["usigma"]= numpy.load("mfcc/%s-usigma.npy" % (song[3]))
        current["n"] = song[4]
        current["name"] = song[0]
        current["key"] = song[3]
        current["centers"] = []
        current["p"] = 0
        data[song[3]] = current

    songs = data.keys()
    return data, songs

if __name__ == "__main__":
    firstrun = False
    # Extract the song data from the database.
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    all_songs = c.execute("SELECT * FROM songs").fetchall()
    conn.close()

    # Create user association database
    conn = sqlite3.connect("users.db")
    c = conn.cursor()
    c.execute("CREATE TABLE IF NOT EXISTS users (name varchar(255));")
    c.execute("CREATE TABLE IF NOT EXISTS songs (userid int, songid int);")
    firstids = c.execute("SELECT songid FROM songs where userid=1").fetchall()
    data, songs = get_songs(firstids)

    centers = center(data, songs, 15)
    #print_centers(centers, songs)

    secondids = c.execute("SELECT songid FROM songs where userid=2").fetchall()
    secondids = secondids[:16] + secondids[17:]
    data2, songs2 = get_songs(secondids)

    centers2 = center(data2, songs2, 15)
    similar_centers = []
    for center2 in centers2:
        dists = [distance(centers[center1]["data"], centers2[center2]["data"]) for
                center1 in centers]
        similar_centers.append(dists.index(max(dists)))
    # find the most common center:
    best = max(set(similar_centers), key=similar_centers.count)
    for song in centers[best]["cluster"]:
        print song["name"]
    print "=" * 80
    for i in range(15):
        if similar_centers[i] == best:
            for song in centers2[i]["cluster"]:
                print song["name"]
    #print_centers(centers2, songs2)

    conn.commit()
    conn.close()
