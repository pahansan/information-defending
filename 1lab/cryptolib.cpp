#include "cryptolib.h"

#include <cstdint>
#include <random>
#include <utility>

int64_t modExp(int64_t a, int64_t x, int64_t p)
{
    int64_t y = 1;

    while (x != 0) {
        int64_t mod_a = a % p;
        y = x & 1 ? y * mod_a : y;
        a = mod_a * mod_a;
        x >>= 1;
    }

    return y % p;
}

static std::random_device dev;
static std::mt19937_64 gen{dev()};

int testFerma(int64_t a, int64_t p)
{
    if (a >= p)
        return -1;

    return modExp(a, p - 1, p) == 1 ? 0 : 1;
}

bool isProbablyPrime(int64_t x)
{
    if (x <= 1 || x % 2 == 0)
        return false;

    const int iters = 100;
    std::uniform_int_distribution<int64_t> dis(2, x - 1);
    for (int i = 0; i < iters && i < x; i++) {
        int64_t a = dis(gen);
        if (a % x == 0)
            continue;
        if (testFerma(a, x) == 0)
            return false;
    }

    return true;
}

Euclid extendedGCD(int64_t a, int64_t b)
{
    if (a < 1 || b < 1)
        return Euclid{0, 0, 0};

    int64_t u1 = a, u2 = 1, u3 = 0;
    int64_t v1 = b, v2 = 0, v3 = 1;
    if (a < b) {
        std::swap(u1, v1);
        std::swap(u2, v2);
        std::swap(u3, v3);
    }

    while (v1 != 0) {
        int64_t q = u1 / v1;

        int64_t t1 = u1 % v1;
        int64_t t2 = u2 - q * v2;
        int64_t t3 = u3 - q * v3;

        u1 = v1;
        u2 = v2;
        u3 = v3;

        v1 = t1;
        v2 = t2;
        v3 = t3;
    }

    return Euclid{u1, u2, u3};
}

EuclidAndNumbers extendedGCDRandoms()
{
    std::uniform_int_distribution<int64_t> dis(1, 1000);

    int64_t a = 0, b = dis(gen);

    while (a < b)
        a = dis(gen);

    return {a, b, extendedGCD(a, b)};
}

int64_t generatePrime()
{
    std::uniform_int_distribution<int64_t> dis(2, 1000);

    int64_t x = dis(gen);
    while (!isProbablyPrime(x)) {
        x = dis(gen);
    }

    return x;
}

EuclidAndNumbers extendedGCDPrimes()
{
    int64_t a = 0, b = generatePrime();

    while (a < b)
        a = generatePrime();

    return {a, b, extendedGCD(a, b)};
}