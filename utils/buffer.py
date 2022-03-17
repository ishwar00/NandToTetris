
file = open("utils\output.txt", 'w')

for i in range(2, 16):
    store = f"    Or(a = in[{i}], b = a{i - 1}, out = a{i});\n"
    file.writelines(store)

file.close()