#include "cryptoliba.h"

#include <random>
#include <chrono>
#include <limits>

int64_t mod_exp(int64_t a, int64_t x, int64_t p)
{
    int64_t y = 1;

    while (x != 0)
    {
        int64_t mod_a = a % p;
        y = x & 1 ? y * mod_a : y;
        a = mod_a * mod_a;
        x >>= 1;
    }

    return y % p;
}

int ferma(int64_t a, int64_t p)
{
    if (a >= p)
        return -1;

    return mod_exp(a, p - 1, p) == 1 ? 0 : 1;
}

Evklid_result Ext_Euc_alg(int64_t a, int64_t b)
{
    if (a < b || a < 1 || b < 1)
        return Evklid_result{0, 0, 0};

    int64_t u1 = a, u2 = 1, u3 = 0;
    int64_t v1 = b, v2 = 0, v3 = 1;

    while (v1 != 0)
    {
        int64_t q = u1 / v1;

        int64_t t1 = u1 % v1, t2 = u2 - q * v2, t3 = u3 - q * v3;

        u1 = v1;
        u2 = v2;
        u3 = v3;

        v1 = t1;
        v2 = t2;
        v3 = t3;
    }

    return Evklid_result{u1, u2, u3};
}

static std::mt19937_64 gen{
    static_cast<unsigned long long>(
        std::chrono::system_clock::now().time_since_epoch().count())};

Evklid_result Ext_Euc_alg_with_number_generator()
{
    std::uniform_int_distribution<int64_t> dis(1, 1000);

    int64_t a = 0, b = dis(gen);

    while (a < b)
        a = dis(gen);

    return Ext_Euc_alg(a, b);
}

int64_t generate_prime_number()
{
    std::uniform_int_distribution<int64_t> dis(2, std::numeric_limits<int64_t>::max());

    int64_t x = dis(gen);

    if (x == 2)
        return x;

    while (ferma(x - 2, x) != 0)
    {
        int64_t x = dis(gen);
    }

    return x;
}

Evklid_result Ext_Euc_alg_with_prime_generator()
{
    int64_t a = 0, b = generate_prime_number();

    while (a < b)
        a = generate_prime_number();

    return Ext_Euc_alg(a, b);
}