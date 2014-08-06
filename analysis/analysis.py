#!/usr/bin/env python 

import numpy as np
import cv2
import cv2.cv as cv
import sys
import time

if __name__ == "__main__":
    # Debug
    debug = "-d" in sys.argv
    if debug:
        start_time = time.time()

    filename = sys.argv[1]
    result = {}
    image = cv2.imread(filename)
    gray_image = cv2.cvtColor(image, cv.CV_RGB2GRAY)
    gray_image = cv2.equalizeHist(gray_image)

    #Face Detection
    face_cascade = cv2.CascadeClassifier('analysis/haarcascade_frontalface_default.xml')
    faces = face_cascade.detectMultiScale(gray_image, 1.1, 2, cv.CV_HAAR_SCALE_IMAGE,
            (20,20))
    result.setdefault("Faces", len(faces))

    #Colour Analysis
    Z = np.float32(image.reshape((-1, 3)))
    criteria = (cv2.TERM_CRITERIA_EPS, 7, 1.0)
    _, label, center = cv2.kmeans(Z,2,criteria,10,cv2.KMEANS_PP_CENTERS)
    result.setdefault("Primary", [center[0][2], center[0][1], center[0][0]])
    result.setdefault("Secondary", [center[1][2], center[1][1], center[1][0]])

    # Debug
    if debug:
        print("--- %f seconds ---" % (time.time() - start_time))
        center = np.uint8(center)
        res = center[label.flatten()]
        res2 = res.reshape((image.shape))
        cv2.imshow('res2',res2)
        cv2.waitKey(0)
        cv2.destroyAllWindows()

    print(result)
