import io
import sys
from PIL import Image
import matplotlib.pyplot as plt

if  __name__ == "__main__":
    data = sys.stdin.buffer.read()
    image = Image.open(io.BytesIO(data))
    plt.imshow(image)
    plt.show()