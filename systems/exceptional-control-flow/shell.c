#include <stdio.h>
#include <stdlib.h>

#define MAXLINE 1024

int main(int argc, char *argv[]) {
  char cmdline[MAXLINE];

  while (1) {
    printf("🦊 ");
    fgets(cmdline, MAXLINE, stdin);
    if (feof(stdin)) {
      printf("\nBye for now!\n");
      exit(0);
    }
    printf("%s", cmdline);
  }
  return 0;
}
