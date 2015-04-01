import json
import numpy
import matplotlib.pyplot as plt
import sqlite3, sys
import random

def pdf_mulvariate_gauss(x, mu, cov):
    part1 = 1 /( (2 * numpy.pi)**(13) * numpy.linalg.det(cov)) ** (0.5)
    part2 = (-1 / 2) * (x - mu).T.dot(numpy.linalg.inv(cov).dot(x - mu))
    return float(part1 * numpy.exp(part2))

def combine_points(alpha, beta, dim = 13):
    sigma = alpha["sigma"] + beta["sigma"]
    usigma = alpha["usigma"] + beta["usigma"]
    n = alpha["n"] + beta["n"]
    mean = usigma[0] / n

    cov = (sigma -\
          (numpy.transpose([mean]*dim) * usigma) -\
          (numpy.array([mean]*dim) * usigma.T) +\
          (n * numpy.transpose([mean]) * mean)
          )/(n - 1.0)
    new = {}
    new["mean"] = mean
    new["cov"] = cov
    new["sigma"] = sigma
    new["usigma"] = usigma
    new["n"] = n
    return new

def distance(center, point):
    #return numpy.linalg.norm(center["mean"] - point["mean"])
    return pdf_mulvariate_gauss(point["mean"], center["mean"], center["cov"])

def center(data, songs, n_centers):
    # generate the centers (center, points assigned to it)
    centers = {}
    for i in range(n_centers):
        song = random.choice(songs)
        centers[i] = {}
        centers[i]["data"] = data[song]
        centers[i]["cluster"] = []

    delta = 200
    olddelta = 201
    rounds = 0
    while True:
        rounds += 1
        # Cluster according to the centers
        distances = {}
        totals = [0 for i in range(n_centers)]
        for song in songs:
            # Calculate the distance between each center and a song
            distances[song] = [distance(center["data"], data[song]) for center in centers.values()]
            totals = [x + y for x, y in zip(totals, distances[song])]

        for song in songs:
            distances[song] = [i/total for i, total in zip(distances[song],
                totals)]
            # Append this song to the center it's closest to
            centers[distances[song].index(max(distances[song]))]["cluster"].append(data[song])
        if delta == olddelta or rounds > 20:
            break
        else:
            olddelta = delta
        
        delta = 0
        # Re center according to new clusters
        for i in range(int(n_centers)):
            # check whether this center has anything clustered to it
            if centers[i]["cluster"]:
                new_center = centers[i]["cluster"][0]
                for song in centers[i]["cluster"][1:]:
                    new_center = combine_points(new_center, song)
                delta += distance(centers[i]["data"], new_center)
                centers[i] = {}
                centers[i]["data"] = new_center
                centers[i]["cluster"] = []
    return centers

def print_centers(centers, song):
    for center in centers:
        print "Cluster %d" % center
        for song in centers[center]["cluster"]:
            print "\t%s" % song["name"][:76]

if __name__ == "__main__":
    n_centers = int(sys.argv[1])

    # Analysis of how many centers to use
#    lg = []
#    lg1 = []
#    for j in range(2,n_centers + 1):
#        temp_lg = []
#        temp_lg1 = []
#        for x in range(10):
#            centers = center(songs[0::2], data, j)
#            for i in range(int(j)):
#                for song in songs:
#                    n = distance(data[song], centers[i]["data"])
#                    data[song]["centers"].append(n)
#
#            logp1 = 0
#            for song in songs[0::2]:
#                total = sum(data[song]["centers"])
#                for cen in data[song]["centers"]:
#                    if total != 0:
#                        data[song]["p"] += cen/total * cen
#            for song in songs[0::2]:
#                if data[song]["p"] != 0:
#                    logp1 += numpy.log(data[song]["p"])
#
#            logp = 0
#            for song in songs[1::2]:
#                total = sum(data[song]["centers"])
#                for cen in data[song]["centers"]:
#                    if total != 0:
#                        data[song]["p"] += cen/total * cen
#            for song in songs[1::2]:
#                if data[song]["p"] != 0:
#                    logp += numpy.log(data[song]["p"])
#            temp_lg1.append(logp1)
#            temp_lg.append(logp)
#        lg1.append(max(temp_lg1))
#        lg.append(max(temp_lg))
#    for i in lg1:
#        print i
#    print
#    print 
#    for i in lg:
#        print i
    centers = center(n_centers)
    # print the results
