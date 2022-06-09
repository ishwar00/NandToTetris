
file = open("utils/output.txt", 'w')

for i in range(0, 8):
    store = f'    RAM8(in = i, load = R{i}, address = address[2..0], out = out{i});\n'
    file.writelines(store)

file.close()