using BenchmarkTools

# TODO(@tjp): Show matching on 2nd param as str or io?
function readrows(channel::Channel)
    println("a")
    open("data/cities500.txt") do io
        for line in eachline(io)
            put!(channel, split(line, "\t"))
        end
    end
    println("b")
end

function main()
    channel = Channel(readrows)
    println("Hi!")
    n = 0
    for row in channel
        n += 1
    end
    println("there: $n")
end

function main2()
    n = 0
    open("data/cities500.txt") do io
        for line in eachline(io)
            split(line, '\t')
            n += 1
        end
    end
    # println(n)
end

# main()
@time main2()
@btime main2()
