#include <libgocoap.h>
#include <stdio.h>

int main () {
    char *version;

    version = coapVersion();
    printf("GoCoapLIB version: %s\n", version);
    return 0;
}