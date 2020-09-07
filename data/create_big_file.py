import os

with open("./test_data.txt", "w") as f:
    f.seek(1024 * 1024 * 1024)
    f.write("end")

