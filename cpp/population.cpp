#include <fstream>
#include <iostream>
#include <string>
#include <vector>
#include <cppcoro/generator.hpp>

template<typename String>
auto split_into(std::vector<String>& container, const String& text) -> void;

// auto read_rows(std::istream& in)
//     -> std::vector<std::vector<std::string>> {
//   auto result = std::vector<std::vector<std::string>>{};
//   auto line = std::string{};
//   auto cells = std::vector<std::string>{};
//   while (std::getline(in, line)) {
//     split_into(cells, line);
//     result.push_back(cells);
//   }
//   return result;
// }

auto read_rows(std::istream& in)
    -> cppcoro::generator<const std::vector<std::string_view>&> {
  auto line = std::string{};
  auto cells = std::vector<std::string_view>{};
  while (std::getline(in, line)) {
    split_into(cells, std::string_view{line});
    co_yield cells;
  }
}

auto main() -> int {
  using Population = std::int64_t;
  auto file = std::ifstream{"data/cities500.txt"};
  auto population_south = Population{0};
  auto population_total = Population{0};
  auto row_count = 0;
  for (auto& row: read_rows(file)) {
    // Since libc++ has no float from_chars, rely on small value optimization.
    auto latitude = std::stod(std::string(row.at(4)));
    auto population = std::stol(std::string(row.at(14)));
    population_total += population;
    if (latitude < 0) {
      population_south += population;
    }
    row_count += 1;
  }
  std::cout << "South population: " << population_south << std::endl;
  std::cout << "Total population: " << population_total << std::endl;
  std::cout << "# of rows: " << row_count << std::endl;
}


// Helper code.

template<typename String>
auto split_into(std::vector<String>& container, const String& text) -> void {
  container.clear();
  auto start = std::size_t{0};
  while (start <= text.size()) {
    auto stop = text.find('\t', start);
    if (stop == String::npos) {
      stop = text.size();
    }
    container.push_back(text.substr(start, stop - start));
    start = stop + 1;
  }
}
