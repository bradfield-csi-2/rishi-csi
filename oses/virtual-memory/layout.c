#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int a = 1, b;
static int c = 1, d;

int main () {
  int e = 1, f;
  static int g = 2, h;
  int *i = malloc(10 * sizeof(int));
  int *j = malloc(10 * sizeof(int));

  /*
  printf("name    location\n");
  printf("&a:     %p\n", &a);
  printf("&b:     %p\n", &b);
  printf("&c:     %p\n", &c);
  printf("&d:     %p\n", &d);
  printf("&e:     %p\n", &e);
  printf("&f:     %p\n", &f);
  printf("&g:     %p\n", &g);
  printf("&h:     %p\n", &h);
  printf("i:      %p\n", i);
  printf("i[1]:   %p\n", &(i[1]));
  printf("i[9]:   %p\n", &(i[9]));
  printf("&i:     %p\n", &i);
  printf("j:      %p\n", j);
  printf("&j:     %p\n", &j);
  printf("main:   %p\n", main);
  printf("printf: %p\n", printf);
*/

  printf("\npid:    %d\n", getpid());
  while (1)
    ;
}
