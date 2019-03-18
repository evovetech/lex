#include <iostream>

extern "C" {
    double sum(double, double);
    double average(double, double);
    double min(double, double);
    double max(double, double);
    double fib(double);
    double callAtan2(double, double);
}

int main() {
    std::cout << "sum of 3.0 and 4.0: " << sum(3, 4) << std::endl;
    std::cout << "average of 3.0 and 4.0: " << average(3.0, 4) << std::endl;
    std::cout << "min of 3.0 and 4.0: " << min(3, 4.0) << std::endl;
    std::cout << "max of 3.0 and 4.0: " << max(3.0, 4.0) << std::endl;
    std::cout << "Compute the 40th fibonacci number: " << fib(40) << std::endl;
    std::cout << "external atan2 of .4 and 42: " << callAtan2(.4, 42) << std::endl;
}
