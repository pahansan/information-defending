#include <cstdint>

struct Euclid {
    int64_t gcd;
    int64_t x;
    int64_t y;
};

struct EuclidAndNumbers {
    int64_t a, b;
    Euclid res;
};

int64_t modExp(int64_t a, int64_t x, int64_t p);
bool isProbablyPrime(int64_t x);

Euclid extendedGCD(int64_t a, int64_t b);
EuclidAndNumbers extendedGCDRandoms();
EuclidAndNumbers extendedGCDPrimes();