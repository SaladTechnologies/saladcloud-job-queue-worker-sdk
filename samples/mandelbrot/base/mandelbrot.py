import matplotlib
import numpy as np
from PIL import Image


def generate(w, h, iter, re_min, re_max, im_min, im_max, kind):
    x = np.linspace(re_min, re_max, num=w).reshape((1, w))
    y = np.linspace(im_min, im_max, num=h).reshape((h, 1))
    c = np.tile(x, (h, 1)) + 1j * np.tile(y, (1, w))
    z = np.zeros((h, w), dtype="complex")
    m = np.full((h, w), True, dtype="bool")
    n = np.full((h, w), 0)
    for i in range(iter):
        z[m] = z[m] * z[m] + c[m]
        m[np.abs(z) > 2] = False
        n[m] = i
    h = n / iter
    s = np.ones(h.shape)
    v = np.ones(h.shape)
    v[m] = 0
    hsv = np.dstack((h, s, v))
    rgb = matplotlib.colors.hsv_to_rgb(hsv)
    return Image.fromarray((rgb * 255).astype("uint8"))
