#include <iostream>

extern "C" {
    double sum(double, double);
    double average(double, double);
}

int main() {
    std::cout << "sum of 3.0 and 4.0: " << sum(3, 4) << std::endl;
    std::cout << "average of 3.0 and 4.0: " << average(3.0, 4.0) << std::endl;
}
