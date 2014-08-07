#!/usr/bin/env python 

import numpy as np
import cv2
import cv2.cv as cv
import sys
import time
import colorsys

if __name__ == "__main__":
    # Debug
    debug = "-d" in sys.argv
    if debug:
        start_time = time.time()

    filename = sys.argv[1]
    result = {}
    image = cv2.imread(filename)
    gray_image = cv2.equalizeHist(cv2.cvtColor(image, cv.CV_RGB2GRAY))

    #Face Detection
    face_cascade = cv2.CascadeClassifier('analysis/haarcascade_frontalface_default.xml')
    faces = face_cascade.detectMultiScale(gray_image, 1.1, 2, cv.CV_HAAR_SCALE_IMAGE, (20,20))
    result.setdefault("Faces", len(faces))

    #Colour Analysis
    Z = np.float32(image.reshape((-1, 3)))
    criteria = (cv2.TERM_CRITERIA_EPS, 6, 1.0)
    _, label, center = cv2.kmeans(Z,2,criteria, 7,cv2.KMEANS_PP_CENTERS)
    #Correcting for BGR ordering ):
    primary_rgb = list(center[0])
    secondary_rgb = list(center[1])
    primary_rgb.reverse()
    secondary_rgb.reverse()
    #Calculating HLS
    (primary_hue, lighting, _) = colorsys.rgb_to_hls(primary_rgb[0]/256.0, primary_rgb[1]/256.0, primary_rgb[2]/256.0)
    (secondary_hue, _, _) = colorsys.rgb_to_hls(secondary_rgb[0]/256.0, secondary_rgb[1]/256.0, secondary_rgb[2]/256.0)

    result.setdefault("Primary", {"rgb" : primary_rgb, "hue" : primary_hue*360})
    result.setdefault("Secondary", {"rgb" : secondary_rgb, "hue" : secondary_hue*360})
    result.setdefault("Lighting", lighting*100)

    # Debug
    if debug:
        # Average run time: 2.8s
        print("--- %f seconds ---" % ((time.time() - start_time)/1.0))
        center = np.uint8(center)
        res = center[label.flatten()]
        res2 = res.reshape((image.shape))
        cv2.imshow('res2',res2)
        cv2.waitKey(0)
        cv2.destroyAllWindows()

    print(result)
