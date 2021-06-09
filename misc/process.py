from math import sin, cos, radians

v = []

servo = 0

with open("misc/out.txt") as f:
    for line in f:
        if line[0] == "S":
            servo = int(line.split()[2])
            servo = (servo - 2500) * 0.05
        elif line[0] == "L":
            a = float(line.split()[1])
            b = float(line.split()[2])
            c = servo
            if b != 0:
                v.append((a, b, c))

f = open("misc/out2.txt", "w")
v.sort(key=lambda x: x[2])
for e in v:
    f.write("{} {} {}\n".format(e[0], e[1], e[2]))
f.close()
