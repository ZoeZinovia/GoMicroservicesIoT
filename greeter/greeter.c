#include "greeter.h"
#include <stdio.h>

int greet(const char *name, int year) {
    int n = 2;
    
    printf("Greetings, %s from %d! We come in peace :)", name, year);

    return n;
}