#include "cryptolib.h"

#include <cstdint>
#include <cstdio>
#include <iostream>

#define BOOL(x) ((x) ? "true" : "false")

int main()
{
    // 1.
    printf("1. Функция быстрого возведения числа в степень по модулю:\n");
    printf("\t%d^%d mod %2d = %ld\n", 5, 12, 7, modExp(5, 12, 7));
    printf("\t%d^%d mod %2d = %ld\n", 3, 21, 11, modExp(3, 21, 11));
    printf("\t%d^%d mod %2d = %ld\n", 7, 31, 17, modExp(7, 31, 17));

    // 2.
    printf("\n2. Тест ферма:\n");
    printf("\t%4d is probably prime: %s\n", 3, BOOL(isProbablyPrime(3)));
    printf("\t%4d is probably prime: %s\n", 2377, BOOL(isProbablyPrime(2377)));
    printf("\t%4d is probably prime: %s\n", 10, BOOL(isProbablyPrime(10)));
    printf("\t%4d is probably prime: %s\n", 11, BOOL(isProbablyPrime(11)));

    // 3.
    printf("\n3. Расширенный алгоритм Евклида:\n");
    auto res = extendedGCD(10, 35);
    printf("\ta = 10, b = 35:\n");
    printf("\tgcd(a, b) = %ld, x = %ld, y = %ld\n", res.gcd, res.x, res.y);

    auto res1 = extendedGCDRandoms();
    printf("\n\tRandom numbers:\n");
    printf("\ta = %ld, b = %ld:\n", res1.a, res1.b);
    printf("\tgcd(a, b) = %ld, x = %ld, y = %ld\n", res1.res.gcd, res1.res.x, res1.res.y);

    auto res2 = extendedGCDPrimes();
    printf("\n\tProbably prime numbers:\n");
    printf("\ta = %ld, b = %ld:\n", res2.a, res2.b);
    printf("\tgcd(a, b) = %ld, x = %ld, y = %ld\n", res2.res.gcd, res2.res.x, res2.res.y);

    printf("\n\tYour numbers:\n\tEnter a, b: ");
    int64_t a, b;
    std::cin >> a >> b;
    auto res3 = extendedGCD(a, b);
    printf("\tgcd(a, b) = %ld, x = %ld, y = %ld\n", res3.gcd, res3.x, res3.y);

    return 0;
}
