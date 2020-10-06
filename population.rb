def read_rows(input)
  return to_enum(:read_rows, input) unless block_given?
  input.each do |line|
    yield line.split("\t")
  end
end

def calc(input)
  population_south = 0
  population_total = 0
  row_count = 0
  rows = read_rows(input)
  # rows.each do |row|
  tally = proc do |row|
  # read_rows(input) do |row|
    latitude = row[4].to_f
    population = row[14].to_i
    population_total += population
    if latitude < 0
      population_south += population
    end
    row_count += 1
  end
  rows.each do |row|
    tally.call(row)
    break if population_total > 1e9
  end
  puts "Reached a billion!"
  rows.each do |row|
    tally.call(row)
  end
  puts "South population: #{population_south}"
  puts "Total population: #{population_total}"
  puts "# of rows: #{row_count}"
end

def read_rows_fiber(input)
  Fiber.new do
    input.each do |line|
      Fiber.yield line.split("\t")
    end
  end
end

def step(input)
  # rows = read_rows_fiber(input)
  # first = rows.resume
  rows = read_rows(input)
  first = rows.next
  puts "lat #{first[4]} with pop #{first[14]}"
end

File.open("data/cities500.txt") do |input|
  calc(input)
  # step(input)
end

# [1, 2, 3].each do |i|
#   puts i
#   break
# end

# numbers = [1, 2, 3].each
# puts numbers.next
