function some_function()
    print("função número 1")
end

function other_funciton()
    print("função número 2")
end

function add(a, b)
    return a + b
end

function fancy_add(msg, a, b)
    print(">>> Fancy add adds but also prints a message: " .. msg)
    return add(a, b)
end
