import numpy
import sqlite3
import sys
import random
# Perform user analysis
def pdf_mulvariate_gauss(x, mu, cov):
    dim = x.shape[0]
    part1 = 1 /( (2 * numpy.pi)**(dim) * numpy.linalg.det(cov)) ** (0.5)
    part2 = (-1 / 2) * (x - mu).T.dot(numpy.linalg.inv(cov).dot(x - mu))
    return float(part1 * numpy.exp(part2))

# class to store data which we have a covariance stored for
class Data():
    def __init__(self, dtype=None, name=None, n=None):
        if dtype and name and n:
            self.mean = numpy.load("%s/%s-mean.npy" % (dtype, name))
            self.cov = numpy.load("%s/%s-cov.npy" % (dtype, name))
            self.sigma = numpy.load("%s/%s-sigma.npy" % (dtype, name))
            self.usigma = numpy.load("%s/%s-usigma.npy" % (dtype, name))
            self.n = n
            self.dim = self.mean.shape[0]

    # combining two data points together
    def __add__(self, other):
        sigma = self.sigma + other.sigma
        usigma = self.usigma + other.usigma
        n = self.n + other.n
        mean = usigma[0] / n
        cov = (sigma -\
              (numpy.transpose([mean]*self.dim) * usigma) -\
              (numpy.array([mean]*self.dim) * usigma.T) +\
              (n * numpy.transpose([mean]) * mean)
              )/(n - 1.0)
        new = Data()
        new.mean = mean
        new.cov = cov
        new.sigma = sigma
        new.usigma = usigma
        new.n = n
        new.dim = self.dim
        return new

    def __sub__(self, other):
        return pdf_mulvariate_gauss(other.mean, self.mean, self.cov)
        

class Song():
    def __init__(self, data=None):
        if data:
            self.mfcc = Data("mfcc", data[3], data[4])
            self.loudness = Data("loudness", data[5], data[6])
            self.name = data[0]
            self.key = data[3]
            self.centers = []
            self.p = 0

    def __str__(self):
        return self.name

    # Distance function between two songs/centers
    def __sub__(self, other):
        return (self.mfcc - other.mfcc)# + (self.loudness - other.loudness)

class Center(Song):
    def __init__(self, data):
        self.mfcc = data.mfcc
        self.loudness = data.loudness
        self.name = "center"
        self.cluster = []

    def __str__(self):
        result = "Center:\n"
        for song in self.cluster:
            result += "\t%s\n" % song
        return result

    def __sub__(self, other):
        dist1 = (self.mfcc - other.mfcc)# + (self.loudness - other.loudness)
        dist2 = (other.mfcc - self.mfcc)# + (self.loudness - other.loudness)
        return (dist1 + dist2)/2.0

    # Combining two songs to creat a new center
    def __add__(self, other):
        self.mfcc += other.mfcc
        self.cluster = []
        return self


def center(songs, n_centers):
    # generate the centers
    centers = []
    for i in range(n_centers):
        new = Center(random.choice(songs))
        centers.append(new)        

    # initiate variables to end this loop
    delta = 200
    prev_delta = 201
    rounds = 0
    
    # Let's center the songs
    while True:
        rounds += 1

        distances = {}
        totals = [0 for i in range(n_centers)]
        # Match each song with it's most likely center
        for song in songs:
            # Calculate the distance between each center and a song
            distances[song.key] = [center - song for center in centers]
            # Calculate the total distances for normalization
            totals = [x + y for x, y in zip(totals, distances[song.key])]
        
        for song in songs:
            distances[song.key] = [i/total for i, total in
                zip(distances[song.key], totals)]
            centers[distances[song.key].index(max(distances[song.key]))].cluster.append(song)

        # Check if we can exit now
        if delta == prev_delta or rounds > 20:
            break
        else:
            olddelta = delta

        delta = 0
        new_centers = []
        # re-centers with our new clusters
        for i in range(n_centers):
            if centers[i].cluster:
                new_center = Center(centers[i].cluster[0])
                for song in centers[i].cluster[1:]:
                    new_center += song
                delta += center - new_center
                centers[i] = new_center
        

    return centers

def get_songs(ids):
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    songs = c.execute("SELECT * FROM songs").fetchall()
    data = [Song(item) for item in songs]
    #data = [song(songs[i[0]]) for i in ids]
    return data

def model_selection(songs, n_centers):
    lg = []
    lg1 = []
    trials = 3
    count = 0
    for j in range(2,n_centers + 1):
        temp_lg = []
        temp_lg1 = []
        for _ in range(trials):
            centers = center(songs[0::2], j)
            for i in range(int(j)):
                for song in songs:
                    n = centers[i] - song
                    song.centers.append(n)

            logp1 = 0
            for song in songs[0::2]:
                total = sum(song.centers)
                for cen in song.centers:
                    if total != 0:
                        song.p += cen/total * cen
            for song in songs[0::2]:
                if song.p != 0:
                    logp1 += numpy.log(song.p)

            logp = 0
            for song in songs[1::2]:
                total = sum(song.centers)
                for cen in song.centers:
                    if total != 0:
                        song.p += cen/total * cen
            for song in songs[1::2]:
                if song.p != 0:
                    logp += numpy.log(song.p)
            temp_lg1.append(logp1)
            temp_lg.append(logp)
            count += 1
            print "Progress: %f%%" % (count/float(trials * n_centers)*100.0)
        lg1.append(max(temp_lg1))
        lg.append(max(temp_lg))
    for i in lg1:
        print i
    print
    print 
    for i in lg:
        print i

def save_centers(centers, filename):
    f = open(filename, "w")
    for center in centers:
        f.write(str(center))
    f.close()

if __name__ == "__main__":
    if len(sys.argv) > 1:
        # model selection
        if "-m" in sys.argv:
            songs = get_songs(range(1200))
            model_selection(songs, 50)
        if "-s" in sys.argv:
            songs = get_songs(range(2400))
            centers = center(songs, 15)
            save_centers(centers, "user")
        if "-r" in sys.argv:
            # Create user association database
            conn = sqlite3.connect("users.db")
            c = conn.cursor()
            c.execute("CREATE TABLE IF NOT EXISTS users (name varchar(255));")
            c.execute("CREATE TABLE IF NOT EXISTS songs (userid int, songid int);")
            firstids = c.execute("SELECT songid FROM songs where userid=1").fetchall()
            songs = get_songs(range(1200))
            centers = center(songs, 15)
            save_centers(centers, "user1")

            secondids = c.execute("SELECT songid FROM songs where userid=2").fetchall()
            songs2 = get_songs(range(1200,2400))
            centers2 = center(songs2, 15)
            save_centers(centers2, "user2")

            similar_centers = []
            for center2 in centers2:
                dists = [center2 - center1 for center1 in centers]
                similar_centers.append(dists.index(max(dists)))
            # find the most common center:
            best = max(set(similar_centers), key=similar_centers.count)
            f = open("combined", "w")
            f.write(str(centers[best]))
            for i in range(15):
                if similar_centers[i] == best:
                    f.write(str(centers2[i]))
            #print_centers(centers2, songs2)

            conn.commit()
            conn.close()
