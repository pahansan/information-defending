#include <iostream>

typedef struct
{
    int64_t nod;
    int64_t x;
    int64_t y;
} Evklid_result;

int64_t mod_exp(int64_t a, int64_t x, int64_t p);
int ferma(int64_t a, int64_t p);
Evklid_result Ext_Euc_alg(int64_t a, int64_t b);
Evklid_result Ext_Euc_alg_with_number_generator();
Evklid_result Ext_Euc_alg_with_prime_generator();