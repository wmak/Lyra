import numpy
import sqlite3, sys
import random

def pdf_mulvariate_gauss(x, mu, cov):
    part1 = 1 / (2 * numpy.pi)**(6.5) * 1 / numpy.linalg.det(cov) ** (1.0/2.0)
    part2 = (-1 / 2) * (x - mu).T.dot(numpy.linalg.inv(cov).dot(x - mu))
    return float(part1 * numpy.exp(part2))

def combine_points(alpha, beta):
    sigma = alpha["sigma"] + beta["sigma"]
    usigma = alpha["usigma"] + beta["usigma"]
    n = alpha["n"] + beta["n"]
    mean = usigma[0] / n
    cov = (sigma -\
          ((numpy.transpose([mean]*13)) * usigma) -\
          (usigma.T * (numpy.array([mean]*13))) +\
          (mean * numpy.transpose([mean]))) / (n - 1)
    new = {}
    new["mean"] = mean
    new["cov"] = cov
    new["sigma"] = sigma
    new["usigma"] = usigma
    new["n"] = n
    return new

def distance(center, point):
    return pdf_mulvariate_gauss(point["mean"], center["mean"], center["cov"])

if __name__ == "__main__":
    n_centers = sys.argv[1]
    conn = sqlite3.connect("music.db")
    c = conn.cursor()
    songs = c.execute("SELECT * FROM songs").fetchall()
    data = {}
    for song in songs:
        current = {}
        current["mean"]= numpy.load("mfcc/%s-mean.npy" % (song[3]))
        current["cov"]= numpy.load("mfcc/%s-cov.npy" % (song[3]))
        current["sigma"]= numpy.load("mfcc/%s-sigma.npy" % (song[3]))
        current["usigma"]= numpy.load("mfcc/%s-usigma.npy" % (song[3]))
        current["n"] = song[4]
        current["name"] = song[0]
        data[song[3]] = current

    songs = data.keys()

    # generate the centers (center, points assigned to it)
    centers = {}
    for i in range(int(n_centers)):
        song = random.choice(songs)
        centers[i] = {}
        centers[i]["data"] = data[song]
        centers[i]["cluster"] = []

    delta = 10
    while True:
        # Cluster according to the centers
        for song in songs:
            # Calculate the distance between each center and a song
            distances = [distance(center["data"], data[song]) for center in centers.values()]
            print distances
            # Append this song to the center it's closest to
            centers[distances.index(max(distances))]["cluster"].append(data[song])
        if delta < 0.1:
            break
        else:
            print delta
        
        delta = 0
        # Re center according to new clusters
        for i in range(int(n_centers)):
            # check whether this center has anything clustered to it
            if centers[i]["cluster"]:
                new_center = centers[i]["cluster"][0]
                for song in centers[i]["cluster"][1:]:
                    new_center = combine_points(new_center, song)
                centers[i] = {}
                centers[i]["data"] = new_center
                centers[i]["cluster"] = []

    # print the results
    for i in range(int(n_centers)):
        print "Cluster %d" % i
        for song in centers[i]["cluster"]:
            print "\t%s" % song["name"]

