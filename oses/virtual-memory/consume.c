#include <stdlib.h>
#include <string.h>

#define BYTES_TO_MALLOC 1

int main() {
  for (;;) {
    void *p = malloc(BYTES_TO_MALLOC);
    memset(p, 0, BYTES_TO_MALLOC);
  }
  return 0;
}
