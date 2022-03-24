
file = open("utils\output.txt", 'w')

for i in range(1, 16):
    store = f"FullAdder(a = a[{i}], b = b[{i}], c = carry{i - 1}, carry = carry{i}, sum = out[{i}]);\n"
    file.writelines(store)

file.close()